package configs

type Client struct {
	Hostname string
	Port     int
}

func NewClientConfig(hostname string, port int) *Client {
	return &Client{
		Hostname: hostname,
		Port:     port,
	}
}
