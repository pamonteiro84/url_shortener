package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"url_shortener/internal/handlers"
)

func Setup(r *gin.Engine, h *handlers.Handler) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.POST("/shorten", h.Shorten)
	r.GET("/:code", h.Redirect)
	r.DELETE("/:code", h.Delete)
}
