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
	git  repositories.GitRepository
	user interfaces.UserRepository
	req  interfaces.RequestRepository
}

type ServiceEnv struct {
	version services.VersionService
	user    services.UserService
	req     services.RequestService
}

type ControllerEnv struct {
	version controllers.VersionController
	user    controllers.UserController
	req     controllers.RequestController
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
		git:  repositories.NewGitRepository(gitpath),
		user: postgres.NewPgUserRepository(database),
		req:  postgres.NewPgRequestRepository(database),
	}, nil
}

func initServiceEnv() (ServiceEnv, error) {
	repos, err := initRepoEnv()
	if err != nil {
		return ServiceEnv{}, err
	}

	return ServiceEnv{
		version: services.VersionService{Gitrepo: repos.git},
		user:    services.UserService{UserRepository: repos.user},
		req:     services.RequestService{Repo: repos.req},
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
		req:     controllers.RequestController{Serv: servs.req},
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
