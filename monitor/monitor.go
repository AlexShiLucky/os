package monitor

import (
	"time"

	log "github.com/golang/glog"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/monitor-srv/handler"
	"github.com/micro/monitor-srv/monitor"

	gweb "github.com/micro/go-web"
	proto "github.com/micro/monitor-srv/proto/monitor"
	whandler "github.com/micro/monitor-web/handler"
	"github.com/micro/platform/internal/helper"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.monitor"),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
		micro.BeforeStart(func() error {
			monitor.DefaultMonitor.Run()
			return nil
		}),
	)

	// healthchecks
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			monitor.HealthCheckTopic,
			monitor.DefaultMonitor.ProcessHealthCheck,
		),
	)

	// status
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			monitor.StatusTopic,
			monitor.DefaultMonitor.ProcessStatus,
		),
	)

	// stats
	service.Server().Subscribe(
		service.Server().NewSubscriber(
			monitor.StatsTopic,
			monitor.DefaultMonitor.ProcessStats,
		),
	)

	proto.RegisterMonitorHandler(service.Server(), new(handler.Monitor))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func web(ctx *cli.Context) {
	opts := []gweb.Option{
		gweb.Name("go.micro.web.monitor"),
		gweb.Handler(whandler.Router()),
	}

	opts = append(opts, helper.WebOpts(ctx)...)

	templateDir := "monitor/templates"
	if dir := ctx.GlobalString("html_dir"); len(dir) > 0 {
		templateDir = dir
	}

	whandler.Init(
		templateDir,
		proto.NewMonitorClient("go.micro.srv.monitor", *cmd.DefaultOptions().Client),
	)

	service := gweb.NewService(opts...)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func monitorCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the monitor server",
			Action: srv,
		},
		{
			Name:   "web",
			Usage:  "Run the monitor web app",
			Action: web,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "monitor",
			Usage:       "Monitor commands",
			Subcommands: monitorCommands(),
		},
	}
}
