package trace

import (
	log "github.com/golang/glog"
	"github.com/micro/cli"
	micro "github.com/micro/go-micro"

	// trace
	"github.com/micro/trace-srv/handler"
	"github.com/micro/trace-srv/trace"

	// db
	"github.com/micro/trace-srv/db"
	"github.com/micro/trace-srv/db/mysql"

	// proto
	proto "github.com/micro/trace-srv/proto/trace"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.trace"),
	)

	if len(ctx.String("database_url")) > 0 {
		mysql.Url = ctx.String("database_url")
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

func traceCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the trace server",
			Action: srv,
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
