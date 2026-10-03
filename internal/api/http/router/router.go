package router

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

type HealthRepository interface {
	Ping(ctx context.Context) error
}

type AuthController interface {
	RegisterRoutes(router *mux.Router)
}

type UserController interface {
	RegisterRoutes(router *mux.Router)
}

type FamilyController interface {
	RegisterRoutes(router *mux.Router)
}

type AuthMiddleware interface {
	Check(next http.Handler) http.Handler
}

type Controller interface {
	RegisterRoutes(router *mux.Router)
}

type Controllers []Controller

func New(
	healthRepository HealthRepository,
	authController AuthController,
	publicControllers Controllers,
	authMiddleware AuthMiddleware,
) http.Handler {
	r := mux.NewRouter()
	authController.RegisterRoutes(r)

	protectedRouter := r.NewRoute().Subrouter()
	protectedRouter.Use(authMiddleware.Check)

	for _, controller := range publicControllers {
		controller.RegisterRoutes(protectedRouter)
	}

	r.HandleFunc("/api/v1/hello", hello).Methods(http.MethodGet)
	r.HandleFunc("/app/health", health).Methods(http.MethodGet)
	r.HandleFunc("/app/ready", ready(healthRepository)).Methods(http.MethodGet)

	return r
}

func hello(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"hello from familyquest"}`))
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func ready(healthRepository HealthRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := healthRepository.Ping(req.Context()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error":"postgres is not ready"}`))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	}
}
