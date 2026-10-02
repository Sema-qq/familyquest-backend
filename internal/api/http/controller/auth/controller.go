package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"familyquest-backend/internal/api/http/controller"
	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type UserCreator interface {
	Create(ctx context.Context, req entity.UserCreateRequest) (entity.User, error)
}

type UserAuthenticator interface {
	Login(ctx context.Context, req entity.AuthRequest) (entity.AuthResult, error)
}

type Responder interface {
	Error(w http.ResponseWriter, err error)
	Success(w http.ResponseWriter, statusCode int, data any)
}

type Validator interface {
	validateRegister(req registerRequest) error
	validateLogin(req loginRequest) error
}

type Controller struct {
	entityMapper   *toEntityMapper
	protocolMapper *toProtocolMapper
	validator      Validator
	response       Responder
	userCreator    UserCreator
	authenticator  UserAuthenticator
}

func NewController(
	validator Validator,
	response Responder,
	userCreator UserCreator,
	authenticator UserAuthenticator,
) *Controller {
	return &Controller{
		entityMapper:   newToEntityMapper(),
		protocolMapper: newToProtocolMapper(),
		validator:      validator,
		response:       response,
		userCreator:    userCreator,
		authenticator:  authenticator,
	}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(controller.PathAuthRegister, c.register).Methods(http.MethodPost)
	router.HandleFunc(controller.PathAuthLogin, c.login).Methods(http.MethodPost)
}

func (c *Controller) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.response.Error(w, domain.ValidationError(fmt.Sprintf("can't unmarshal request: %v", err)))
		return
	}

	if err := c.validator.validateRegister(req); err != nil {
		c.response.Error(w, domain.ValidationError(err.Error()))
		return
	}

	user, err := c.userCreator.Create(r.Context(), c.entityMapper.mapUserCreateRequest(req))
	if err != nil {
		c.response.Error(w, err)
		return
	}

	c.response.Success(w, http.StatusCreated, c.protocolMapper.mapRegisterResponse(user))
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.response.Error(w, domain.ValidationError(fmt.Sprintf("can't unmarshal request: %v", err)))
		return
	}

	if err := c.validator.validateLogin(req); err != nil {
		c.response.Error(w, domain.ValidationError(err.Error()))
		return
	}

	result, err := c.authenticator.Login(r.Context(), c.entityMapper.mapAuthRequest(req))
	if err != nil {
		c.response.Error(w, err)
		return
	}

	c.response.Success(w, http.StatusOK, c.protocolMapper.mapLoginResponse(result))
}
