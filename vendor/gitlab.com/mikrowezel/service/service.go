package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"sync"

	health "github.com/heptiolabs/healthcheck"
	"gitlab.com/mikrowezel/config"
	logger "gitlab.com/mikrowezel/log"
)

type BaseService struct {
	mux      sync.Mutex
	ctx      context.Context
	cfg      *config.Config
	log      *logger.Logger
	cancel   context.CancelFunc
	name     string
	version  string
	ready    bool
	alive    bool
	handlers map[string]Handler
	health   health.Handler
	worker   Worker
}

// NewService creates a new Service instance.
func NewService(ctx context.Context, cfg *config.Config, log *logger.Logger, cancel context.CancelFunc, name, version string) Service {
	s := &BaseService{
		ctx:      ctx,
		cfg:      cfg,
		log:      log,
		cancel:   cancel,
		name:     name,
		version:  version,
		ready:    false,
		alive:    false,
		handlers: make(map[string]Handler, 1),
	}
	return s
}

// Init service instance.
func (s *BaseService) Init() error {
	// Ensure that service has a worker.
	if s.worker == nil {
		return errors.New("service has no worker")
	}

	// Disable service until everything is setted up.
	s.Disable()

	// Initialize all sevice handlers.
	i := 0
	okChs := make([]chan bool, len(s.handlers))
	for _, h := range s.handlers {
		okChs[i] = h.Init(s)
		i++
	}

	// Check service handlers initialization result.
	allOk := true
	for _, okCh := range okChs {
		allOk = allOk && <-okCh
	}

	if !allOk {
		return fmt.Errorf("cannot initialize '%s' service", s.name)
	}

	// Initialize service worker.
	s.worker.Init()

	// Enable the service
	s.Enable()

	return nil
}

// Lock service mutex.
func (s *BaseService) Lock() {
	s.mux.Lock()
}

// Unlock service mutex.
func (s *BaseService) Unlock() {
	s.mux.Unlock()
}

// Name service.
func (s *BaseService) Name() string {
	return s.name
}

// Version of service.
func (s *BaseService) Version() string {
	return s.version
}

// Enable the service.
func (s *BaseService) Enable() {
	s.ready = true
}

// Disable the service.
func (s *BaseService) Disable() {
	s.ready = false
}

// IsReady tells if service is ready.
func (s *BaseService) IsReady() bool {
	return s.ready
}

// Start service.
func (s *BaseService) Start() {
	go s.checkCancel()
	for _, h := range s.handlers {
		h.Start()
	}
}

// Stop service.
func (s *BaseService) Stop() {
	for _, h := range s.handlers {
		h.Stop()
	}
}

// Handlers available for this service.
func (s *BaseService) Handlers() map[string]Handler {
	return s.handlers
}

// AddHandler for service
func (s *BaseService) AddHandler(h Handler) {
	s.handlers[h.Name()] = h
}

// SetWorker for service
func (s *BaseService) SetWorker(w Worker) {
	w.SetHandlers(s.handlers)
	s.worker = w
}

// Ctx returns service context.
func (s *BaseService) Ctx() context.Context {
	return s.ctx
}

// Cfg returns service configuration.
func (s *BaseService) Cfg() *config.Config {
	return s.cfg
}

// Log for service.
func (s *BaseService) Log() *logger.Logger {
	return s.log
}

func (s *BaseService) checkCancel() {
	<-s.ctx.Done()
	s.Stop()
}

func initChecks(s *BaseService) chan bool {
	ok := make(chan bool)
	go func() {
		defer close(ok)
		s.Lock()
		s.health = health.NewHandler()
		s.health.AddReadinessCheck("ready", s.ReadinessCheck())
		s.health.AddLivenessCheck("heap-threshold", s.HeapLivenessCheck(15))
		s.health.AddLivenessCheck("goroutine-threshold", health.GoroutineCountCheck(100))
		lsa := s.cfg.ValOrDef("liveness.server.address", "0.0.0.0:8086")
		go http.ListenAndServe(lsa, s.health)
		s.Unlock()
		ok <- true
	}()
	return ok
}

// ReadinessCheck lets detect if service is ready to operate.
func (s *BaseService) ReadinessCheck() health.Check {
	return func() error {
		r := s.IsReady()
		if !r {
			msg := fmt.Sprintf("%s is not ready", s.name)
			return errors.New(msg)
		}
		s.log.Info("Readiness check.", "service", s.Name, "status", "ready")
		return nil
	}
}

// HeapLivenessCheck is a heap allocation liveness test for the service.
// After several consecutive and failed attemps of the listener to gather data
// The worker puts itsel in non-ready state.
func (s *BaseService) HeapLivenessCheck(maxMb uint64) health.Check {
	return func() error {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		mb := toMb(m.Alloc)
		r := mb > maxMb
		if r {
			msg := fmt.Sprintf("%s is not in healthy state", s.name)
			return errors.New(msg)
		}
		s.log.Info("Ping liveness check", "service", s.Name, "type", "heap", "status", "alive")
		return nil
	}
}

func toMb(b uint64) uint64 {
	return b / 1024 / 1024

}
