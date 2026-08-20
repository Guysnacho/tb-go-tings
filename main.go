package main

import (
	fmt "fmt"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type routes struct {
	router *gin.Engine
}

func main() {
	fmt.Println("Testiiiiing 123")
	r := initRoutes()
	err := r.Run()

	if err != nil {
		panic(err)
	}
}

func initRoutes() routes {
	r := routes{
		router: gin.Default(),
	}

	hello := r.router.Group("/test")

	r.addTestRoutes(hello)

	return r
}

func (r routes) Run(addr ...string) error {
	// start consumers
	go InitConsumers()
	return r.router.Run()
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "lmaaaaooooooo",
	})
}

func Create(c *gin.Context) {
	producers := NewPublisher()
	publisher := *producers.producer
	// header := &test.TestHeader{
	// 	title:
	// }

	partition, offset, err := publisher.SendMessage(&sarama.ProducerMessage{
		Topic: "test-topic",
		Value: sarama.StringEncoder("sumn sumn"),
	})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("partition=%d, offset=%d\n", partition, offset)
}

func (r routes) addTestRoutes(rg *gin.RouterGroup) {
	rg.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST"},
	}))

	rg.GET("/hello", Hello)
	rg.POST("", Create)
}

// func createMessage(*c.Request.Body)  {

// }
