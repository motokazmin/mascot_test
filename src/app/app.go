package app

import (
	"context"
	"log"
	"mascot/src/config"
	"mascot/src/http"
	"mascot/src/mascot"
	"mascot/src/mngr"
	"sync"
)

type App struct {
	Config  *config.Config
	Http    *http.Service
	Mascot  *mascot.Service
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

	http.AddRouter(mascot.GetMascotRouter())

	app := &App{
		Config:  config,
		Http:    http,
		Mascot:  mascot,
		Context: ctx,
		wg:      &wg,
		cancel:  cancel,
	}
	return app, nil
}

func (a *App) Run() {
	go a.Http.Run()
}

func (a *App) Stop() {
	a.cancel()
	a.wg.Wait()
	log.Println("App stopped")
}
