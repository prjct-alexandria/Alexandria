package server

import (
	"database/sql"
	"flag"
	"fmt"
	"mainServer/controllers"
	"mainServer/db"
	"mainServer/repositories/interfaces"
	"mainServer/repositories/postgres"
	"mainServer/repositories/storer"
	"mainServer/server/config"
	"mainServer/services"
	servinterfaces "mainServer/services/interfaces"
)

type RepoEnv struct {
	storer                interfaces.Storer
	article               interfaces.ArticleRepository
	user                  interfaces.UserRepository
	version               interfaces.VersionRepository
	req                   interfaces.RequestRepository
	thread                interfaces.ThreadRepository
	comment               interfaces.CommentRepository
	commitThread          interfaces.CommitThreadRepository
	requestThread         interfaces.RequestThreadRepository
	reviewThread          interfaces.ReviewThreadRepository
	commitSelectionThread interfaces.CommitSelectionThreadRepository
}

type ServiceEnv struct {
	version               servinterfaces.VersionService
	article               servinterfaces.ArticleService
	user                  servinterfaces.UserService
	req                   servinterfaces.RequestService
	thread                servinterfaces.ThreadService
	comment               servinterfaces.CommentService
	commitThread          servinterfaces.CommitThreadService
	commitSelectionThread servinterfaces.CommitSelectionThreadService
	requestThread         servinterfaces.RequestThreadService
	reviewThread          servinterfaces.ReviewThreadService
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
		storer:                storer.NewStorer(&cfg.Fs),
		article:               postgres.NewPgArticleRepository(database),
		user:                  postgres.NewPgUserRepository(database),
		version:               postgres.NewPgVersionRepository(database),
		req:                   postgres.NewPgRequestRepository(database),
		thread:                postgres.NewPgThreadRepository(database),
		comment:               postgres.NewPgCommentRepository(database),
		commitThread:          postgres.NewPgCommitThreadRepository(database),
		commitSelectionThread: postgres.NewPgCommitSelectionThreadRepository(database),
		requestThread:         postgres.NewPgRequestThreadRepository(database),
		reviewThread:          postgres.NewPgReviewThreadRepository(database),
	}
}

func initServiceEnv(repos RepoEnv) ServiceEnv {
	return ServiceEnv{
		article:               services.NewArticleService(repos.article, repos.version, repos.user, repos.storer),
		user:                  services.NewUserService(repos.user),
		req:                   services.RequestService{Repo: repos.req, Versionrepo: repos.version, Storer: repos.storer},
		version:               services.VersionService{VersionRepo: repos.version, Storer: repos.storer, UserRepo: repos.user},
		thread:                services.ThreadService{ThreadRepository: repos.thread},
		comment:               services.CommentService{CommentRepository: repos.comment},
		commitThread:          services.CommitThreadService{CommitThreadRepository: repos.commitThread},
		commitSelectionThread: services.CommitSelectionThreadService{CommitSelectionThreadRepository: repos.commitSelectionThread},
		requestThread: services.RequestThreadService{RequestThreadRepository: repos.requestThread,
			VersionRepository: repos.version,
			RequestRepository: repos.req},
		reviewThread: services.ReviewThreadService{ReviewThreadRepository: repos.reviewThread},
	}
}

func initControllerEnv(cfg *config.Config, servs ServiceEnv) ControllerEnv {
	return ControllerEnv{
		article: controllers.NewArticleController(servs.article),
		user:    controllers.UserController{UserService: servs.user, Cfg: cfg},
		req:     controllers.RequestController{Serv: servs.req},
		version: controllers.VersionController{Serv: servs.version},
		thread: controllers.ThreadController{ThreadService: servs.thread,
			CommitThreadService:          servs.commitThread,
			CommitSelectionThreadService: servs.commitSelectionThread,
			RequestThreadService:         servs.requestThread,
			CommentService:               servs.comment,
			ReviewThreadService:          servs.reviewThread},
	}
}

func Init() {
	// read config file
	var cfg config.Config

	// check if running in a docker environment
	dockerPtr := flag.Bool("dockerconfig", false, "running in docker environment")
	flag.Parse()

	if *dockerPtr {
		cfg = config.ReadConfig("./dockerconfig.json")
	} else {
		cfg = config.ReadConfig("./config.json")
	}

	// connect to database
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
