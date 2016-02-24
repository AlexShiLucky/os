package event

import (
	"time"

	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/event-srv/db"
	"github.com/micro/event-srv/db/mysql"
	"github.com/micro/event-srv/event"
	"github.com/micro/event-srv/handler"
	proto "github.com/micro/event-srv/proto/event"
	whandler "github.com/micro/event-web/handler"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/server"
	gweb "github.com/micro/go-web"
	"github.com/micro/platform/internal/helper"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.event"),
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

	proto.RegisterEventHandler(service.Server(), new(handler.Event))

	service.Server().Subscribe(
		service.Server().NewSubscriber(
			"micro.event.record",
			event.Process,
			server.SubscriberQueue("event-srv"),
		),
	)

	// For watchers
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			"micro.event.record",
			event.Stream,
		),
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
		gweb.Name("go.micro.web.event"),
		gweb.Handler(whandler.Router()),
	}

	opts = append(opts, helper.WebOpts(ctx)...)

	templateDir := "event/templates"
	if dir := ctx.GlobalString("html_dir"); len(dir) > 0 {
		templateDir = dir
	}

	whandler.Init(
		templateDir,
		proto.NewEventClient("go.micro.srv.event", *cmd.DefaultOptions().Client),
	)

	service := gweb.NewService(opts...)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func eventCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the event server",
			Action: srv,
		},
		{
			Name:   "web",
			Usage:  "Run the event web app",
			Action: web,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "event",
			Usage:       "Event commands",
			Subcommands: eventCommands(),
		},
	}
}
