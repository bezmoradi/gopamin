// Package kafka publishes to and consumes from one Kafka topic.
package kafka

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	// connectTimeout bounds each dial, so an unreachable broker stops the app
	// within seconds.
	connectTimeout = 5 * time.Second
	// publishTimeout bounds one publish. When it expires the caller gets an
	// error, but the message can still be delivered afterwards: kafka-go keeps
	// the batch in flight. A failed publish is "maybe sent", not "not sent".
	publishTimeout = 2 * time.Second
	// handlerTimeout bounds one attempt at processing a message.
	handlerTimeout = 30 * time.Second
	// minBackoff and maxBackoff bound the wait between attempts at a message
	// whose handler failed.
	minBackoff = 500 * time.Millisecond
	maxBackoff = 30 * time.Second
)

// Broker publishes to KAFKA_TOPIC and consumes it as a member of the consumer
// group KAFKA_GROUP_ID.
type Broker struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

// New reads KAFKA_BROKER, KAFKA_TOPIC and KAFKA_GROUP_ID and creates the topic
// if it doesn't exist.
func New() (*Broker, error) {
	broker := os.Getenv("KAFKA_BROKER")
	topic := os.Getenv("KAFKA_TOPIC")
	group := os.Getenv("KAFKA_GROUP_ID")
	if broker == "" || topic == "" || group == "" {
		return nil, errors.New("KAFKA_BROKER, KAFKA_TOPIC and KAFKA_GROUP_ID must all be set")
	}

	// Create the topic up front so a fresh broker doesn't answer the first
	// produce or consume with "Unknown Topic Or Partition".
	if err := ensureTopic(broker, topic); err != nil {
		return nil, fmt.Errorf("cannot reach Kafka: %w", err)
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.Hash{},
		// Wait for every in-sync replica: a publish that returns nil is stored,
		// and one the broker rejects returns an error.
		RequiredAcks: kafka.RequireAll,
		// Send within 10 ms rather than waiting up to a second to fill a batch.
		BatchTimeout: 10 * time.Millisecond,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		// The group stores how far it has read: a restart resumes after the last
		// committed message instead of replaying the topic, and replicas split
		// the partitions instead of each reading everything.
		GroupID: group,
		// Hand over a message as soon as one arrives.
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	return &Broker{writer: writer, reader: reader}, nil
}

func ensureTopic(broker string, topic string) error {
	dialer := &kafka.Dialer{Timeout: connectTimeout}
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	conn, err := dialer.DialContext(ctx, "tcp", broker)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer func() { _ = controllerConn.Close() }()

	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil && !errors.Is(err, kafka.TopicAlreadyExists) {
		return err
	}

	return nil
}

// Produce publishes body and waits for every in-sync replica to store it.
func (b *Broker) Produce(ctx context.Context, body []byte) error {
	ctx, cancel := context.WithTimeout(ctx, publishTimeout)
	defer cancel()
	return b.writer.WriteMessages(ctx, kafka.Message{Value: body})
}

// Consume hands each message to handler and commits it only after handler
// returns nil, so a message whose processing didn't finish is delivered again
// (at least once). A failing handler is retried with backoff; there is no
// dead-letter topic, so a handler should return nil for a message it can never
// process. Consume returns nil when ctx ends, leaving the current message
// uncommitted, and an error if Kafka fails.
func (b *Broker) Consume(ctx context.Context, handler func(context.Context, []byte) error) error {
	for {
		message, err := b.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}

		if !process(ctx, message.Value, handler) {
			return nil
		}

		// The message is processed: commit even if shutdown has begun, so it
		// isn't delivered again.
		commitCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), connectTimeout)
		err = b.reader.CommitMessages(commitCtx, message)
		cancel()
		if err != nil {
			return fmt.Errorf("commit offset: %w", err)
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

// Close leaves the consumer group and flushes pending writes.
func (b *Broker) Close() error {
	return errors.Join(b.reader.Close(), b.writer.Close())
}
