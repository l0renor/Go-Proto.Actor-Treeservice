package service

import (
	"fmt"
	"github.com/AsynkronIT/protoactor-go/actor"
	l "github.com/AsynkronIT/protoactor-go/log"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/logger"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// nolint:gochecknoglobals
var bindAddr string

func Main() {
	actor.SetLogLevel(l.ErrorLevel)
	app := &cli.App{
		Name:  "treeservice",
		Usage: "Start a tree service server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "bind",
				Usage:       "Address `HOST:PORT` to bind service to",
				Value:       "treeservice.actors:8090",
				Destination: &bindAddr,
			},
		},
		Action: func(context *cli.Context) error {
			var wg sync.WaitGroup
			wg.Add(1)
			remote.Start(bindAddr)
			remote.Register("tree", actor.PropsFromProducer(func() actor.Actor {
				return &Service{}
			}))
			wg.Wait()
			return nil
		},
	}
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println()
		os.Exit(0)
	}()
	err := app.Run(os.Args)
	if err != nil {
		logger.GetInstance().Error.Fatal(err)
	}
}
