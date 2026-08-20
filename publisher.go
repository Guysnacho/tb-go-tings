package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/IBM/sarama"
)

type producers struct {
	producer *sarama.SyncProducer
}

var (
	brokers  = ":29092"
	version  = "4.3.1"
	verbose  = false
	Producer sarama.SyncProducer
)

func NewPublisher() producers {
	brokers := strings.Split(brokers, ",")
	parsedVersion, _ := sarama.ParseKafkaVersion(version)
	Producer = newSyncProducer(brokers, parsedVersion)

	return producers{
		producer: &Producer,
	}
}

func newSyncProducer(brokerList []string, version sarama.KafkaVersion) sarama.SyncProducer {
	if Producer != nil {
		fmt.Println("Producer already configured and spun up. Returning og")
		return Producer
	}

	config := sarama.NewConfig()
	config.Version = version
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Net.TLS.Enable = false

	producer, err := sarama.NewSyncProducer(brokerList, config)

	if err != nil {
		log.Panicf("Failed to start Sarama producer - %s", err.Error())
	}

	return producer
}
