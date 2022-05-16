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
	version    controllers.VersionController
	helloWorld controllers.HelloWorldController
}

func initRepoEnv() (RepoEnv, error) {
	gitrepo := repositories.GitRepository{Path: ".../gitfiles"}

	env := RepoEnv{
		git: gitrepo,
	}
	return env, nil
}

func initServiceEnv() (ServiceEnv, error) {
	repos, err := initRepoEnv()
	if err != nil {
		return ServiceEnv{}, err
	}

	env := ServiceEnv{
		version: services.VersionService{Gitrepo: repos.git},
	}

	return env, nil
}

func initControllerEnv() (ControllerEnv, error) {
	servs, err := initServiceEnv()
	if err != nil {
		return ControllerEnv{}, err
	}

	env := ControllerEnv{
		version: controllers.VersionController{Serv: servs.version},
	}

	return env, nil
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
