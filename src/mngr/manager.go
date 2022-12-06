package mngr

import (
	"context"
	"mascot/src/config"
	"mascot/src/http"
	"mascot/src/mascot"
	"sync"
)

type Manager struct {
	once   map[string]*sync.Once
	conf   *config.Config
	ctx    context.Context
	wg     *sync.WaitGroup
	http   *http.Service
	mascot *mascot.Service
}

func New(ctx context.Context, wg *sync.WaitGroup, conf *config.Config) *Manager {
	return &Manager{
		once: map[string]*sync.Once{
			"http":   {},
			"mascot": {},
		},
		conf: conf,
		ctx:  ctx,
		wg:   wg,
	}
}

func (m *Manager) GetHttpService() *http.Service {
	m.once["http"].Do(func() {
		m.http = http.New(m.ctx, m.wg, m.conf)
	})
	return m.http
}

func (m *Manager) GetMascotService() *mascot.Service {
	m.once["mascot"].Do(func() {
		m.mascot = mascot.New(m.ctx, m.wg, m.conf)
	})
	return m.mascot
}
