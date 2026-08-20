package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
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

	<-consumer.ready // for some reason this awaits for some reason
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

func consumeMessage(message string) sarama.ConsumerGroupHandler {

}

type Consumer struct {
	ready chan bool
}
