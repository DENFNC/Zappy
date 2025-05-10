package rabbit

import (
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL    string
	Locale string
}

func DialConfig(url string) *Config {
	return &Config{
		URL: url,
	}
}

type Client struct {
	cfg  *Config
	conn *amqp091.Connection
	mu   sync.Mutex
}

func NewClient(cfg *Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn != nil && !c.conn.IsClosed() {
		return nil
	}

	var err error
	c.conn, err = amqp091.Dial(c.cfg.URL)

	return err
}

func (c *Client) Channel() (*amqp091.Channel, error) {
	if err := c.Connect(); err != nil {
		return nil, err
	}
	return c.conn.Channel()
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil || c.conn.IsClosed() {
		return nil
	}
	return c.conn.Close()
}
