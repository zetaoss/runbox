package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zetaoss/zetarun/pkg/handler/run"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/-/healthy", healthy)
	r.POST("/run/single", run.Single)
	return r
}
