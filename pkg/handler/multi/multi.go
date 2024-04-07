package multi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zetaoss/runbox/pkg/runner/lang/multi"
	"github.com/zetaoss/runbox/pkg/runner/lang/types"
	"k8s.io/klog/v2"
)

type ResponseObj struct {
	Status string       `json:"status"`
	Error  string       `json:"error,omitempty"`
	Data   types.Output `json:"-"`
}

var fakeErr Error = NoError

func Run(c *gin.Context) {
	var input multi.Input
	if err := c.BindJSON(&input); err != nil || fakeErr == ErrBindJSON {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": ErrBindJSON})
		return
	}
	output, err := multi.Run(input)
	if err != nil || fakeErr == ErrUnknown {
		switch err {
		case multi.ErrNoFiles:
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		case multi.ErrInvalidLanguage:
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
