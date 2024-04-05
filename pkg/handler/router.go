package handler

import (
	"github.com/gin-gonic/gin"
	single "github.com/zetaoss/runbox/pkg/handler/single"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/-/healthy", healthy)
	r.POST("/run/single", single.Single)
	return r
}
