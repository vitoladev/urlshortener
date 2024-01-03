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

	originalUrl, err := s.db.GetOriginalUrl(shortUrl)

	if err != nil {
		c.String(500, "Internal Server Error")
		return
	}

	if originalUrl == "" {
		c.String(404, "Short URL not found")
		return
	}

	// Redirect the user to the original URL
	c.Redirect(302, originalUrl)
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

	existingShortUrl, _ := s.db.GetShortUrl(requestBody.Url)

	if existingShortUrl != "" {
		// The original URL already has a short URL, return it
		c.JSON(http.StatusOK, gin.H{
			"url": existingShortUrl,
		})
		return
	}

	// If no existing short URL, generate a new one
	shortUrl := uuid.New().String()[:7]

	err := s.db.ShortenUrl(database.UrlData{OriginalUrl: requestBody.Url, ShortUrl: shortUrl})

	if err != nil {
		c.String(500, "Internal Server Error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": shortUrl,
	})
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
