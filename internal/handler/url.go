package handler

import (
	"log"
	"net/http"
	"urlshortener/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UrlHandler struct {
	urlRepo repository.UrlRepository
}

func NewUrlHandler(urlRepository repository.UrlRepository) *UrlHandler {
	return &UrlHandler{
		urlRepo: urlRepository,
	}
}

func (h *UrlHandler) GetShortUrlHandler(c *gin.Context) {
	shortUrl := c.Param("url")

	originalUrl, err := h.urlRepo.GetOriginalUrl(shortUrl)

	if err != nil {
		log.Println(err)
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

func (h *UrlHandler) ShortenUrlHandler(c *gin.Context) {
	var requestBody ShortenUrlRequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid URL",
		})
		return
	}

	existingShortUrl, _ := h.urlRepo.GetShortUrl(requestBody.Url)

	if existingShortUrl != "" {
		// The original URL already has a short URL, return it
		c.JSON(http.StatusOK, gin.H{
			"url": existingShortUrl,
		})
		return
	}

	// If no existing short URL, generate a new one
	shortUrl := uuid.New().String()[:7]

	err := h.urlRepo.ShortenUrl(repository.ShortenUrlPayload{OriginalUrl: requestBody.Url, ShortUrl: shortUrl})

	if err != nil {
		log.Println(err)
		c.String(500, "Internal Server Error")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": shortUrl,
	})
}
