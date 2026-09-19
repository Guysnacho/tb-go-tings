package events

import (
	"context"
	"errors"
	"log"
	test "main/internal/events/test"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var (
	topics []string = []string{
		"test-topic",
	}
)

func InitConsumers() {
	brokers := strings.Split(brokers, ",")
	parsedVersion, _ := sarama.ParseKafkaVersion(version)

	ctx, cancel := context.WithCancel(context.Background())
	cg, consumer := newConsumerGroup(brokers, parsedVersion)

	keepRunning := true
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := (*cg).Consume(ctx, topics, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // for some reason this awaits
	log.Println("Successfully spun up Sarama consumer. TURN ME UUUUUUPPPP")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: signal received")
			keepRunning = false
		}
	}
	cancel()
	wg.Wait()
	if err := (*cg).Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
	log.Panicln("Successfully finished closing consumers. Less than gracefully closing app.")
}

func newConsumerGroup(brokerList []string, version sarama.KafkaVersion) (*sarama.ConsumerGroup, Consumer) {
	config := sarama.NewConfig()
	config.Version = version

	// setup config
	config.Consumer.Retry.Max = 3
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Net.TLS.Enable = false

	// setup consumer
	consumer := Consumer{
		ready: make(chan bool),
	}

	// setup consumer group
	cg, err := sarama.NewConsumerGroup(brokerList, uuid.NewString(), config)

	if err != nil {
		log.Panicf("Failed to start Sarama consumer - %s", err.Error())
	}

	return &cg, consumer
}

type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Printf("message channel was closed")
				return nil
			}
			value, err := decodeMessage(&message.Value)
			if err != nil {
				log.Printf("Message claimed but errored: timestamp = %v, topic = %s", message.Timestamp, message.Topic)
				session.MarkMessage(message, "errored - "+err.Error())
			} else {
				log.Printf("Message claimed: value = %+v, timestamp = %v, topic = %s", value, message.Timestamp, message.Topic)
				session.MarkMessage(message, "")
			}
		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

func decodeMessage(data *[]byte) (*test.TestMessage, error) {
	var message test.TestMessage
	err := proto.Unmarshal((*data), &message)
	if err != nil {
		log.Printf("Failed to deserialize this one gangy... - %v\n", err)
		return nil, err
	}
	return &message, nil
}
