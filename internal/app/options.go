package app

import (
	"github.com/TheRTK/http-multiplexer/internal/config"
	"github.com/TheRTK/http-multiplexer/internal/services/request"
)

type (
	OptionCreator func() []Option
	Option        func(a *App)
	ServerOption  func(s Server)
)

func ServerConnector(s Server) Option {
	return func(a *App) {
		a.requestService = s.GetRequestService()
		a.cfg = s.GetConfig()
	}
}

func ConfigServer(
	cfg *config.Config,
	requestService request.IRequestService,
) ServerOption {
	return func(s Server) {
		s.SetRequestService(requestService)
		s.SetConfig(cfg)
	}
}
