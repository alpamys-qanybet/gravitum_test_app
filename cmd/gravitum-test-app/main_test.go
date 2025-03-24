package main

import (
	"context"
	"fmt"
	"gravitum-test-app/config"
	"gravitum-test-app/internal/app"
	"gravitum-test-app/internal/repository"
	"gravitum-test-app/internal/repository/postgres"
	"gravitum-test-app/internal/service"
	"gravitum-test-app/pkg/logger"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	once              sync.Once
	cfg               *config.Config
	repos             *repository.Repository
	services          *service.Service
	DbInsertBatchSize int = 10000
)

func setupTestApp(t *testing.T) {
	once.Do(func() {
		gin.SetMode(gin.ReleaseMode)

		projectRoot, err := filepath.Abs("../..") // Adjust to your project structure
		if err != nil {
			fmt.Println("Failed to set working directory:", err)
			os.Exit(1)
		}
		os.Chdir(projectRoot)

		cfgInstance, err := config.New()
		if err != nil {
			t.Fatalf("Failed to load config: %v", err)
		}
		cfg = &cfgInstance
		cfg.Print()

		log := logger.New(logger.GetLevelByString(cfg.Log.Level))

		ctx := context.Background()

		appInstance := app.New(cfg, log)
		err = appInstance.ConnectDB(ctx, cfg.GetDbConfig().GetDsn())
		if err != nil {
			t.Fatalf("couldn't instantiate db: %v", err)
		}

		repos = postgres.NewRepository(cfg, appInstance.Db)
		services = service.NewService(cfg, repos)
	})
}

func handleTestError(t *testing.T, err error) {
	t.Error(err)
}

// go clean -testcache && go test -v -run ^TestScenario$ cmd/gravitum-test-app/main_test.go
func TestScenario(t *testing.T) {
	setupTestApp(t)

	if cfg.App.Profile != "dev" {
		return
	}

	ctx := context.Background()

	list, err := services.User.GetList(ctx)
	if err != nil {
		handleTestError(t, err)
		return
	}

	n := len(list)

	// create
	name := "John"
	surname := "Smith"

	err = services.User.Create(ctx, name, &surname)
	if err != nil {
		handleTestError(t, err)
		return
	}

	list, err = services.User.GetList(ctx)
	if err != nil {
		handleTestError(t, err)
		return
	}

	assert.Equal(t, n+1, len(list), "user creation does not work")

	id := list[len(list)-1].Id

	user, err := services.User.Get(ctx, id)
	if err != nil {
		handleTestError(t, err)
		return
	}

	assert.Equal(t, name, user.Name, "name is incorrect")

	if user.Surname == nil {
		t.Errorf("surname must be nil, it must be=%s", surname)
	} else {
		assert.Equal(t, surname, *user.Surname, "surname is incorrect")
	}

	// update
	newName := "Kevin"
	newSurname := "Tierney"
	err = services.User.Update(ctx, id, newName, &newSurname)
	if err != nil {
		handleTestError(t, err)
		return
	}

	user, err = services.User.Get(ctx, id)
	if err != nil {
		handleTestError(t, err)
		return
	}

	assert.Equal(t, newName, user.Name, "name is incorrect")

	if user.Surname == nil {
		t.Errorf("surname must be nil, it must be=%s", surname)
	} else {
		assert.Equal(t, newSurname, *user.Surname, "surname is incorrect")
	}

}
