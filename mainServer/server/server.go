package server

import (
	"mainServer/controllers"
	"mainServer/repositories"
	"mainServer/services"
)

type RepoEnv struct {
	git repositories.GitRepository
}

type ServiceEnv struct {
	version services.VersionService
}

type ControllerEnv struct {
	version controllers.VersionController
}

func initRepoEnv() (RepoEnv, error) {
	return RepoEnv{
		git: repositories.GitRepository{Path: ".../gitfiles"},
	}, nil
}

func initServiceEnv() (ServiceEnv, error) {
	repos, err := initRepoEnv()
	if err != nil {
		return ServiceEnv{}, err
	}

	return ServiceEnv{
		version: services.VersionService{Gitrepo: repos.git},
	}, nil
}

func initControllerEnv() (ControllerEnv, error) {
	servs, err := initServiceEnv()
	if err != nil {
		return ControllerEnv{}, err
	}

	return ControllerEnv{
		version: controllers.VersionController{Serv: servs.version},
	}, nil
}

func Init() {
	env, err := initControllerEnv()
	if err != nil {
		return
	}

	router := SetUpRouter(env)
	err = router.Run("localhost:8080")
	if err != nil {
		return
	}
}
