package option

import "github.com/rabbitmq/amqp091-go"

type CommonArgs struct {
	Durable    bool
	AutoDelete bool
	NoWait     bool
	Args       amqp091.Table
}

func DefaultCommonArgs() CommonArgs {
	return CommonArgs{
		Durable:    true,
		AutoDelete: false,
		NoWait:     false,
		Args:       amqp091.Table{},
	}
}

type ExchangeOption struct {
	CommonArgs
	Name     string
	Kind     string
	Internal bool
}

func NewExchangeOption(name, kind string) ExchangeOption {
	return ExchangeOption{
		CommonArgs: DefaultCommonArgs(),
		Name:       name,
		Kind:       kind,
		Internal:   false,
	}
}

type QueueOption struct {
	CommonArgs
	Name      string
	Exclusive bool
}

func NewQueueOption(name string) QueueOption {
	return QueueOption{
		CommonArgs: DefaultCommonArgs(),
		Name:       name,
		Exclusive:  false,
	}
}
