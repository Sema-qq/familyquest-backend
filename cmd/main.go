package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	authcontroller "familyquest-backend/internal/api/http/controller/auth"
	familycontroller "familyquest-backend/internal/api/http/controller/family"
	seasoncontroller "familyquest-backend/internal/api/http/controller/season"
	taskcontroller "familyquest-backend/internal/api/http/controller/task"
	usercontroller "familyquest-backend/internal/api/http/controller/user"
	authmiddleware "familyquest-backend/internal/api/http/middleware/auth"
	"familyquest-backend/internal/api/http/response"
	httprouter "familyquest-backend/internal/api/http/router"
	familyrepository "familyquest-backend/internal/repository/family"
	familymemberrepository "familyquest-backend/internal/repository/familymember"
	healthrepository "familyquest-backend/internal/repository/postgres/health"
	seasonrepository "familyquest-backend/internal/repository/season"
	taskrepository "familyquest-backend/internal/repository/task"
	userrepository "familyquest-backend/internal/repository/user"
	authlogin "familyquest-backend/internal/usecase/auth/login"
	familycreate "familyquest-backend/internal/usecase/family/create"
	familyget "familyquest-backend/internal/usecase/family/get"
	familymembercreate "familyquest-backend/internal/usecase/family/membercreate"
	familymemberlist "familyquest-backend/internal/usecase/family/memberlist"
	seasoncreate "familyquest-backend/internal/usecase/season/create"
	seasonget "familyquest-backend/internal/usecase/season/get"
	seasonlist "familyquest-backend/internal/usecase/season/list"
	taskadd "familyquest-backend/internal/usecase/task/add"
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
	taskRepository := taskrepository.NewRepository(txConn)
	seasonRepository := seasonrepository.NewRepository(txConn)
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
	taskAddUC := taskadd.New(taskRepository, familyMemberRepository, uuidGenerator)
	seasonCreateUC := seasoncreate.New(seasonRepository, familyRepository, familyMemberRepository, uuidGenerator)
	seasonListUC := seasonlist.New(seasonRepository, familyMemberRepository)
	seasonGetUC := seasonget.New(seasonRepository, familyMemberRepository)

	httpResponse := response.NewPublicResponder()
	authController := authcontroller.NewController(authcontroller.NewRequestValidator(), httpResponse, userCreateUC, authLoginUC)
	userController := usercontroller.NewController(httpResponse, userGetUC)
	familyController := familycontroller.NewController(
		familycontroller.NewRequestValidator(),
		httpResponse,
		familyCreateUC,
		familyGetUC,
		familyMemberCreateUC,
		familyMemberListUC,
	)
	taskController := taskcontroller.NewController(taskAddUC, httpResponse, taskcontroller.NewRequestValidator())
	seasonController := seasoncontroller.NewController(
		seasonCreateUC,
		seasonListUC,
		seasonGetUC,
		httpResponse,
		seasoncontroller.NewRequestValidator(),
	)
	authMiddleware := authmiddleware.New(tokenGenerator, userRepository, httpResponse)

	healthRepository := healthrepository.NewRepository(postgresConn)
	publicControllers := []httprouter.Controller{userController, familyController, taskController, seasonController}
	router := httprouter.New(healthRepository, authController, publicControllers, authMiddleware)
	server := NewHTTPServer(cfg, router)

	RunHTTPServer(ctx, server, cfg.HTTP.ShutdownTimeout)
}
