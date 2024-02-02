package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"urlshortener/internal/handler"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func NewServer(urlHandler *handler.UrlHandler) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/:url", urlHandler.GetShortUrlHandler)
	r.POST("/shorten", urlHandler.ShortenUrlHandler)

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
