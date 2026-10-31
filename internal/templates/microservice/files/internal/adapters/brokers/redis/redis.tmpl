// Package redis publishes to and subscribes to one Redis pub/sub channel.
//
// Pub/sub is at-most-once: Redis delivers a message only to the subscribers
// connected at that moment and keeps nothing, so a message published while this
// process is down or reconnecting is lost, and every replica receives every
// message. Use Kafka or RabbitMQ when a message must not be lost.
package redis

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// connectTimeout bounds the dial and the subscription, so an unreachable
	// server stops the app within seconds.
	connectTimeout = 5 * time.Second
	// publishTimeout bounds one publish.
	publishTimeout = 2 * time.Second
	// handlerTimeout bounds one attempt at processing a message.
	handlerTimeout = 30 * time.Second
	// minBackoff and maxBackoff bound the wait between attempts at a message
	// whose handler failed.
	minBackoff = 500 * time.Millisecond
	maxBackoff = 30 * time.Second
)

// Broker publishes to and subscribes to the channel REDIS_CHANNEL.
type Broker struct {
	client  *redis.Client
	sub     *redis.PubSub
	channel string
}

// New connects with the REDIS_* environment (REDIS_TLS=true for TLS) and
// subscribes before returning, because pub/sub drops a message published while
// nobody is subscribed: the sample message main publishes would be lost.
func New() (*Broker, error) {
	channel := os.Getenv("REDIS_CHANNEL")
	if channel == "" {
		return nil, errors.New("REDIS_CHANNEL is not set")
	}

	host := os.Getenv("REDIS_HOST")
	options := &redis.Options{
		Addr:        net.JoinHostPort(host, os.Getenv("REDIS_PORT")),
		Password:    os.Getenv("REDIS_PASSWORD"),
		DialTimeout: connectTimeout,
	}
	if os.Getenv("REDIS_TLS") == "true" {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12, ServerName: host}
	}
	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()
	sub := client.Subscribe(ctx, channel)
	if _, err := sub.Receive(ctx); err != nil {
		_ = sub.Close()
		_ = client.Close()
		return nil, fmt.Errorf("cannot reach Redis: %w", err)
	}

	return &Broker{client: client, sub: sub, channel: channel}, nil
}

// Produce publishes body to the channel's current subscribers.
func (b *Broker) Produce(ctx context.Context, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, publishTimeout)
	defer cancel()
	return b.client.Publish(ctx, b.channel, body).Err()
}

// Consume hands each message to handler, retrying a failing handler with
// backoff. Messages that arrive meanwhile wait in the client's buffer, which
// drops them once full. Consume returns nil when ctx ends.
func (b *Broker) Consume(ctx context.Context, handler func(context.Context, []byte) error) error {
	messages := b.sub.Channel()
	for {
		select {
		case <-ctx.Done():
			return nil
		case message, ok := <-messages:
			if !ok {
				return nil
			}
			if !process(ctx, []byte(message.Payload), handler) {
				return nil
			}
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

// Close unsubscribes and closes the connection pool.
func (b *Broker) Close() error {
	return errors.Join(b.sub.Close(), b.client.Close())
}
