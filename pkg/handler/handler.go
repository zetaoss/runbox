package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func healthy(c *gin.Context) {
	c.String(http.StatusOK, "Healthy.\n")
}
