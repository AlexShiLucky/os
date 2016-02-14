package router

import (
	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/router-srv/db"
	"github.com/micro/router-srv/db/mysql"
	"github.com/micro/router-srv/handler"
	"github.com/micro/router-srv/router"

	label "github.com/micro/router-srv/proto/label"
	proto "github.com/micro/router-srv/proto/router"
	rule "github.com/micro/router-srv/proto/rule"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.router"),
	)

	if len(ctx.String("database_url")) > 0 {
		mysql.Url = ctx.String("database_url")
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

func routerCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the router server",
			Action: srv,
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
