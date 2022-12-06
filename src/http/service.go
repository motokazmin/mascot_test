package http

import (
	"context"
	"log"
	"mascot/src/config"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Service struct {
	Name   string
	config *config.HttpService
	server *http.Server
	router *mux.Router
	ctx    context.Context
	wg     *sync.WaitGroup
}

func New(ctx context.Context, wg *sync.WaitGroup, conf *config.Config) *Service {
	return &Service{
		Name:   "HttpService",
		config: &conf.HttpService,
		router: mux.NewRouter(),
		ctx:    ctx,
		wg:     wg,
	}
}

func (s *Service) AddRouter(path string, handler http.HandlerFunc, method string) {
	s.router.Handle(path, handler).Methods(method)
}

func (s *Service) Run() {
	s.server = &http.Server{
		Addr:         s.config.Listen,
		Handler:      s.router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s: %v", s.Name, err)
		}
	}()
	s.wg.Add(1)
	log.Printf("%s started", s.Name)

	<-s.ctx.Done()
	s.server.Shutdown(context.Background())
	s.wg.Done()
	log.Printf("%s stopped", s.Name)
}
