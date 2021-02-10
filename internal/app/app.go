package app

import (
	"github.com/TheRTK/http-multiplexer/internal/config"
	"github.com/TheRTK/http-multiplexer/internal/services/request"
)

type App struct {
	cfg *config.Config
	requestService request.IRequestService
}

func New(opts ...Option) *App {
	app := &App{}

	for _, option := range opts {
		option(app)
	}

	return app
}

func (a *App) GetRequestService() request.IRequestService  {
	return a.requestService
}

