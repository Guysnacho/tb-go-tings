package events

import (
	"fmt"
	"log"
	"main/internal/events/test"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
)

type Producers struct {
	Producer *sarama.SyncProducer
}

var (
	brokers  = ":29092"
	version  = "4.3.1"
	verbose  = false
	Producer sarama.SyncProducer
)

func NewPublisher() Producers {
	brokers := strings.Split(brokers, ",")
	parsedVersion, _ := sarama.ParseKafkaVersion(version)
	Producer = newSyncProducer(brokers, parsedVersion)

	return Producers{
		Producer: &Producer,
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

func CreateMessage(c *gin.Context) *[]byte {
	var reqBody TestRequestBody
	if err := c.BindJSON(&reqBody); err != nil {
		log.Panicf("Womp womp, invalid request body - %v\n", err)
	}

	header := test.TestHeader{
		Title:     reqBody.Header.Title,
		Author:    reqBody.Header.Author,
		Timestamp: reqBody.Header.Timestamp,
	}
	body := test.TestBody{
		Content: reqBody.Body.Content,
	}

	message := &test.TestMessage{
		Header: &header,
		Body:   &body,
	}

	data, err := proto.Marshal(message)
	if err != nil {
		log.Panicf("Failed to marshal message into byte array - %v\n", err)
	}
	return &data
}

type TestRequestBody struct {
	Header TestRequestBody_Header `json:"header" binding:"required"`
	Body   TestRequestBody_Body   `json:"body" binding:"required"`
}
type TestRequestBody_Header struct {
	Title     string `json:"title" binding:"required"`
	Author    string `json:"author" binding:"required"`
	Timestamp int64  `json:"timestamp" binding:"required"`
}
type TestRequestBody_Body struct {
	Content string `json:"content"`
}
