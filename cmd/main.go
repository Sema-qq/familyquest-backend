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
	taskcompletioncontroller "familyquest-backend/internal/api/http/controller/taskcompletion"
	usercontroller "familyquest-backend/internal/api/http/controller/user"
	authmiddleware "familyquest-backend/internal/api/http/middleware/auth"
	"familyquest-backend/internal/api/http/response"
	httprouter "familyquest-backend/internal/api/http/router"
	familyrepository "familyquest-backend/internal/repository/family"
	familymemberrepository "familyquest-backend/internal/repository/familymember"
	healthrepository "familyquest-backend/internal/repository/postgres/health"
	seasonrepository "familyquest-backend/internal/repository/season"
	seasontaskrepository "familyquest-backend/internal/repository/seasontask"
	seasontaskslotrepository "familyquest-backend/internal/repository/seasontaskslot"
	taskrepository "familyquest-backend/internal/repository/task"
	taskcompletionrepository "familyquest-backend/internal/repository/taskcompletion"
	userrepository "familyquest-backend/internal/repository/user"
	authlogin "familyquest-backend/internal/usecase/auth/login"
	familycreate "familyquest-backend/internal/usecase/family/create"
	familyget "familyquest-backend/internal/usecase/family/get"
	familymembercreate "familyquest-backend/internal/usecase/family/membercreate"
	familymemberlist "familyquest-backend/internal/usecase/family/memberlist"
	seasonactivate "familyquest-backend/internal/usecase/season/activate"
	seasoncomplete "familyquest-backend/internal/usecase/season/complete"
	seasoncreate "familyquest-backend/internal/usecase/season/create"
	seasonget "familyquest-backend/internal/usecase/season/get"
	seasonlist "familyquest-backend/internal/usecase/season/list"
	seasonresults "familyquest-backend/internal/usecase/season/results"
	seasontaskadd "familyquest-backend/internal/usecase/season/taskadd"
	taskadd "familyquest-backend/internal/usecase/task/add"
	tasklist "familyquest-backend/internal/usecase/task/list"
	taskcompletionapprove "familyquest-backend/internal/usecase/taskcompletion/approve"
	taskcompletionlist "familyquest-backend/internal/usecase/taskcompletion/list"
	taskcompletionreject "familyquest-backend/internal/usecase/taskcompletion/reject"
	taskcompletionsubmit "familyquest-backend/internal/usecase/taskcompletion/submit"
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
	seasonTaskRepository := seasontaskrepository.NewRepository(txConn)
	seasonTaskSlotRepository := seasontaskslotrepository.NewRepository(txConn)
	taskCompletionRepository := taskcompletionrepository.NewRepository(txConn)
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
	taskListUC := tasklist.New(taskRepository, familyMemberRepository)
	seasonCreateUC := seasoncreate.New(seasonRepository, familyRepository, familyMemberRepository, uuidGenerator)
	seasonListUC := seasonlist.New(seasonRepository, familyMemberRepository)
	seasonGetUC := seasonget.New(
		seasonRepository,
		familyMemberRepository,
		seasonTaskRepository,
		seasonTaskSlotRepository,
		taskCompletionRepository,
	)
	seasonActivateUC := seasonactivate.New(
		seasonRepository,
		familyMemberRepository,
		seasonTaskRepository,
		seasonTaskSlotRepository,
	)
	seasonResultsUC := seasonresults.New(familyMemberRepository, seasonRepository, taskCompletionRepository)
	seasonCompleteUC := seasoncomplete.New(familyMemberRepository, seasonRepository, taskCompletionRepository, timeGenerator)
	seasonTaskAddUC := seasontaskadd.New(
		familyMemberRepository,
		seasonRepository,
		taskRepository,
		seasonTaskRepository,
		seasonTaskSlotRepository,
		txManager,
		uuidGenerator,
	)
	taskCompletionSubmitUC := taskcompletionsubmit.New(
		familyMemberRepository,
		seasonTaskRepository,
		seasonTaskSlotRepository,
		taskCompletionRepository,
		txManager,
		timeGenerator,
		uuidGenerator,
	)
	taskCompletionListUC := taskcompletionlist.New(familyMemberRepository, seasonRepository, taskCompletionRepository)
	taskCompletionApproveUC := taskcompletionapprove.New(familyMemberRepository, taskCompletionRepository, timeGenerator)
	taskCompletionRejectUC := taskcompletionreject.New(familyMemberRepository, taskCompletionRepository, timeGenerator)

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
	taskController := taskcontroller.NewController(taskAddUC, taskListUC, httpResponse, taskcontroller.NewRequestValidator())
	seasonController := seasoncontroller.NewController(
		seasonCreateUC,
		seasonListUC,
		seasonGetUC,
		seasonTaskAddUC,
		seasonActivateUC,
		seasonCompleteUC,
		seasonResultsUC,
		httpResponse,
		seasoncontroller.NewRequestValidator(),
	)
	taskCompletionController := taskcompletioncontroller.NewController(
		taskCompletionSubmitUC,
		taskCompletionListUC,
		taskCompletionApproveUC,
		taskCompletionRejectUC,
		httpResponse,
		taskcompletioncontroller.NewRequestValidator(),
	)
	authMiddleware := authmiddleware.New(tokenGenerator, userRepository, httpResponse)

	healthRepository := healthrepository.NewRepository(postgresConn)
	publicControllers := []httprouter.Controller{
		userController,
		familyController,
		taskController,
		seasonController,
		taskCompletionController,
	}
	router := httprouter.New(healthRepository, authController, publicControllers, authMiddleware)
	server := NewHTTPServer(cfg, router)

	RunHTTPServer(ctx, server, cfg.HTTP.ShutdownTimeout)
}
