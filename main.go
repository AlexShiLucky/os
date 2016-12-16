package main

import (
	"os"

	ccli "github.com/micro/cli"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/os/auth"
	"github.com/micro/os/config"
	"github.com/micro/os/db"
	"github.com/micro/os/discovery"
	"github.com/micro/os/event"
	"github.com/micro/os/kv"
	//	"github.com/micro/os/log"
	//	"github.com/micro/os/metrics"
	"github.com/micro/os/monitor"
	"github.com/micro/os/router"
	//	"github.com/micro/os/sync"
	"github.com/micro/os/trace"
)

func init() {
	os.Setenv("MICRO_REGISTRY", "os")
}

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
	app.Commands = append(app.Commands, kv.Commands()...)
	//	app.Commands = append(app.Commands, log.Commands()...)
	//	app.Commands = append(app.Commands, metrics.Commands()...)
	app.Commands = append(app.Commands, monitor.Commands()...)
	app.Commands = append(app.Commands, router.Commands()...)
	//	app.Commands = append(app.Commands, sync.Commands()...)
	app.Commands = append(app.Commands, trace.Commands()...)
	app.Action = func(context *ccli.Context) { ccli.ShowAppHelp(context) }

	setup(app)

	cmd.Init(
		cmd.Name("os"),
		cmd.Description("A microservices operating system"),
		cmd.Version("latest"),
	)
}
