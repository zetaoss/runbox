package multi

import "github.com/gin-gonic/gin"

var router1 *gin.Engine

func init() {
	gin.SetMode(gin.TestMode)
	router1 = gin.Default()
	router1.POST("/multi", Run)
}
