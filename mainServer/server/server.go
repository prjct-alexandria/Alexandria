package server

import (
	"fmt"
	"mainServer/controllers"
	"mainServer/db"
	"mainServer/repositories"
	"mainServer/repositories/interfaces"
	"mainServer/repositories/postgres"
	"mainServer/services"
)

type RepoEnv struct {
	git     repositories.GitRepository
	article interfaces.ArticleRepository
	user    interfaces.UserRepository
	version interfaces.VersionRepository
}

type ServiceEnv struct {
	article services.ArticleService
	user    services.UserService
	version services.VersionService
}

type ControllerEnv struct {
	article controllers.ArticleController
	version controllers.VersionController
	user    controllers.UserController
}

func initRepoEnv() (RepoEnv, error) {
	// TODO: gitfiles path in config file
	gitpath := "../../gitfiles"

	gitrepo, err := repositories.NewGitRepository(gitpath)
	if err != nil {
		return RepoEnv{}, err
	}

	database := db.Connect()

	return RepoEnv{
		git:     gitrepo,
		article: postgres.NewPgArticleRepository(database),
		user:    postgres.NewPgUserRepository(database),
		version: postgres.NewPgVersionRepository(database),
	}, nil
}

func initServiceEnv() (ServiceEnv, error) {
	repos, err := initRepoEnv()
	if err != nil {
		return ServiceEnv{}, err
	}

	return ServiceEnv{
		article: services.NewArticleService(repos.article, repos.version, repos.git),
		user:    services.UserService{UserRepository: repos.user},
		version: services.VersionService{Gitrepo: repos.git, Versionrepo: repos.version},
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
		version: controllers.VersionController{Serv: servs.version},
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
