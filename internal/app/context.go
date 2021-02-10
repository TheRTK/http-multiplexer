package app

import (
	"context"

	"github.com/TheRTK/http-multiplexer/internal/services/request"
)

type Context struct {
	App    IAppFace
	Ctx    context.Context
}

type IAppFace interface {
	GetRequestService() request.IRequestService
}

func NewContext(ctx context.Context, app IAppFace) *Context {
	return &Context{
		App: app,
		Ctx: ctx,
	}
}

