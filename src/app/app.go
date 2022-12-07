// Обеспечивает инициализацию и запуск/остановку сервисов приложения.
package app

import (
	"context"
	"log"
	"mascot/src/config"
	"mascot/src/db"
	"mascot/src/http"
	"mascot/src/mngr"
	"sync"
)

type App struct {
	Config  *config.Config
	Db      *db.Service
	Http    *http.Service
	Context context.Context
	wg      *sync.WaitGroup
	cancel  context.CancelFunc
}

func New(config *config.Config) (*App, error) {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	manager := mngr.New(ctx, &wg, config)
	http := manager.GetHttpService()
	mascot := manager.GetMascotService()
	db := manager.GetPostgresService()

	http.AddRouter(mascot.GetMascotRouter())

	app := &App{
		Config:  config,
		Db:      db,
		Http:    http,
		Context: ctx,
		wg:      &wg,
		cancel:  cancel,
	}
	return app, nil
}

func (a *App) Run() {
	go a.Db.Run()
	go a.Http.Run()
}

func (a *App) Stop() {
	a.cancel()
	a.wg.Wait()
	log.Println("App stopped")
}
