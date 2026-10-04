package season

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"familyquest-backend/internal/api/http/controller"
	"familyquest-backend/internal/api/http/httpcontext"
	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type SeasonCreator interface {
	Create(ctx context.Context, req entity.SeasonCreateRequest) (entity.Season, error)
}

type SeasonLister interface {
	List(ctx context.Context, userID entity.UserID) ([]entity.Season, error)
}

type SeasonGetter interface {
	Get(ctx context.Context, userID entity.UserID, seasonID entity.SeasonID) (entity.Season, error)
}

type Responder interface {
	Error(w http.ResponseWriter, err error)
	Success(w http.ResponseWriter, statusCode int, data any)
}

type Validator interface {
	ValidateCreate(req createRequest) error
}

type Controller struct {
	seasonCreator    SeasonCreator
	seasonLister     SeasonLister
	seasonGetter     SeasonGetter
	responder        Responder
	validator        Validator
	toEntityMapper   toEntityMapper
	toProtocolMapper toProtocolMapper
}

func NewController(
	seasonCreator SeasonCreator,
	seasonLister SeasonLister,
	seasonGetter SeasonGetter,
	responder Responder,
	validator Validator,
) *Controller {
	return &Controller{
		seasonCreator:    seasonCreator,
		seasonLister:     seasonLister,
		seasonGetter:     seasonGetter,
		responder:        responder,
		validator:        validator,
		toEntityMapper:   newToEntityMapper(),
		toProtocolMapper: newToProtocolMapper(),
	}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(controller.PathSeasons, c.create).Methods(http.MethodPost)
	router.HandleFunc(controller.PathSeasons, c.list).Methods(http.MethodGet)
	router.HandleFunc(controller.PathSeason, c.get).Methods(http.MethodGet)
}

func (c *Controller) create(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.responder.Error(w, domain.AuthorizationError())
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.responder.Error(w, domain.ValidationError(fmt.Sprintf("can't unmarshal request: %v", err)))
		return
	}

	if err := c.validator.ValidateCreate(req); err != nil {
		c.responder.Error(w, domain.ValidationError(err.Error()))
		return
	}

	season, err := c.seasonCreator.Create(r.Context(), c.toEntityMapper.mapCreateRequest(req, currentUser.ID))
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	c.responder.Success(w, http.StatusCreated, c.toProtocolMapper.mapSeasonResponse(season))
}

func (c *Controller) list(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.responder.Error(w, domain.AuthorizationError())
		return
	}

	seasons, err := c.seasonLister.List(r.Context(), currentUser.ID)
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	c.responder.Success(w, http.StatusOK, c.toProtocolMapper.mapSeasonsResponse(seasons))
}

func (c *Controller) get(w http.ResponseWriter, r *http.Request) {
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

	season, err := c.seasonGetter.Get(r.Context(), currentUser.ID, seasonID)
	if err != nil {
		c.responder.Error(w, err)
		return
	}

	c.responder.Success(w, http.StatusOK, c.toProtocolMapper.mapSeasonDetailResponse(season))
}

func parseSeasonID(value string) (entity.SeasonID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return entity.SeasonID{}, domain.ValidationError(`the "season_id" path parameter must be uuid`)
	}

	return entity.SeasonID(id), nil
}
