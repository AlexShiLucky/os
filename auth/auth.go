package auth

import (
	"log"
	"time"

	"github.com/micro/auth-srv/db"
	"github.com/micro/auth-srv/db/mysql"
	"github.com/micro/auth-srv/handler"
	account "github.com/micro/auth-srv/proto/account"
	oauth2 "github.com/micro/auth-srv/proto/oauth2"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

func srv(ctx *cli.Context) {
	service := micro.NewService(
		micro.Name("go.micro.srv.auth"),
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

	// register account handler
	account.RegisterAccountHandler(service.Server(), new(handler.Account))
	// register oauth2 handler
	oauth2.RegisterOauth2Handler(service.Server(), new(handler.Oauth2))

	// initialise database
	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func authCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "srv",
			Usage:  "Run the auth server",
			Action: srv,
		},
	}
}

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:        "auth",
			Usage:       "Auth commands",
			Subcommands: authCommands(),
		},
	}
}
