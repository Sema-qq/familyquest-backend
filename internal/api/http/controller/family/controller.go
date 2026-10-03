package family

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

type FamilyCreator interface {
	Create(ctx context.Context, req entity.FamilyCreateRequest) (entity.FamilyWithRole, error)
}

type FamilyGetter interface {
	Get(ctx context.Context, userID entity.UserID) (entity.FamilyWithRole, error)
}

type FamilyMemberCreator interface {
	Create(ctx context.Context, req entity.FamilyMemberCreateRequest) (entity.FamilyMemberProfile, error)
}

type FamilyMemberLister interface {
	List(ctx context.Context, userID entity.UserID) ([]entity.FamilyMemberProfile, error)
}

type Responder interface {
	Error(w http.ResponseWriter, err error)
	Success(w http.ResponseWriter, statusCode int, data any)
}

type Validator interface {
	ValidateCreate(req createRequest) error
	ValidateCreateMember(req createMemberRequest) error
}

type Controller struct {
	entityMapper        toEntityMapper
	protocolMapper      toProtocolMapper
	validator           Validator
	response            Responder
	familyCreator       FamilyCreator
	familyGetter        FamilyGetter
	familyMemberCreator FamilyMemberCreator
	familyMemberLister  FamilyMemberLister
}

func NewController(
	validator Validator,
	response Responder,
	familyCreator FamilyCreator,
	familyGetter FamilyGetter,
	familyMemberCreator FamilyMemberCreator,
	familyMemberLister FamilyMemberLister,
) *Controller {
	return &Controller{
		entityMapper:        newToEntityMapper(),
		protocolMapper:      newToProtocolMapper(),
		validator:           validator,
		response:            response,
		familyCreator:       familyCreator,
		familyGetter:        familyGetter,
		familyMemberCreator: familyMemberCreator,
		familyMemberLister:  familyMemberLister,
	}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(controller.PathFamily, c.create).Methods(http.MethodPost)
	router.HandleFunc(controller.PathFamily, c.get).Methods(http.MethodGet)
	router.HandleFunc(controller.PathFamilyMembers, c.createMember).Methods(http.MethodPost)
	router.HandleFunc(controller.PathFamilyMembers, c.listMembers).Methods(http.MethodGet)
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.response.Error(w, domain.AuthorizationError())
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.response.Error(w, domain.ValidationError(fmt.Sprintf("can't unmarshal request: %v", err)))
		return
	}

	if err := c.validator.ValidateCreate(req); err != nil {
		c.response.Error(w, domain.ValidationError(err.Error()))
		return
	}

	family, err := c.familyCreator.Create(r.Context(), c.entityMapper.mapFamilyCreateRequest(req, currentUser.ID))
	if err != nil {
		c.response.Error(w, err)
		return
	}

	c.response.Success(w, http.StatusCreated, c.protocolMapper.mapFamilyResponse(family))
}

func (c *Controller) get(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.response.Error(w, domain.AuthorizationError())
		return
	}

	family, err := c.familyGetter.Get(r.Context(), currentUser.ID)
	if err != nil {
		c.response.Error(w, err)
		return
	}

	c.response.Success(w, http.StatusOK, c.protocolMapper.mapFamilyResponse(family))
}

func (c *Controller) createMember(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.response.Error(w, domain.AuthorizationError())
		return
	}

	var req createMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.response.Error(w, domain.ValidationError(fmt.Sprintf("can't unmarshal request: %v", err)))
		return
	}

	if err := c.validator.ValidateCreateMember(req); err != nil {
		c.response.Error(w, domain.ValidationError(err.Error()))
		return
	}

	member, err := c.familyMemberCreator.Create(r.Context(), c.entityMapper.mapFamilyMemberCreateRequest(req, currentUser.ID))
	if err != nil {
		c.response.Error(w, err)
		return
	}

	c.response.Success(w, http.StatusCreated, c.protocolMapper.mapFamilyMemberResponse(member))
}

func (c *Controller) listMembers(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.response.Error(w, domain.AuthorizationError())
		return
	}

	members, err := c.familyMemberLister.List(r.Context(), currentUser.ID)
	if err != nil {
		c.response.Error(w, err)
		return
	}

	c.response.Success(w, http.StatusOK, c.protocolMapper.mapFamilyMembersResponse(members))
}
