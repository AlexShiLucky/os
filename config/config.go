package config

import (
	"time"

	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/config-srv/config"
	"github.com/micro/config-srv/handler"

	// db
	"github.com/micro/config-srv/db"
	"github.com/micro/config-srv/db/mysql"
	proto "github.com/micro/config-srv/proto/config"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.config"),
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

	proto.RegisterConfigHandler(service.Server(), new(handler.Config))

	// subcriber to watches
	service.Server().Subscribe(service.Server().NewSubscriber(config.WatchTopic, config.Watcher))

	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func configCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the config server",
			Action: srv,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "config",
			Usage:       "Config commands",
			Subcommands: configCommands(),
		},
	}
}
