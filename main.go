package main

import (
	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/platform/auth"
	"github.com/micro/platform/config"
	"github.com/micro/platform/db"
	"github.com/micro/platform/discovery"
	"github.com/micro/platform/event"
	//	"github.com/micro/platform/kv"
	//	"github.com/micro/platform/log"
	//	"github.com/micro/platform/metrics"
	"github.com/micro/platform/monitor"
	"github.com/micro/platform/router"
	//	"github.com/micro/platform/sync"
	"github.com/micro/platform/trace"
)

func setup(app *ccli.App) {
	// common flags
	app.Flags = append(app.Flags,
		ccli.StringFlag{
			Name:   "database_url",
			EnvVar: "MICRO_DATABASE_URL",
			Usage:  "The database URL e.g root@tcp(127.0.0.1:3306)/database",
		},
		ccli.IntFlag{
			Name:   "register_ttl",
			EnvVar: "MICRO_REGISTER_TTL",
			Usage:  "Register TTL in seconds",
		},
		ccli.IntFlag{
			Name:   "register_interval",
			EnvVar: "MICRO_REGISTER_INTERVAL",
			Usage:  "Register interval in seconds",
		},
		ccli.StringFlag{
			Name:   "html_dir",
			EnvVar: "MICRO_HTML_DIR",
			Usage:  "The html directory for a web app",
		},
	)
}

func main() {
	app := cmd.App()
	app.Commands = append(app.Commands, auth.Commands()...)
	app.Commands = append(app.Commands, config.Commands()...)
	app.Commands = append(app.Commands, discovery.Commands()...)
	app.Commands = append(app.Commands, db.Commands()...)
	app.Commands = append(app.Commands, event.Commands()...)
	//	app.Commands = append(app.Commands, kv.Commands()...)
	//	app.Commands = append(app.Commands, log.Commands()...)
	//	app.Commands = append(app.Commands, metrics.Commands()...)
	app.Commands = append(app.Commands, monitor.Commands()...)
	app.Commands = append(app.Commands, router.Commands()...)
	//	app.Commands = append(app.Commands, sync.Commands()...)
	app.Commands = append(app.Commands, trace.Commands()...)
	app.Action = func(context *ccli.Context) { ccli.ShowAppHelp(context) }

	setup(app)

	cmd.Init(
		cmd.Name("platform"),
		cmd.Description("A microservices platform"),
		cmd.Version("latest"),
	)
}
