package service

import (
	"context"

	"gitlab.com/mikrowezel/config"
	logger "gitlab.com/mikrowezel/log"
)

type BaseWorker struct {
	ctx      context.Context
	cfg      *config.Config
	log      *logger.Logger
	cancel   context.CancelFunc
	name     string
	ready    bool
	handlers map[string]Handler
}

// NewWorker creates a new base worker instance.
// This is a bare implementtion of Worker interface
// just for mocking and/or testing purposes.
func NewWorker(ctx context.Context, cfg *config.Config, log *logger.Logger, name string) *BaseWorker {
	w := &BaseWorker{
		ctx:   ctx,
		cfg:   cfg,
		log:   log,
		name:  name,
		ready: false,
	}
	return w
}

// Init a new worker.
func (w *BaseWorker) Init() bool {
	return true
}

// Name returns the server name.
func (w *BaseWorker) Name() string {
	return w.name
}

// Enable handler.
func (w *BaseWorker) Enable() {
	w.ready = true
}

// Disable handler.
func (w *BaseWorker) Disable() {
	w.ready = false
}

// IsReady returns the current state of handler.
func (w *BaseWorker) IsReady() bool {
	return w.ready
}

// Start handler.
func (w *BaseWorker) Start() error {
	return nil
}

// Stop handler.
func (w *BaseWorker) Stop() {
	return
}

// AttachTo service.
func (w *BaseWorker) AttachTo(s Service) {
	s.SetWorker(w)
}

// SetHandlers for worker.
func (w *BaseWorker) SetHandlers(handlers map[string]Handler) {
	w.handlers = handlers
}

// Handler returns worker handler by name.
func (w *BaseWorker) Handler(name string) (h Handler, ok bool) {
	if w.handlers == nil {
		return nil, false
	}
	h, ok = w.handlers[name]
	return h, ok
}

// Ctx returns service context.
func (w *BaseWorker) Ctx() context.Context {
	return w.ctx
}

// Cfg returns service configuration.
func (w *BaseWorker) Cfg() *config.Config {
	return w.cfg
}

// Log for worker.
func (w *BaseWorker) Log() *logger.Logger {
	return w.log
}
