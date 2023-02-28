package application

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	transport2 "github.com/yvv4git/task-wow/internal/transport"
)

var ErrGracefullyShutdown = errors.New("got system signal for gracefully shutdown")

type ServerWOW struct {
	runGroup  run.Group
	transport transport2.Transport
}

func NewServerApplication(transport transport2.Transport) *ServerWOW {
	var g run.Group

	return &ServerWOW{
		runGroup:  g,
		transport: transport,
	}
}

func (s *ServerWOW) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	s.runGroup.Add(func() error {
		err := s.transport.Listen(ctx)
		if err != nil {
			return fmt.Errorf("error on transport listen: %w", err)
		}

		return nil
	}, func(error) {
		cancel()
	})

	s.runGroup.Add(func() error {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)
		defer signal.Stop(signals)
		<-signals

		return ErrGracefullyShutdown
	}, func(error) {
		cancel()
	})

	err := s.runGroup.Run()

	return fmt.Errorf("error on run server app: %w", err)
}
