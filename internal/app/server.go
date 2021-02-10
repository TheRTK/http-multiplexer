package app

import (
	"github.com/TheRTK/http-multiplexer/internal/config"
	"github.com/TheRTK/http-multiplexer/internal/services/request"
)

type Server interface {
	SetRequestService(rs request.IRequestService)
	GetRequestService() request.IRequestService
	SetConfig(cfg *config.Config)
	GetConfig() *config.Config
}
