package discovery

import (
	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/discovery-srv/discovery"
	"github.com/micro/discovery-srv/handler"

	proto "github.com/micro/discovery-srv/proto/discovery"
	proto2 "github.com/micro/discovery-srv/proto/registry"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.discovery"),
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

func discoveryCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the discovery server",
			Action: srv,
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
