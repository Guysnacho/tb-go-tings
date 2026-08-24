package main

import (
	fmt "fmt"

	events "main/internal/events"
	handlers "main/internal/handlers"

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
	go events.InitConsumers()
	return r.router.Run()
}

func (r routes) addTestRoutes(rg *gin.RouterGroup) {
	rg.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST"},
	}))

	rg.GET("/hello", handlers.Hello)
	rg.POST("", handlers.Create)
}
