package user

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	"familyquest-backend/internal/api/http/controller"
	"familyquest-backend/internal/api/http/httpcontext"
	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type UserGetter interface {
	Get(ctx context.Context, id entity.UserID) (entity.User, error)
}

type Responder interface {
	Error(w http.ResponseWriter, err error)
	Success(w http.ResponseWriter, statusCode int, data any)
}

type Controller struct {
	protocolMapper toProtocolMapper
	response       Responder
	userGetter     UserGetter
}

func NewController(response Responder, userGetter UserGetter) *Controller {
	return &Controller{
		protocolMapper: newToProtocolMapper(),
		response:       response,
		userGetter:     userGetter,
	}
}

func (c *Controller) RegisterRoutes(router *mux.Router) {
	router.HandleFunc(controller.PathUserMe, c.me).Methods(http.MethodGet)
}

func (c *Controller) me(w http.ResponseWriter, r *http.Request) {
	currentUser, ok := httpcontext.UserFromRequest(r)
	if !ok {
		c.response.Error(w, domain.AuthorizationError())
		return
	}

	user, err := c.userGetter.Get(r.Context(), currentUser.ID)
	if err != nil {
		c.response.Error(w, err)
		return
	}

	c.response.Success(w, http.StatusOK, c.protocolMapper.mapUserResponse(user))
}
