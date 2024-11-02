package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ----- Add new routes and other required functionalities in http server below ----

// Function to start the Gin server
func startServer(dataChannel chan Data) {
	r := gin.Default()

	r.GET("/services", func(ctx *gin.Context) {
		fmt.Println("--------- GOT UPDATE MESSAGE ---------")
		ctx.JSON(http.StatusOK, gin.H{"status": "data received"})
	})

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
