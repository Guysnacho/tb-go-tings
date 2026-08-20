package main

import (
	fmt "fmt"
	"log"
	"net/http"

	"main/protobuf/tunjiproductions.com/events/test"

	"google.golang.org/protobuf/proto"

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
	message := createMessage(c)

	partition, offset, err := publisher.SendMessage(&sarama.ProducerMessage{
		Topic: "test-topic",
		Value: sarama.ByteEncoder(*message),
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

func createMessage(c *gin.Context) *[]byte {
	var reqBody TestRequestBody
	if err := c.BindJSON(&reqBody); err != nil {
		log.Panicln("Womp womp, invalid request body - %v", err)
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
		log.Panicln("Failed to marshal message into byte array - %v", err)
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
