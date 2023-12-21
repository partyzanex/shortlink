package service

import (
	"sync"

	"github.com/partyzanex/shutdown"
)

type Service struct {
	status    Status
	debugHost string

	mu sync.RWMutex

	init  sync.Once
	start sync.Once

	closer *shutdown.Lifo
}

func New() *Service {
	return &Service{
		status:    0,
		debugHost: "",
		mu:        sync.RWMutex{},
		init:      sync.Once{},
		start:     sync.Once{},
		closer:    new(shutdown.Lifo),
	}
}

func (s *Service) Close() error {
	return s.closer.Close()
}

func (s *Service) GetStatus() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.status
}

func (s *Service) setStatus(status Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status = status
}

func (s *Service) setDebugHost(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.debugHost = host
}
