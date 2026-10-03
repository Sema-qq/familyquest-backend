package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"familyquest-backend/internal/api/http/controller"
	"familyquest-backend/internal/api/http/httpcontext"
	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type TaskAdder interface {
	Add(ctx context.Context, model entity.TaskAddRequest) (entity.Task, error)
}

type Responder interface {
	Error(w http.ResponseWriter, err error)
	Success(w http.ResponseWriter, statusCode int, data any)
}

type Validator interface {
	ValidateAdd(req addRequest) error
}

type Controller struct {
	taskAdder        TaskAdder
	responder        Responder
	validator        Validator
	toEntityMapper   toEntityMapper
	toProtocolMapper toProtocolMapper
}

func NewController(taskAdder TaskAdder, responder Responder, validator Validator) *Controller {
	return &Controller{
		taskAdder:        taskAdder,
		responder:        responder,
		validator:        validator,
		toEntityMapper:   newToEntityMapper(),
		toProtocolMapper: newToProtocolMapper(),
	}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(controller.PathTaskAdd, c.add).Methods(http.MethodPost)
}

func (c *Controller) add(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.responder.Error(w, domain.AuthorizationError())
		return
	}

	var req addRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.responder.Error(w, domain.ValidationError(fmt.Sprintf("can't unmarshal request: %v", err)))
		return
	}

	if err := c.validator.ValidateAdd(req); err != nil {
		c.responder.Error(w, domain.ValidationError(err.Error()))
		return
	}

	task, err := c.taskAdder.Add(r.Context(), c.toEntityMapper.mapAddTask(req, currentUser.ID))
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	c.responder.Success(w, http.StatusCreated, c.toProtocolMapper.mapAddResponse(task))
}
