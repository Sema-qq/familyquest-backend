package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"familyquest-backend/internal/api/http/httpcontext"
	"familyquest-backend/internal/domain"
	"familyquest-backend/internal/domain/entity"
)

type TokenParser interface {
	Parse(token string) (entity.UserID, error)
}

type UserGetter interface {
	GetByID(ctx context.Context, id entity.UserID) (entity.User, error)
}

type Responder interface {
	Error(w http.ResponseWriter, err error)
}

type Middleware struct {
	tokenParser TokenParser
	userGetter  UserGetter
	responder   Responder
}

func New(tokenParser TokenParser, userGetter UserGetter, responder Responder) *Middleware {
	return &Middleware{
		tokenParser: tokenParser,
		userGetter:  userGetter,
		responder:   responder,
	}
}

func (m *Middleware) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r.Header.Get("Authorization"))
		if token == "" {
			m.responder.Error(w, domain.AuthorizationError())
			return
		}

		userID, err := m.tokenParser.Parse(token)
		if err != nil {
			m.responder.Error(w, domain.AuthorizationError())
			return
		}

		user, err := m.userGetter.GetByID(r.Context(), userID)
		if err != nil {
			m.responder.Error(w, fmt.Errorf("can't get user by id: %w", err))
			return
		}

		next.ServeHTTP(w, r.WithContext(httpcontext.WithUser(r.Context(), user)))
	})
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
