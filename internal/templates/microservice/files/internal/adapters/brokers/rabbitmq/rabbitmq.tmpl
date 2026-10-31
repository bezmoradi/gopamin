// Package rabbitmq publishes to and consumes from one durable RabbitMQ queue.
package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// connectTimeout bounds the dial, so an unreachable broker stops the app
	// within seconds.
	connectTimeout = 5 * time.Second
	// publishTimeout bounds one publish, including the broker's confirm.
	publishTimeout = 2 * time.Second
	// handlerTimeout bounds one attempt at processing a message.
	handlerTimeout = 30 * time.Second
	// minBackoff and maxBackoff bound the wait between attempts at a message
	// whose handler failed.
	minBackoff = 500 * time.Millisecond
	maxBackoff = 30 * time.Second
	// consumerTag names this process's one consumer on its channel.
	consumerTag = "consumer"
)

// Broker publishes to and consumes from the queue RABBITMQ_QUEUE. Publishing and
// consuming use separate channels, as RabbitMQ recommends.
type Broker struct {
	connection *amqp.Connection
	publisher  *amqp.Channel
	consumer   *amqp.Channel
	queue      string
}

// New connects with the RABBITMQ_* environment and declares the queue durable,
// so it and its persistent messages survive a broker restart.
func New() (*Broker, error) {
	queue := os.Getenv("RABBITMQ_QUEUE")
	if queue == "" {
		return nil, errors.New("RABBITMQ_QUEUE is not set")
	}

	// url.URL escapes the credentials, so a password containing @ or / can't
	// change which host the client dials.
	address := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASSWORD")),
		Host:   net.JoinHostPort(os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT")),
		Path:   "/",
	}
	connection, err := amqp.DialConfig(address.String(), amqp.Config{Dial: amqp.DefaultDial(connectTimeout)})
	if err != nil {
		return nil, fmt.Errorf("cannot reach RabbitMQ: %w", err)
	}

	b := &Broker{connection: connection, queue: queue}
	if err := b.open(); err != nil {
		_ = connection.Close()
		return nil, err
	}

	return b, nil
}

func (b *Broker) open() error {
	var err error
	if b.publisher, err = b.connection.Channel(); err != nil {
		return err
	}
	// Confirm mode: the broker acknowledges each publish once it has taken
	// responsibility for the message.
	if err := b.publisher.Confirm(false); err != nil {
		return err
	}

	if b.consumer, err = b.connection.Channel(); err != nil {
		return err
	}
	// Deliver one unacknowledged message at a time.
	if err := b.consumer.Qos(1, 0, false); err != nil {
		return err
	}

	// A queue that already exists as non-durable can't be redeclared durable
	// (PRECONDITION_FAILED, 406): delete it once, or pick another RABBITMQ_QUEUE.
	_, err = b.consumer.QueueDeclare(b.queue, true, false, false, false, nil)
	return err
}

// Produce publishes body as a persistent message and waits for the broker's
// confirm.
func (b *Broker) Produce(ctx context.Context, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, publishTimeout)
	defer cancel()

	confirmation, err := b.publisher.PublishWithDeferredConfirmWithContext(ctx, "", b.queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
	if err != nil {
		return err
	}

	acked, err := confirmation.WaitContext(ctx)
	if err != nil {
		return err
	}
	if !acked {
		return errors.New("rabbitmq refused the message")
	}
	return nil
}

// Consume hands each message to handler and acknowledges it only after handler
// returns nil, so a message whose processing didn't finish is delivered again
// (at least once). A failing handler is retried with backoff; there is no
// dead-letter queue, so a handler should return nil for a message it can never
// process. Consume returns nil when ctx ends (an unacknowledged message goes
// back to the queue) and an error if the connection or channel is lost.
func (b *Broker) Consume(ctx context.Context, handler func(context.Context, []byte) error) error {
	deliveries, err := b.consumer.Consume(b.queue, consumerTag, false, false, false, false, nil)
	if err != nil {
		return err
	}
	defer b.stopConsuming(deliveries)

	for {
		select {
		case <-ctx.Done():
			return nil
		case delivery, ok := <-deliveries:
			if !ok {
				if ctx.Err() != nil {
					return nil
				}
				return errors.New("rabbitmq delivery channel closed: the connection or channel was lost")
			}

			if !process(ctx, delivery.Body, handler) {
				_ = delivery.Nack(false, true)
				return nil
			}
			if err := delivery.Ack(false); err != nil {
				return fmt.Errorf("ack: %w", err)
			}
		}
	}
}

// stopConsuming cancels the consumer and returns to the queue whatever RabbitMQ
// had already pushed, until amqp091 closes the delivery channel. amqp091 holds the
// channel while a delivery waits to be received, so without this Close can hang.
func (b *Broker) stopConsuming(deliveries <-chan amqp.Delivery) {
	go func() { _ = b.consumer.Cancel(consumerTag, false) }()

	timeout := time.After(connectTimeout)
	for {
		select {
		case delivery, ok := <-deliveries:
			if !ok {
				return
			}
			_ = delivery.Nack(false, true)
		case <-timeout:
			return
		}
	}
}

// process runs handler until it succeeds, waiting longer after each failure.
// Each attempt gets its own timeout and isn't cancelled by ctx, so shutdown
// never interrupts a write halfway; it only stops further attempts. process
// reports false if ctx ended before handler succeeded.
func process(ctx context.Context, body []byte, handler func(context.Context, []byte) error) bool {
	backoff := minBackoff
	for {
		attemptCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), handlerTimeout)
		err := handler(attemptCtx, body)
		cancel()
		if err == nil {
			return true
		}

		select {
		case <-ctx.Done():
			return false
		case <-time.After(backoff):
		}
		backoff = min(2*backoff, maxBackoff)
	}
}

// Close closes both channels and the connection; RabbitMQ requeues any message
// left unacknowledged.
func (b *Broker) Close() error {
	return errors.Join(b.consumer.Close(), b.publisher.Close(), b.connection.Close())
}
