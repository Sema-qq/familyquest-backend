package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	authcontroller "familyquest-backend/internal/api/http/controller/auth"
	"familyquest-backend/internal/api/http/response"
	httprouter "familyquest-backend/internal/api/http/router"
	healthrepository "familyquest-backend/internal/repository/postgres/health"
	userrepository "familyquest-backend/internal/repository/user"
	authlogin "familyquest-backend/internal/usecase/auth/login"
	usercreate "familyquest-backend/internal/usecase/user/create"
	"familyquest-backend/pkg/db"
	"familyquest-backend/pkg/password"
	"familyquest-backend/pkg/timeprovider"
	"familyquest-backend/pkg/token"
	"familyquest-backend/pkg/uuidprovider"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg, err := LoadConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	postgresConn, err := NewPostgresConn(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer postgresConn.Close()

	txManager := db.NewTxManager(postgresConn)
	_ = txManager

	userRepository := userrepository.NewRepository(postgresConn)
	passwordHasher := password.NewHasher()
	uuidGenerator := uuidprovider.New()
	timeGenerator := timeprovider.New()
	tokenGenerator := token.NewHMACGenerator(
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTIssuer,
		cfg.Auth.AccessTokenTTL,
	)

	userCreateUC := usercreate.New(userRepository, passwordHasher, uuidGenerator)
	authLoginUC := authlogin.New(userRepository, passwordHasher, timeGenerator, tokenGenerator)
	httpResponse := response.NewPublicResponder()
	authValidator := authcontroller.NewRequestValidator()
	authController := authcontroller.NewController(authValidator, httpResponse, userCreateUC, authLoginUC)

	healthRepository := healthrepository.NewRepository(postgresConn)
	router := httprouter.New(healthRepository, authController)
	server := NewHTTPServer(cfg, router)

	RunHTTPServer(ctx, server, cfg.HTTP.ShutdownTimeout)
}
