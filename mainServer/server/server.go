package server

import (
	"database/sql"
	"fmt"
	"mainServer/controllers"
	"mainServer/db"
	"mainServer/repositories"
	"mainServer/repositories/interfaces"
	"mainServer/repositories/postgres"
	"mainServer/server/config"
	"mainServer/services"
	servinterfaces "mainServer/services/interfaces"
)

type RepoEnv struct {
	git           repositories.GitRepository
	article       interfaces.ArticleRepository
	user          interfaces.UserRepository
	version       interfaces.VersionRepository
	req           interfaces.RequestRepository
	thread        interfaces.ThreadRepository
	comment       interfaces.CommentRepository
	commitThread  interfaces.CommitThreadRepository
	requestThread interfaces.RequestThreadRepository
	reviewThread  interfaces.ReviewThreadRepository
}

type ServiceEnv struct {
	version       servinterfaces.VersionService
	article       services.ArticleService
	user          services.UserService
	req           services.RequestService
	thread        services.ThreadService
	comment       services.CommentService
	commitThread  services.CommitThreadService
	requestThread services.RequestThreadService
	reviewThread  services.ReviewThreadService
}

type ControllerEnv struct {
	article controllers.ArticleController
	version controllers.VersionController
	user    controllers.UserController
	req     controllers.RequestController
	thread  controllers.ThreadController
}

func initRepoEnv(cfg *config.Config, database *sql.DB) RepoEnv {
	return RepoEnv{
		git:           repositories.NewGitRepository(&cfg.Git),
		article:       postgres.NewPgArticleRepository(database),
		user:          postgres.NewPgUserRepository(database),
		version:       postgres.NewPgVersionRepository(database),
		req:           postgres.NewPgRequestRepository(database),
		thread:        postgres.NewPgThreadRepository(database),
		comment:       postgres.NewPgCommentRepository(database),
		commitThread:  postgres.NewPgCommitThreadRepository(database),
		requestThread: postgres.NewPgRequestThreadRepository(database),
		reviewThread:  postgres.NewPgReviewThreadRepository(database),
	}
}

func initServiceEnv(repos RepoEnv) ServiceEnv {
	return ServiceEnv{
		article:       services.NewArticleService(repos.article, repos.version, repos.git),
		user:          services.UserService{UserRepository: repos.user},
		req:           services.RequestService{Repo: repos.req, Versionrepo: repos.version, Gitrepo: repos.git},
		version:       services.VersionService{Gitrepo: repos.git, Versionrepo: repos.version},
		thread:        services.ThreadService{ThreadRepository: repos.thread},
		comment:       services.CommentService{CommentRepository: repos.comment},
		commitThread:  services.CommitThreadService{CommitThreadRepository: repos.commitThread},
		requestThread: services.RequestThreadService{RequestThreadRepository: repos.requestThread},
		reviewThread:  services.ReviewThreadService{ReviewThreadRepository: repos.reviewThread},
	}
}

func initControllerEnv(cfg *config.Config, servs ServiceEnv) ControllerEnv {
	return ControllerEnv{
		article: controllers.NewArticleController(servs.article),
		user:    controllers.UserController{UserService: servs.user, Cfg: cfg},
		req:     controllers.RequestController{Serv: servs.req},
		version: controllers.VersionController{Serv: servs.version},
		thread: controllers.ThreadController{ThreadService: servs.thread,
			CommitThreadService:  servs.commitThread,
			RequestThreadService: servs.requestThread,
			CommentService:       servs.comment,
			ReviewThreadService:  servs.reviewThread},
	}
}

func Init() {
	// read config file
	cfg := config.ReadConfig("./config.json")

	database := db.Connect(&cfg.Database)

	// create environments in order
	repoEnv := initRepoEnv(&cfg, database)
	serviceEnv := initServiceEnv(repoEnv)
	controllerEnv := initControllerEnv(&cfg, serviceEnv)

	// set up routing for the endpoint URLsS
	router := SetUpRouter(&cfg, controllerEnv)
	err := router.Run(fmt.Sprintf("%s:%d", cfg.Hosting.Backend.Host, cfg.Hosting.Backend.Port))
	if err != nil {
		panic(err)
	}
}
