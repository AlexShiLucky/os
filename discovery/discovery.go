package discovery

import (
	"log"
	"time"

	"github.com/micro/cli"
	"github.com/micro/discovery-srv/discovery"
	"github.com/micro/discovery-srv/handler"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"

	whandler "github.com/micro/discovery-web/handler"
	gweb "github.com/micro/go-web"
	"github.com/micro/os/internal/helper"

	proto "github.com/micro/discovery-srv/proto/discovery"
	proto2 "github.com/micro/discovery-srv/proto/registry"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.discovery"),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
		micro.BeforeStart(func() error {
			discovery.Run()
			return nil
		}),
	)

	discovery.Init(service)

	service.Server().Subscribe(
		service.Server().NewSubscriber(
			discovery.HeartbeatTopic,
			discovery.DefaultDiscovery.ProcessHeartbeat,
		),
	)

	service.Server().Subscribe(
		service.Server().NewSubscriber(
			discovery.WatchTopic,
			discovery.DefaultDiscovery.ProcessResult,
		),
	)

	proto.RegisterDiscoveryHandler(service.Server(), new(handler.Discovery))
	proto2.RegisterRegistryHandler(service.Server(), new(handler.Registry))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func web(ctx *cli.Context) {
	opts := []gweb.Option{
		gweb.Name("go.micro.web.discovery"),
		gweb.Handler(whandler.Router()),
	}

	opts = append(opts, helper.WebOpts(ctx)...)

	templateDir := "discovery/templates"
	if dir := ctx.GlobalString("html_dir"); len(dir) > 0 {
		templateDir = dir
	}

	whandler.Init(
		templateDir,
		proto.NewDiscoveryClient("go.micro.srv.discovery", *cmd.DefaultOptions().Client),
		proto2.NewRegistryClient("go.micro.srv.discovery", *cmd.DefaultOptions().Client),
	)

	service := gweb.NewService(opts...)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func discoveryCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the discovery server",
			Action: srv,
		},
		{
			Name:   "web",
			Usage:  "Run the discovery web app",
			Action: web,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "discovery",
			Usage:       "Discovery commands",
			Subcommands: discoveryCommands(),
		},
	}
}
