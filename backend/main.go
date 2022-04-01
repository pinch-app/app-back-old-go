package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Build the container, which will automatically register
	// all the services, contorllers, and middleware & repos
	container := buildContainer()

	// Invoke the server with the container
	err := container.Invoke(func(router *gin.Engine) {
		router.Run(":8080")
	})

	if err != nil {
		log.Fatal(err)
	}
}
