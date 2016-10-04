package trace

import (
	"log"
	"time"

	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"

	// trace
	"github.com/micro/trace-srv/handler"
	"github.com/micro/trace-srv/trace"

	// db
	"github.com/micro/trace-srv/db"
	"github.com/micro/trace-srv/db/mysql"

	gweb "github.com/micro/go-web"
	"github.com/micro/os/internal/helper"
	whandler "github.com/micro/trace-web/handler"

	// proto
	proto "github.com/micro/trace-srv/proto/trace"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.trace"),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	if len(ctx.GlobalString("database_url")) > 0 {
		mysql.Url = ctx.GlobalString("database_url")
	}

	proto.RegisterTraceHandler(service.Server(), new(handler.Trace))

	service.Server().Subscribe(
		service.Server().NewSubscriber(trace.TraceTopic, trace.ProcessSpan),
	)

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func web(ctx *cli.Context) {
	opts := []gweb.Option{
		gweb.Name("go.micro.web.trace"),
		gweb.Handler(whandler.Router()),
	}

	opts = append(opts, helper.WebOpts(ctx)...)

	templateDir := "trace/templates"
	if dir := ctx.GlobalString("html_dir"); len(dir) > 0 {
		templateDir = dir
	}

	whandler.Init(
		templateDir,
		proto.NewTraceClient("go.micro.srv.trace", *cmd.DefaultOptions().Client),
	)

	service := gweb.NewService(opts...)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func traceCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the trace server",
			Action: srv,
		},
		{
			Name:   "web",
			Usage:  "Run the trace web app",
			Action: web,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "trace",
			Usage:       "Trace commands",
			Subcommands: traceCommands(),
		},
	}
}
