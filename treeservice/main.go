package service

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/ob-vss-ws19/blatt-3-chupa-chups/logger"
	"github.com/urfave/cli/v2"
	"os"
	"sync"
)

var bindAddr string

func Main() {
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
			defer wg.Wait()
			remote.Start(bindAddr)
			remote.Register("tree", actor.PropsFromProducer(func() actor.Actor {
				return &Service{}
			}))
			return nil
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.GetInstance().Error.Fatal(err)
	}
}
