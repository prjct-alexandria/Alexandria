package server

import (
	"fmt"
	"mainServer/controllers"
	"mainServer/db"
	"mainServer/repositories"
	"mainServer/repositories/interfaces"
	"mainServer/repositories/postgres"
	"mainServer/services"
	servinterfaces "mainServer/services/interfaces"
	"mainServer/utils/clock"
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

func initRepoEnv() (RepoEnv, error) {
	// TODO: gitfiles path in config file
	gitpath := "../../gitfiles"

	gitrepo, err := repositories.NewGitRepository(gitpath)
	if err != nil {
		return RepoEnv{}, err
	}

	database := db.Connect()
	clock := clock.RealClock{}

	return RepoEnv{
		git:           gitrepo,
		article:       postgres.NewPgArticleRepository(database),
		user:          postgres.NewPgUserRepository(database),
		version:       postgres.NewPgVersionRepository(database),
		req:           postgres.NewPgRequestRepository(database),
		thread:        postgres.NewPgThreadRepository(database),
		comment:       postgres.NewPgCommentRepository(database, clock),
		commitThread:  postgres.NewPgCommitThreadRepository(database),
		requestThread: postgres.NewPgRequestThreadRepository(database),
		reviewThread:  postgres.NewPgReviewThreadRepository(database),
	}, nil
}

func initServiceEnv() (ServiceEnv, error) {
	repos, err := initRepoEnv()
	if err != nil {
		return ServiceEnv{}, err
	}

	return ServiceEnv{
		article:       services.NewArticleService(repos.article, repos.version, repos.git),
		user:          services.UserService{UserRepository: repos.user},
		req:           services.RequestService{Repo: repos.req, Versionrepo: repos.version},
		version:       services.VersionService{Gitrepo: repos.git, Versionrepo: repos.version},
		thread:        services.ThreadService{ThreadRepository: repos.thread},
		comment:       services.CommentService{CommentRepository: repos.comment},
		commitThread:  services.CommitThreadService{CommitThreadRepository: repos.commitThread},
		requestThread: services.RequestThreadService{RequestThreadRepository: repos.requestThread},
		reviewThread:  services.ReviewThreadService{ReviewThreadRepository: repos.reviewThread},
	}, nil
}

func initControllerEnv() (ControllerEnv, error) {
	servs, err := initServiceEnv()
	if err != nil {
		return ControllerEnv{}, err
	}

	return ControllerEnv{
		article: controllers.NewArticleController(servs.article),
		user:    controllers.UserController{UserService: servs.user},
		req:     controllers.RequestController{Serv: servs.req},
		version: controllers.VersionController{Serv: servs.version},
		thread: controllers.ThreadController{ThreadService: servs.thread,
			CommitThreadService:  servs.commitThread,
			RequestThreadService: servs.requestThread,
			CommentService:       servs.comment,
			ReviewThreadService:  servs.reviewThread},
	}, nil
}

func Init() {
	env, err := initControllerEnv()
	if err != nil {
		fmt.Println(err)
		return
	}

	router := SetUpRouter(env)
	err = router.Run("localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
}
