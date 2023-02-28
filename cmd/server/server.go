package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/withmandala/go-log"
	"github.com/yvv4git/task-wow/internal/application"
	"github.com/yvv4git/task-wow/internal/infrastructure/configs"
	"github.com/yvv4git/task-wow/internal/services"
	"github.com/yvv4git/task-wow/internal/transport"
	"github.com/yvv4git/task-wow/pkg/pow"
)

const (
	defaultDifficulty = 2
	defaultPort       = 8095
)

func main() {
	logger := log.New(os.Stdout)

	app := &cli.App{
		Name:    "WOW ServerTCP server",
		Usage:   `The server sends phrases from the book "Word of Wisdom" to clients`,
		Version: "v0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "hostname",
				Value: "localhost",
				Usage: "binding server on this host",
			},
			&cli.IntFlag{
				Name:  "port",
				Value: defaultPort,
				Usage: "setup port number",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "enable debug level",
			},
			&cli.IntFlag{
				Name:  "difficulty",
				Value: defaultDifficulty,
				Usage: "setup difficulty for pow algorithm",
			},
		},
		Action: func(cCtx *cli.Context) error {
			cfg := configs.NewServerConfig(
				cCtx.String("hostname"),
				cCtx.Int("port"),
				cCtx.Bool("debug"),
				cCtx.Int("difficulty"),
			)
			logger.Infof("Config: %#v", cfg)

			processorPOW := pow.NewSHA256(cfg.Difficulty)
			serviceWOW := services.NewWOW(processorPOW)

			transportTCP := transport.NewTCP(cfg, logger, serviceWOW)
			wowApplication := application.NewServerApplication(transportTCP)

			if cfg.DebugLevel {
				logger = logger.WithDebug()
			}

			if err := wowApplication.Run(context.Background()); err != nil {
				if !errors.Is(err, application.ErrGracefullyShutdown) {
					return fmt.Errorf("error on run wow app: %w", err)
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatalf("error on run application: %v", err)
	}
}
