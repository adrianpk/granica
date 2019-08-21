package service

import (
	"context"

	"gitlab.com/mikrowezel/config"
	logger "gitlab.com/mikrowezel/log"
)

type (
	// Service defines a service interface
	Service interface {
		// Lock mutex.
		Lock()

		// Unlock mutex.
		Unlock()

		// Init service instance.
		Init() error

		// Name service.
		Name() string

		// Version of service.
		Version() string

		// Enable the service.
		Enable()

		// Disable the service.
		Disable()

		// IsReady returns true is service is ready.
		IsReady() bool

		// Start service.
		Start()

		// Stop service.
		Stop()

		// Handlers available for this service.
		Handlers() map[string]Handler

		// AddHandler for service.
		AddHandler(h Handler)

		// SetWorker for service
		SetWorker(w Worker)

		// Ctx returns service context.
		Ctx() context.Context

		// Cfg returns service config.
		Cfg() *config.Config

		// Log returns current service logger.
		Log() *logger.Logger
	}

	// Handler defines a service handler interface
	Handler interface {
		// Init a new handler.
		Init(s Service) chan bool

		// Name returns the server name.
		Name() string

		// Enable handler.
		Enable()

		// Disable handler.
		Disable()

		// IsReady returns the current state of handler.
		IsReady() bool

		// Start handler
		Start() error

		// Stop handler
		Stop()

		// Ctx returns handler context.
		Ctx() context.Context

		// Ctx returns handler context.
		Cfg() *config.Config

		// Log returns current handler logger.
		Log() *logger.Logger
	}

	// Worket defines a service worker interface
	Worker interface {
		// Init a new worker.
		Init() bool

		// Name returns the server name.
		Name() string

		// Enable handler.
		Enable()

		// Disable handler.
		Disable()

		// IsReady returns the current state of handler.
		IsReady() bool

		// Start handler.
		Start() error

		// Stop handler.
		Stop()

		// AttachTo service.
		AttachTo(svc Service)

		// SetHandlers for worker.
		SetHandlers(handlers map[string]Handler)

		// Handler returns worker handler by name.
		Handler(name string) (h Handler, ok bool)

		// Ctx returns worker context.
		Ctx() context.Context

		// Cfg returns worker config.
		Cfg() *config.Config

		// Log returns current worker logger.
		Log() *logger.Logger
	}
)
