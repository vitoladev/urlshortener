package server

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.GET("/:url", s.GetShortUrlHandler)
	r.GET("/health", s.healthHandler)
	r.POST("/shorten", s.ShortenUrlHandler)

	return r
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
