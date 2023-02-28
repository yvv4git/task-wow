package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/withmandala/go-log"
	"github.com/yvv4git/task-wow/internal/application"
	"github.com/yvv4git/task-wow/internal/infrastructure/configs"
	"github.com/yvv4git/task-wow/pkg/pow"
)

const (
	defaultDifficulty = 2
	defaultPort       = 8095
)

func main() {
	logger := log.New(os.Stdout)

	app := &cli.App{
		Name:    "WOW TCP client",
		Usage:   `The client sends request to server "Word of Wisdom"`,
		Version: "v0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "hostname",
				Value: "localhost",
				Usage: "setup server host name",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: defaultPort,
				Usage: "setup server port number",
			},
		},
		Action: func(cCtx *cli.Context) error {
			cfg := configs.NewClientConfig(cCtx.String("hostname"), cCtx.Int("port"))
			logger.Infof("Config: %#v", cfg)

			powProcessor := pow.NewSHA256(defaultDifficulty)
			clientApplication := application.NewClientApplication(cfg.Hostname, cfg.Port, logger, powProcessor)

			err := clientApplication.Connect()
			if err != nil {
				return fmt.Errorf("error on connect client application: %w", err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatalf("error on run application: %v", err)
	}
}
