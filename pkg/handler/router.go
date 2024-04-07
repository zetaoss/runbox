package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zetaoss/runbox/pkg/handler/multi"
	"github.com/zetaoss/runbox/pkg/handler/single"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/-/healthy", healthy)
	r.POST("/single", single.Run)
	r.POST("/multi", multi.Run)
	return r
}
