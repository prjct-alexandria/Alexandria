package server

import (
	"fmt"
	"mainServer/controllers"
	"mainServer/db"
	"mainServer/repositories"
	"mainServer/repositories/interfaces"
	"mainServer/repositories/postgres"
	"mainServer/services"
	"os"
)

type RepoEnv struct {
	git           repositories.GitRepository
	user          interfaces.UserRepository
	thread        interfaces.ThreadRepository
	comment       interfaces.CommentRepository
	commitThread  interfaces.CommitThreadRepository
	requestThread interfaces.RequestThreadRepository
	reviewThread  interfaces.ReviewThreadRepository
}

type ServiceEnv struct {
	version       services.VersionService
	user          services.UserService
	thread        services.ThreadService
	comment       services.CommentService
	commitThread  services.CommitThreadService
	requestThread services.RequestThreadService
	reviewThread  services.ReviewThreadService
}

type ControllerEnv struct {
	version controllers.VersionController
	user    controllers.UserController
	thread  controllers.ThreadController
}

func initRepoEnv() (RepoEnv, error) {
	// TODO: gitfiles path in config file
	gitpath := "../../gitfiles"

	// make folder for git files
	err := os.MkdirAll(gitpath, os.ModePerm)
	if err != nil {
		return RepoEnv{}, err
	}

	database := db.Connect()

	return RepoEnv{
		git:           repositories.NewGitRepository(gitpath),
		user:          postgres.NewPgUserRepository(database),
		thread:        postgres.NewPgThreadRepository(database),
		comment:       postgres.NewPgCommentRepository(database),
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
		version:       services.VersionService{Gitrepo: repos.git},
		user:          services.UserService{UserRepository: repos.user},
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
		version: controllers.VersionController{Serv: servs.version},
		user:    controllers.UserController{UserService: servs.user},
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
