package main

import (
	"log"

	"go.uber.org/dig"

	"github.com/EQUISEED-WEALTH/pinch/backend/controller"
	"github.com/EQUISEED-WEALTH/pinch/backend/domain/utils"
	"github.com/EQUISEED-WEALTH/pinch/backend/repo"
	"github.com/EQUISEED-WEALTH/pinch/backend/service"
)

// Handle Build Error
func handleBuildError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Provide all build functions
func provide(container *dig.Container, buildFuncs ...utils.Any) {
	for _, buildFunc := range buildFuncs {
		err := container.Provide(buildFunc)
		handleBuildError(err)
	}
}

// Invoke all build functions
func invoke(container *dig.Container, buildFuncs ...utils.Any) {
	for _, buildFunc := range buildFuncs {
		err := container.Invoke(buildFunc)
		handleBuildError(err)
	}
}

// Build all dependencies
func buildContainer() *dig.Container {
	container := dig.New()

	provide(container,
		controller.BuildGinEngine,
		repo.NewPgDB,
		repo.NewRedisClient,

		// Repositories
		repo.NewUserRepo,
		repo.NewAugmontUserRepo,
		repo.NewAugmontOrderRepo,
		repo.NewAugmontInMemRepo,

		// Services
		service.NewAugmondService,
		service.NewUserService,
	)

	invoke(container,
		// Controllers
		controller.NewUserController,
		controller.NewGoldController,
	)

	return container
}
