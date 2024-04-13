package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zetaoss/runbox/pkg/handler/lang"
	"github.com/zetaoss/runbox/pkg/handler/notebook"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/-/healthy", healthy)
	r.POST("/lang", lang.Run)
	r.POST("/notebook", notebook.Run)
	return r
}
