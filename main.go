package main

import (
	fmt "fmt"
	"net/http"

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

	hello := r.router.Group("/hello")

	r.addTestRoutes(hello)

	return r
}

func (r routes) Run(addr ...string) error {
	return r.router.Run()
}

func Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "lmaaaaooooooo",
	})
}

func (r routes) addTestRoutes(rg *gin.RouterGroup) {
	rg.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET"},
	}))

	rg.GET("/", Hello)
}
