package helper

import (
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-web"
)

func WebOpts(ctx *cli.Context) []web.Option {
	var opts []web.Option

	if ttl := ctx.GlobalInt("register_ttl"); ttl > 0 {
		opts = append(opts, web.RegisterTTL(time.Duration(ttl)*time.Second))
	}

	if interval := ctx.GlobalInt("register_interval"); interval > 0 {
		opts = append(opts, web.RegisterInterval(time.Duration(interval)*time.Second))
	}

	if name := ctx.GlobalString("server_name"); len(name) > 0 {
		opts = append(opts, web.Name(name))
	}

	if ver := ctx.GlobalString("server_version"); len(ver) > 0 {
		opts = append(opts, web.Version(ver))
	}

	if id := ctx.GlobalString("server_id"); len(id) > 0 {
		opts = append(opts, web.Id(id))
	}

	if addr := ctx.GlobalString("server_address"); len(addr) > 0 {
		opts = append(opts, web.Address(addr))
	}

	if adv := ctx.GlobalString("server_advertise"); len(adv) > 0 {
		opts = append(opts, web.Advertise(adv))
	}

	return opts
}
