package taskcompletion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"familyquest-backend/internal/api/http/controller"
	"familyquest-backend/internal/api/http/httpcontext"
	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type TaskCompletionSubmitter interface {
	Submit(ctx context.Context, req entity.TaskCompletionSubmitRequest) (entity.TaskCompletion, error)
}

type TaskCompletionLister interface {
	List(ctx context.Context, req entity.TaskCompletionListRequest) ([]entity.TaskCompletionView, error)
}

type TaskCompletionApprover interface {
	Approve(ctx context.Context, req entity.TaskCompletionReviewRequest) (entity.TaskCompletion, error)
}

type TaskCompletionRejecter interface {
	Reject(ctx context.Context, req entity.TaskCompletionReviewRequest) (entity.TaskCompletion, error)
}

type Responder interface {
	Error(w http.ResponseWriter, err error)
	Success(w http.ResponseWriter, statusCode int, data any)
}

type Validator interface {
	ValidateSubmit(req submitRequest) error
	ValidateStatus(status string) (*entity.TaskCompletionStatus, error)
	ValidateReject(req rejectRequest) error
}

type Controller struct {
	submitter        TaskCompletionSubmitter
	lister           TaskCompletionLister
	approver         TaskCompletionApprover
	rejecter         TaskCompletionRejecter
	responder        Responder
	validator        Validator
	toEntityMapper   toEntityMapper
	toProtocolMapper toProtocolMapper
}

func NewController(
	submitter TaskCompletionSubmitter,
	lister TaskCompletionLister,
	approver TaskCompletionApprover,
	rejecter TaskCompletionRejecter,
	responder Responder,
	validator Validator,
) *Controller {
	return &Controller{
		submitter:        submitter,
		lister:           lister,
		approver:         approver,
		rejecter:         rejecter,
		responder:        responder,
		validator:        validator,
		toEntityMapper:   newToEntityMapper(),
		toProtocolMapper: newToProtocolMapper(),
	}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(controller.PathSeasonCompletions, c.list).Methods(http.MethodGet)
	router.HandleFunc(controller.PathSeasonTaskSubmit, c.submit).Methods(http.MethodPost)
	router.HandleFunc(controller.PathTaskCompletionApprove, c.approve).Methods(http.MethodPost)
	router.HandleFunc(controller.PathTaskCompletionReject, c.reject).Methods(http.MethodPost)
}

func (c *Controller) list(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.responder.Error(w, domain.AuthorizationError())
		return
	}

	seasonID, err := parseSeasonID(mux.Vars(r)["season_id"])
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	status, err := c.validator.ValidateStatus(r.URL.Query().Get("status"))
	if err != nil {
		c.responder.Error(w, domain.ValidationError(err.Error()))
		return
	}

	completions, err := c.lister.List(r.Context(), entity.TaskCompletionListRequest{
		UserID:   currentUser.ID,
		SeasonID: seasonID,
		Status:   status,
	})
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	c.responder.Success(w, http.StatusOK, c.toProtocolMapper.mapCompletionsResponse(completions))
}

func (c *Controller) submit(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.responder.Error(w, domain.AuthorizationError())
		return
	}

	seasonTaskID, err := parseSeasonTaskID(mux.Vars(r)["season_task_id"])
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	var req submitRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.responder.Error(w, domain.ValidationError(fmt.Sprintf("can't unmarshal request: %v", err)))
		return
	}

	if err = c.validator.ValidateSubmit(req); err != nil {
		c.responder.Error(w, domain.ValidationError(err.Error()))
		return
	}

	completion, err := c.submitter.Submit(
		r.Context(),
		c.toEntityMapper.mapSubmitRequest(req, currentUser.ID, seasonTaskID),
	)
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	c.responder.Success(w, http.StatusOK, c.toProtocolMapper.mapCompletionResponse(completion))
}

func (c *Controller) approve(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.responder.Error(w, domain.AuthorizationError())
		return
	}

	completionID, err := parseCompletionID(mux.Vars(r)["completion_id"])
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	completion, err := c.approver.Approve(r.Context(), entity.TaskCompletionReviewRequest{
		UserID:       currentUser.ID,
		CompletionID: completionID,
	})
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	c.responder.Success(w, http.StatusOK, c.toProtocolMapper.mapCompletionResponse(completion))
}

func (c *Controller) reject(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.responder.Error(w, domain.AuthorizationError())
		return
	}

	completionID, err := parseCompletionID(mux.Vars(r)["completion_id"])
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	var req rejectRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.responder.Error(w, domain.ValidationError(fmt.Sprintf("can't unmarshal request: %v", err)))
		return
	}
	if err = c.validator.ValidateReject(req); err != nil {
		c.responder.Error(w, domain.ValidationError(err.Error()))
		return
	}

	completion, err := c.rejecter.Reject(r.Context(), entity.TaskCompletionReviewRequest{
		UserID:       currentUser.ID,
		CompletionID: completionID,
		Comment:      strings.TrimSpace(req.Comment),
	})
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	c.responder.Success(w, http.StatusOK, c.toProtocolMapper.mapCompletionResponse(completion))
}

func parseSeasonID(value string) (entity.SeasonID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return entity.SeasonID{}, domain.ValidationError(`the "season_id" path parameter must be uuid`)
	}

	return entity.SeasonID(id), nil
}

func parseCompletionID(value string) (entity.TaskCompletionID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return entity.TaskCompletionID{}, domain.ValidationError(`the "completion_id" path parameter must be uuid`)
	}

	return entity.TaskCompletionID(id), nil
}

func parseSeasonTaskID(value string) (entity.SeasonTaskID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return entity.SeasonTaskID{}, domain.ValidationError(`the "season_task_id" path parameter must be uuid`)
	}

	return entity.SeasonTaskID(id), nil
}
