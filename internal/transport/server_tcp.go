package transport

import (
	"bufio"
	"context"
	"fmt"
	"net"

	"github.com/withmandala/go-log"
	"github.com/yvv4git/task-wow/internal/infrastructure/configs"
	"github.com/yvv4git/task-wow/internal/services"
)

type ServerTCP struct {
	cfg      *configs.Server
	logger   *log.Logger
	listener net.Listener
	svcWOW   *services.WOW
}

func NewTCP(cfg *configs.Server, logger *log.Logger, svcWOW *services.WOW) *ServerTCP {
	return &ServerTCP{
		cfg:    cfg,
		logger: logger,
		svcWOW: svcWOW,
	}
}

func (t *ServerTCP) Listen(ctx context.Context) error {
	listenConfig := net.ListenConfig{}

	listener, err := listenConfig.Listen(ctx, "tcp", fmt.Sprintf("%s:%d", t.cfg.Hostname, t.cfg.Port))
	if err != nil {
		return fmt.Errorf("error on listen server: %w", err)
	}

	t.listener = listener

	go func() {
		<-ctx.Done()
		t.logger.Infof("Shutdown server")
		if err := t.Close(); err != nil {
			t.logger.Errorf("error on close server: %v", err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if conn == nil {
				// Exit because the server is shutdown.
				return nil
			}

			if err := conn.Close(); err != nil {
				t.logger.Errorf("error on close client connection: %v", err)
			}

			continue
		}

		t.logger.Infof("accept incoming connection: %v", conn.RemoteAddr())
		go t.handleConnection(ctx, conn)
	}
}

func (t *ServerTCP) Close() error {
	if err := t.listener.Close(); err != nil {
		return fmt.Errorf("error on close listenner: %w", err)
	}

	return nil
}

func (t *ServerTCP) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		t.logger.Info("close client connection")
		if err := conn.Close(); err != nil {
			t.logger.Errorf("error on handling connecting: %v", err)
		}
	}()

	for {
		msg, err := t.svcWOW.Send(ctx)
		if err != nil {
			t.logger.Errorf("error on processing message for send: %v", err)

			return
		}

		if _, err := fmt.Fprint(conn, msg); err != nil {
			t.logger.Errorf("error on write msg to client connection: %v %T", err, err)

			return
		}

		receive, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			t.logger.Debugf("error on read from client connection: %v %T", err, err)

			return
		}

		if err := t.svcWOW.Receive(ctx, receive); err != nil {
			t.logger.Debugf("error on processing receive message: %v", err)

			return
		}

		t.logger.Infof("got from client: %s", receive)
	}
}
