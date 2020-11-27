package main

import (
	"flag"

	"github.com/yileCJW/go-zero/core/conf"
	"github.com/yileCJW/go-zero/example/graceful/dns/api/config"
	"github.com/yileCJW/go-zero/example/graceful/dns/api/handler"
	"github.com/yileCJW/go-zero/example/graceful/dns/api/svc"
	"github.com/yileCJW/go-zero/rest"
	"github.com/yileCJW/go-zero/zrpc"
)

var configFile = flag.String("f", "etc/graceful-api.json", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	client := zrpc.MustNewClient(c.Rpc)
	ctx := &svc.ServiceContext{
		Client: client,
	}

	engine := rest.MustNewServer(c.RestConf)
	defer engine.Stop()

	handler.RegisterHandlers(engine, ctx)
	engine.Start()
}
