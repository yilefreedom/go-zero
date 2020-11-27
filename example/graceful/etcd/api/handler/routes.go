// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"github.com/yilefreedom/go-zero/example/graceful/etcd/api/svc"
	"github.com/yilefreedom/go-zero/rest"
)

func RegisterHandlers(engine *rest.Server, ctx *svc.ServiceContext) {
	engine.AddRoutes([]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/api/graceful",
			Handler: gracefulHandler(ctx),
		},
	})
}
