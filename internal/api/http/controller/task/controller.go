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

type TaskLister interface {
	List(ctx context.Context, userID entity.UserID) ([]entity.Task, error)
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
	taskLister       TaskLister
	responder        Responder
	validator        Validator
	toEntityMapper   toEntityMapper
	toProtocolMapper toProtocolMapper
}

func NewController(taskAdder TaskAdder, taskLister TaskLister, responder Responder, validator Validator) *Controller {
	return &Controller{
		taskAdder:        taskAdder,
		taskLister:       taskLister,
		responder:        responder,
		validator:        validator,
		toEntityMapper:   newToEntityMapper(),
		toProtocolMapper: newToProtocolMapper(),
	}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(controller.PathTasks, c.list).Methods(http.MethodGet)
	router.HandleFunc(controller.PathTaskAdd, c.add).Methods(http.MethodPost)
}

func (c *Controller) list(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.responder.Error(w, domain.AuthorizationError())
		return
	}

	tasks, err := c.taskLister.List(r.Context(), currentUser.ID)
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	c.responder.Success(w, http.StatusOK, c.toProtocolMapper.mapTasksResponse(tasks))
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
