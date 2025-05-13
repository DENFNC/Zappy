package redis

import (
	"github.com/redis/go-redis/v9"
)

type Option func(*clientOption)

type Client struct {
	*redis.Client
}

type clientOption struct {
	addr     string
	password string
	db       int
	protocol int
}

func WithAddr(addr string) Option {
	return func(co *clientOption) {
		co.addr = addr
	}
}

func WithPassword(password string) Option {
	return func(co *clientOption) {
		co.password = password
	}
}

func WithDB(db int) Option {
	return func(co *clientOption) {
		co.db = db
	}
}

func WithProtocol(protocol int) Option {
	return func(co *clientOption) {
		co.protocol = protocol
	}
}

func NewClient(opts ...Option) *Client {
	var options = &clientOption{
		addr:     "localhost:6379",
		password: "",
		db:       0,
		protocol: 2,
	}

	for _, opt := range opts {
		opt(options)
	}

	client := redis.NewClient(
		&redis.Options{
			Addr:     options.addr,
			Password: options.password,
			DB:       options.db,
			Protocol: options.protocol,
		},
	)

	return &Client{
		Client: client,
	}
}
