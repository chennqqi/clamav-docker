package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/chennqqi/goutils/closeevent"
	unet "github.com/chennqqi/goutils/net"
	utime "github.com/chennqqi/goutils/time"
	"github.com/google/subcommands"
)

type webCmd struct {
	port     int
	zipto    string
	callback string

	w        *Web
	port     int
	fileto   string
	zipto    string
	callback string
	datadir  string
	indexdir string
	dns      string
}

func (p *webCmd) Name() string {
	return "web"
}

func (p *webCmd) Synopsis() string {
	return "web"
}

func (p *webCmd) Usage() string {
	return "web -p port"
}

func (p *webCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&p.port, "p", 8080, "set port")
	f.StringVar(&p.zipto, "timeout", "60s", "set scan timeout")
	f.StringVar(&p.fileto, "fileto", "20s", "set scan file timeout")
	f.StringVar(&p.datadir, "data", "/dev/shm", "set data dir")
	f.StringVar(&p.indexdir, "index", "/dev/shm/.persist", "set index dir")
	f.StringVar(&p.dns, "nameserver", "", "set ns server, can be list split by coma")
}

func (p *webCmd) Execute(context.Context, *flag.FlagSet, ...interface{}) subcommands.ExitStatus {
	zipTo, err := time.ParseDuration(p.zipto)
	if err != nil {
		zipTo, _ = time.ParseDuration("60s")
	}
	fileTo, err := time.ParseDuration(p.fileto)
	if err != nil {
		fileTo, _ = time.ParseDuration("20s")
	}

	if p.callback == "" {
		p.callback = os.Getenv("HMBD_CALLBACK")
	}

	if p.dns != "" {
		dns := strings.Split(p.dns, ",")
		if len(dns) > 0 {
			net.DefaultResolver = unet.NewResolver(dns)
		}
	} else {
		dns := os.Getenv("HMBD_DNS")
		dnslist := strings.Split(p.dns, ",")
		if dns != "" && len(dnslist) > 0 {
			net.DefaultResolver = unet.NewResolver(dnslist)
		}
	}

	var w Web
	w.clav, _ = NewClamAV("", false)
	w.fileto = utime.Duration(to)
	w.zipto = utime.Duration(to)
	w.callback = p.callback
	ctx, cancel := context.WithCancel(context.Background())
	go w.Run(p.port, ctx)

	closeevent.Wait(func(s os.Signal) {
		defer cancel()
		ctx := context.Background()
		w.Shutdown(ctx)
	}, os.Interrupt)

	return subcommands.ExitSuccess
}
