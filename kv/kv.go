package kv

import (
	"log"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-platform/kv"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.kv"),
		micro.Version("latest"),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
		micro.Flags(),
	)

	opts := []kv.Option{
		kv.Client(service.Client()),
		kv.Server(service.Server()),
		kv.Internal(false),
	}

	if len(ctx.String("namespace")) > 0 {
		opts = append(opts, kv.Namespace(ctx.String("namespace")))
	}

	keyval := kv.NewKV(opts...)
	defer keyval.Close()

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func kvCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the kv server",
			Action: srv,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "namespace",
					EnvVar: "NAMESPACE",
					Usage:  "Namespace for the key-value store",
				},
			},
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "kv",
			Usage:       "KV commands",
			Subcommands: kvCommands(),
		},
	}
}
