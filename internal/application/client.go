package application

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/withmandala/go-log"
	"github.com/yvv4git/task-wow/pkg/pow"
)

type Client struct {
	hostname     string
	port         int
	conn         net.Conn
	logger       *log.Logger
	powProcessor pow.POW
}

func NewClientApplication(hostname string, port int, logger *log.Logger, powProcessor pow.POW) *Client {
	return &Client{
		hostname:     hostname,
		port:         port,
		logger:       logger,
		powProcessor: powProcessor,
	}
}

func (c *Client) Connect() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", c.hostname, c.port))
	if err != nil {
		return fmt.Errorf("error on connect to server: %w", err)
	}
	defer func() {
		if errCloseConn := conn.Close(); errCloseConn != nil {
			c.logger.Errorf("error on close connection: %v", errCloseConn)
		}
	}()

	c.conn = conn

	receive, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("error on receive message from server: %w", err)
	}

	// Receive challenge.
	receive = strings.TrimSuffix(receive, "\n")
	c.logger.Infof("receive challenge: %s ", receive)

	proof := c.powProcessor.SolveChallenge([]byte(receive))
	rawMsg, err := json.Marshal(proof)
	if err != nil {
		return fmt.Errorf("error on marshal proof: %w", err)
	}

	// Send challenge result.
	messageToServer := fmt.Sprintf("%s\n", rawMsg)
	if _, err := fmt.Fprint(conn, messageToServer); err != nil {
		return fmt.Errorf("error on send message to server: %w", err)
	}

	// Got wow phrase message.
	receive, err = bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("error on receive message from server: %w", err)
	}

	c.logger.Infof("receive wow phrase: %s", receive)

	return nil
}
