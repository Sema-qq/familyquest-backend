package httpcontext

import (
	"context"
	"net/http"

	"familyquest-backend/internal/domain/entity"
)

type userContextKey struct{}

func WithUser(ctx context.Context, user entity.User) context.Context {
	return context.WithValue(ctx, userContextKey{}, user)
}

func UserFromContext(ctx context.Context) (entity.User, bool) {
	user, ok := ctx.Value(userContextKey{}).(entity.User)
	return user, ok
}

func UserFromRequest(r *http.Request) (entity.User, bool) {
	return UserFromContext(r.Context())
}
