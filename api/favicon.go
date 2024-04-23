package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func faviconMiddleware(c *gin.Context) {
	if c.Request.URL.Path == "/favicon.ico" {
		c.String(http.StatusOK, "")
		c.Abort()
	}
	c.Next()
}
