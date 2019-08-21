package service

import (
	"context"

	"gitlab.com/mikrowezel/config"
	logger "gitlab.com/mikrowezel/log"
)

type (
	// Handler is a DB handler.
	BaseHandler struct {
		ctx   context.Context
		cfg   *config.Config
		log   *logger.Logger
		name  string
		ready bool
		alive bool
	}
)

// NewBaseHandler creates a new base handler.
func NewBaseHandler(ctx context.Context, cfg *config.Config, log *logger.Logger, name string) *BaseHandler {

	return &BaseHandler{
		ctx:   ctx,
		cfg:   cfg,
		log:   log,
		name:  name,
		ready: false,
		alive: false,
	}
}

// Init a new handler.
func (h *BaseHandler) Init(s Service) chan bool {
	return make(chan bool)
}

// Name returns the server name.
func (h *BaseHandler) Name() string {
	return h.name
}

// Enable handler.
func (h *BaseHandler) Enable() {
	h.ready = true
}

// Disable handler.
func (h *BaseHandler) Disable() {
	h.ready = false
}

// IsReady returns the current state of handler.
func (h *BaseHandler) IsReady() bool {
	return h.ready
}

// Start handler
func (h *BaseHandler) Start() error {
	return nil
}

// Stop handler
func (h *BaseHandler) Stop() {
	return
}

// Ctx returns service context.
func (h *BaseHandler) Ctx() context.Context {
	return h.ctx
}

// Cfg returns service config.
func (h *BaseHandler) Cfg() *config.Config {
	return h.cfg
}

// Log returns current service logger.
func (h *BaseHandler) Log() *logger.Logger {
	return h.log
}
