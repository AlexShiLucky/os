package db

import (
	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/db-srv/db"
	_ "github.com/micro/db-srv/db/mysql"
	"github.com/micro/db-srv/handler"

	proto "github.com/micro/db-srv/proto/db"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.db"),
	)

	if len(ctx.String("database_service_namespace")) > 0 {
		db.DBServiceNamespace = ctx.String("database_service_namespace")
	}

	proto.RegisterDBHandler(service.Server(), new(handler.DB))

	if err := db.Init(service.Client().Options().Selector); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func dbCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the db server",
			Action: srv,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "database_service_namespace",
					EnvVar: "DATABASE_SERVICE_NAMESPACE",
					Usage:  "The namespace used when looking up databases in registry e.g go.micro.db",
				},
			},
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "db",
			Usage:       "DB commands",
			Subcommands: dbCommands(),
		},
	}
}
