package main

import (
	"os"
	infrastructure "websocket/infraestructure"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	r := gin.Default()

	infrastructure.Routes(r)

	godotenv.Load()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
