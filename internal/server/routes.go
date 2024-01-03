package server

import (
	"net/http"
	"urlshortener/internal/database"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.GET("/:url", s.GetShortUrlHandler)
	r.GET("/health", s.healthHandler)
	r.POST("/shorten", s.ShortenUrlHandler)

	return r
}

func (s *Server) GetShortUrlHandler(c *gin.Context) {
	shortUrl := c.Param("url")
	resp := make(map[string]string)
	resp["message"] = shortUrl

	c.JSON(http.StatusOK, resp)
}

type ShortenUrlRequestBody struct {
	Url string `json:"url" binding:"required,url"`
}

func (s *Server) ShortenUrlHandler(c *gin.Context) {
	var requestBody ShortenUrlRequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid URL",
		})
		return
	}

	shortUrl := uuid.New().String()[:7]

	err := s.db.ShortenUrl(database.UrlData{ OriginalUrl: requestBody.Url, ShortUrl: shortUrl})
	
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
		return
	}

	c.JSON(200, gin.H{
		"url": shortUrl,
	})
}


func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
