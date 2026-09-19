package handlers

import (
	"fmt"
	"net/http"

	events "main/internal/events"
	services "main/internal/services"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "lmaaaaooooooo",
	})
}

func Create(c *gin.Context) {
	producers := events.NewPublisher()
	publisher := *producers.Producer
	message := events.CreateMessage(c)

	partition, offset, err := publisher.SendMessage(&sarama.ProducerMessage{
		Topic: "test-topic",
		Value: sarama.ByteEncoder(*message),
	})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("partition=%d, offset=%d\n", partition, offset)
}

//  admin handlers

func GetAccounts(c *gin.Context) {
	accounts := services.FetchAllAccounts()
	// accountString := fmt.Sprintf("%v", (*accounts))
	// data, err := json.Marshal(accountString)

	// if err != nil {
	// 	log.Fatalf("Failed to serialize response - %v\n", err)
	// }

	c.JSON(http.StatusOK, gin.H{
		"status": (*accounts),
	})
}
