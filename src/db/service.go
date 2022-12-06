package db

import (
	"context"
	"database/sql"
	"log"
	"mascot/src/config"
	"sync"

	_ "github.com/lib/pq"
)

type Service struct {
	Name   string
	config *config.PostgresService
	client *sql.DB
	ctx    context.Context
	wg     *sync.WaitGroup
}

func New(ctx context.Context, wg *sync.WaitGroup, config *config.Config) *Service {
	return &Service{
		Name:   "PostgresService",
		config: &config.PostgresService,
		ctx:    ctx,
		wg:     wg,
	}
}

func (s *Service) Run() {
	connStr := "user=" + s.config.User + " dbname=" + s.config.Dbname +
		" sslmode=" + s.config.SslMode + " password=" + s.config.Password

	client, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping()
	if err != nil {
		log.Fatalf("could not connect to postgres after ping: %v", err)
	}

	s.client = client
	s.wg.Add(1)
	log.Printf("%s started", s.Name)

	<-s.ctx.Done()
	s.Close()
	s.wg.Done()
	log.Printf("%s stopped", s.Name)
}

func (s *Service) Close() {
	s.client.Close()
}

/*
func (s *Service) Select(space string, index string, offset, limit uint64, key []interface{}) ([][]interface{}, error) {
	resp, err := s.client.Select(space, index, uint32(offset), uint32(limit), tarantool.IterEq, key)
	if err != nil {
		return nil, err
	}
	return resp.Tuples(), nil
}*/
