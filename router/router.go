package router

import (
	"log"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"

	"github.com/micro/router-srv/db"
	"github.com/micro/router-srv/db/mysql"
	"github.com/micro/router-srv/handler"
	"github.com/micro/router-srv/router"

	gweb "github.com/micro/go-web"
	"github.com/micro/os/internal/helper"
	whandler "github.com/micro/router-web/handler"

	label "github.com/micro/router-srv/proto/label"
	proto "github.com/micro/router-srv/proto/router"
	rule "github.com/micro/router-srv/proto/rule"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.router"),
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

	router.Init(service)

	proto.RegisterRouterHandler(service.Server(), new(handler.Router))
	label.RegisterLabelHandler(service.Server(), new(handler.Label))
	rule.RegisterRuleHandler(service.Server(), new(handler.Rule))

	// subcriber to stats
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			router.StatsTopic,
			router.ProcessStats,
		),
	)

	// initialise database
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func web(ctx *cli.Context) {
	opts := []gweb.Option{
		gweb.Name("go.micro.web.router"),
		gweb.Handler(whandler.Router()),
	}

	opts = append(opts, helper.WebOpts(ctx)...)

	templateDir := "router/templates"
	if dir := ctx.GlobalString("html_dir"); len(dir) > 0 {
		templateDir = dir
	}

	whandler.Init(
		templateDir,
		label.NewLabelClient("go.micro.srv.router", *cmd.DefaultOptions().Client),
		rule.NewRuleClient("go.micro.srv.router", *cmd.DefaultOptions().Client),
		proto.NewRouterClient("go.micro.srv.router", *cmd.DefaultOptions().Client),
	)

	service := gweb.NewService(opts...)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func routerCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the router server",
			Action: srv,
		},
		{
			Name:   "web",
			Usage:  "Run the router web app",
			Action: web,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "router",
			Usage:       "Router commands",
			Subcommands: routerCommands(),
		},
	}
}
