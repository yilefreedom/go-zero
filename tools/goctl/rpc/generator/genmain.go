package generator

import (
	"fmt"
	"path/filepath"
	"strings"

	conf "github.com/yilefreedom/go-zero/tools/goctl/config"
	"github.com/yilefreedom/go-zero/tools/goctl/rpc/parser"
	"github.com/yilefreedom/go-zero/tools/goctl/util"
	"github.com/yilefreedom/go-zero/tools/goctl/util/format"
	"github.com/yilefreedom/go-zero/tools/goctl/util/stringx"
)

const mainTemplate = `{{.head}}

package main

import (
	"flag"
	"fmt"

	{{.imports}}

	"github.com/yilefreedom/go-zero/core/conf"
	"github.com/yilefreedom/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/{{.serviceName}}.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	srv := server.New{{.serviceNew}}Server(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		{{.pkg}}.Register{{.service}}Server(grpcServer, srv)
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
`

func (g *defaultGenerator) GenMain(ctx DirContext, proto parser.Proto, cfg *conf.Config) error {
	dir := ctx.GetMain()
	mainFilename, err := format.FileNamingFormat(cfg.NamingFormat, ctx.GetMain().Base)
	if err != nil {
		return err
	}

	fileName := filepath.Join(dir.Filename, fmt.Sprintf("%v.go", mainFilename))
	imports := make([]string, 0)
	pbImport := fmt.Sprintf(`"%v"`, ctx.GetPb().Package)
	svcImport := fmt.Sprintf(`"%v"`, ctx.GetSvc().Package)
	remoteImport := fmt.Sprintf(`"%v"`, ctx.GetServer().Package)
	configImport := fmt.Sprintf(`"%v"`, ctx.GetConfig().Package)
	imports = append(imports, configImport, pbImport, remoteImport, svcImport)
	head := util.GetHead(proto.Name)
	text, err := util.LoadTemplate(category, mainTemplateFile, mainTemplate)
	if err != nil {
		return err
	}

	return util.With("main").GoFmt(true).Parse(text).SaveTo(map[string]interface{}{
		"head":        head,
		"serviceName": strings.ToLower(stringx.From(ctx.GetMain().Base).ToCamel()),
		"imports":     strings.Join(imports, util.NL),
		"pkg":         proto.PbPackage,
		"serviceNew":  stringx.From(proto.Service.Name).ToCamel(),
		"service":     parser.CamelCase(proto.Service.Name),
	}, fileName, false)
}
