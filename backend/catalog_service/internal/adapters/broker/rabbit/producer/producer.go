package producer

import (
	"context"

	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/broker/rabbit/option"
	"github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	channel *amqp091.Channel
	excOpts option.ExchangeOption
}

func New(ch *amqp091.Channel, opts option.ExchangeOption) (*Producer, error) {
	return &Producer{
		channel: ch,
		excOpts: opts,
	}, nil
}

func (rp *Producer) Publish(
	ctx context.Context,
	key string,
	body []byte,
) error {
	if err := rp.ensureExchange(); err != nil {
		return err
	}

	if err := rp.sendMessage(ctx, key, body); err != nil {
		return err
	}

	return nil
}

func (rp *Producer) Close() error {
	return rp.channel.Close()
}

func (rp *Producer) ensureExchange() error {
	return rp.channel.ExchangeDeclare(
		rp.excOpts.Name,
		rp.excOpts.Kind,
		rp.excOpts.Durable,
		rp.excOpts.AutoDelete,
		rp.excOpts.Internal,
		rp.excOpts.NoWait,
		rp.excOpts.Args,
	)
}

func (rp *Producer) sendMessage(
	ctx context.Context,
	key string,
	body []byte,
) error {
	return rp.channel.PublishWithContext(
		ctx,
		rp.excOpts.Name,
		key,
		false, //* mandatory //
		false, //* immediate //
		amqp091.Publishing{
			Body: body,
		},
	)
}
