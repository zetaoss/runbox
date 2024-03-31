package run

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zetaoss/runbox/pkg/run/lang/single"
	"github.com/zetaoss/runbox/pkg/run/lang/types"
)

type ResponseObj struct {
	Status string       `json:"status"`
	Error  string       `json:"error,omitempty"`
	Data   types.Output `json:"-"`
}

var fakeErr Error = NoError

func Single(c *gin.Context) {
	var input types.SingleInput
	if err := c.BindJSON(&input); err != nil || fakeErr == ErrBindJSON {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": ErrBindJSON})
		return
	}
	output, err := single.Run(input)
	if err != nil || fakeErr == ErrUnknown {
		switch err {
		case types.ErrNoSource:
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		case types.ErrInvalidLanguage:
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": ErrUnknown})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   output,
	})
}
