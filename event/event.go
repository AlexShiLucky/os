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
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
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

func eventCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the event server",
			Action: srv,
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
