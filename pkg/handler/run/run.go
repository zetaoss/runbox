package run

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zetaoss/runbox/pkg/runner/run"
	"k8s.io/klog/v2"
)

type ResponseObj struct {
	Status string     `json:"status"`
	Error  string     `json:"error,omitempty"`
	Data   run.Output `json:"-"`
}

var fakeErr Error = NoError

func Run(c *gin.Context) {
	var input run.Input
	if err := c.BindJSON(&input); err != nil || fakeErr == ErrBindJSON {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": ErrBindJSON})
		return
	}
	output, err := run.Run(input)
	if err != nil || fakeErr == ErrUnknown {
		switch err {
		case run.ErrNoFiles:
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		case run.ErrInvalidLanguage:
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		default:
			klog.Warningf("unknown err: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": ErrUnknown})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   output,
	})
}
