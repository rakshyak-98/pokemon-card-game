package main

import (
	"github.com/rakshyak-98/pokemonapi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.Use(middleware.AuthMiddleware())
}
