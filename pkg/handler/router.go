package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zetaoss/runbox/pkg/handler/notebook"
	"github.com/zetaoss/runbox/pkg/handler/run"
	"github.com/zetaoss/runbox/pkg/handler/status"
)

func NewRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/-/healthy", status.Healthy)
	r.POST("/api/notebook", notebook.Run)
	r.POST("/api/run", run.Run)
	return r
}
