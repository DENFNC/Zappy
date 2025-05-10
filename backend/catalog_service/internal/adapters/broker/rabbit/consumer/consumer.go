package consumer

import (
	"context"
	"sync"

	"github.com/DENFNC/Zappy/catalog_service/internal/adapters/broker/rabbit/option"
	"github.com/rabbitmq/amqp091-go"
)

type (
	CallbackHandleMessage func(context.Context, amqp091.Delivery) error
)

type Consumer struct {
	channel   *amqp091.Channel
	excOpts   option.ExchangeOption
	queueOpts option.QueueOption
}

func New(
	ch *amqp091.Channel,
	excOpts option.ExchangeOption,
	queueOpts option.QueueOption,
) *Consumer {
	return &Consumer{
		channel:   ch,
		excOpts:   excOpts,
		queueOpts: queueOpts,
	}
}

func (co *Consumer) Consume(
	ctx context.Context,
	routeKey, consumerName string,
	autoAck, exclusive bool,
	handler CallbackHandleMessage,
	wg *sync.WaitGroup,
) error {
	if err := co.declareExchange(); err != nil {
		return err
	}

	queueName, err := co.setupQueueAndBind(routeKey)
	if err != nil {
		return err
	}

	msgs, err := co.startConsuming(queueName, consumerName, autoAck, exclusive)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		co.handleMessages(
			ctx,
			consumerName,
			msgs,
			autoAck,
			handler,
		)
	}()

	return nil
}

func (co *Consumer) Close() error {
	return co.channel.Close()
}

func (co *Consumer) declareExchange() error {
	return co.channel.ExchangeDeclare(
		co.excOpts.Name, co.excOpts.Kind,
		co.excOpts.Durable, co.excOpts.AutoDelete,
		co.excOpts.Internal, co.excOpts.NoWait,
		co.excOpts.Args,
	)
}

func (co *Consumer) setupQueueAndBind(routeKey string) (string, error) {
	q, err := co.channel.QueueDeclare(
		co.queueOpts.Name, co.queueOpts.Durable,
		co.queueOpts.AutoDelete, co.queueOpts.Exclusive,
		co.queueOpts.NoWait, co.queueOpts.Args,
	)
	if err != nil {
		return "", err
	}
	if err := co.channel.QueueBind(
		q.Name, routeKey, co.excOpts.Name,
		false, co.queueOpts.Args,
	); err != nil {
		return "", err
	}
	return q.Name, nil
}

func (co *Consumer) startConsuming(
	queueName, consumer string,
	autoAck, exclusive bool,
) (<-chan amqp091.Delivery, error) {
	return co.channel.Consume(
		queueName, consumer,
		autoAck, exclusive,
		false, //* noLocal *//
		false, //* noWait *//
		co.queueOpts.Args,
	)
}

func (co *Consumer) handleMessages(
	ctx context.Context,
	consumerName string,
	msgs <-chan amqp091.Delivery,
	autoAck bool,
	handler CallbackHandleMessage,
) error {
	defer func() {
		_ = co.channel.Cancel(consumerName, false)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			if err := handler(ctx, msg); err != nil {
				return err
			}
			if !autoAck {
				msg.Ack(false)
			}
		}
	}
}
