package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	authcontroller "familyquest-backend/internal/api/http/controller/auth"
	familycontroller "familyquest-backend/internal/api/http/controller/family"
	usercontroller "familyquest-backend/internal/api/http/controller/user"
	authmiddleware "familyquest-backend/internal/api/http/middleware/auth"
	"familyquest-backend/internal/api/http/response"
	httprouter "familyquest-backend/internal/api/http/router"
	familyrepository "familyquest-backend/internal/repository/family"
	familymemberrepository "familyquest-backend/internal/repository/familymember"
	healthrepository "familyquest-backend/internal/repository/postgres/health"
	userrepository "familyquest-backend/internal/repository/user"
	authlogin "familyquest-backend/internal/usecase/auth/login"
	familycreate "familyquest-backend/internal/usecase/family/create"
	familyget "familyquest-backend/internal/usecase/family/get"
	familymembercreate "familyquest-backend/internal/usecase/family/membercreate"
	familymemberlist "familyquest-backend/internal/usecase/family/memberlist"
	usercreate "familyquest-backend/internal/usecase/user/create"
	userget "familyquest-backend/internal/usecase/user/get"
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

	txConn := db.NewTxAwareConn(postgresConn)
	txManager := db.NewTxManager(postgresConn)

	userRepository := userrepository.NewRepository(txConn)
	familyRepository := familyrepository.NewRepository(txConn)
	familyMemberRepository := familymemberrepository.NewRepository(txConn)
	passwordHasher := password.NewHasher()
	uuidGenerator := uuidprovider.New()
	timeGenerator := timeprovider.New()
	tokenGenerator := token.NewHMACGenerator(
		cfg.Auth.JWTSecret,
		cfg.Auth.JWTIssuer,
		cfg.Auth.AccessTokenTTL,
	)

	userCreateUC := usercreate.New(userRepository, passwordHasher, uuidGenerator)
	userGetUC := userget.New(userRepository)
	authLoginUC := authlogin.New(userRepository, passwordHasher, timeGenerator, tokenGenerator)
	familyCreateUC := familycreate.New(familyRepository, familyMemberRepository, txManager, uuidGenerator)
	familyGetUC := familyget.New(familyRepository)
	familyMemberCreateUC := familymembercreate.New(
		userRepository,
		familyMemberRepository,
		txManager,
		passwordHasher,
		uuidGenerator,
	)
	familyMemberListUC := familymemberlist.New(familyMemberRepository)

	httpResponse := response.NewPublicResponder()
	authValidator := authcontroller.NewRequestValidator()
	authController := authcontroller.NewController(authValidator, httpResponse, userCreateUC, authLoginUC)
	userController := usercontroller.NewController(httpResponse, userGetUC)
	familyValidator := familycontroller.NewRequestValidator()
	familyController := familycontroller.NewController(
		familyValidator,
		httpResponse,
		familyCreateUC,
		familyGetUC,
		familyMemberCreateUC,
		familyMemberListUC,
	)
	authMiddleware := authmiddleware.New(tokenGenerator, userRepository, httpResponse)

	healthRepository := healthrepository.NewRepository(postgresConn)
	router := httprouter.New(healthRepository, authController, userController, familyController, authMiddleware)
	server := NewHTTPServer(cfg, router)

	RunHTTPServer(ctx, server, cfg.HTTP.ShutdownTimeout)
}
