package run

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zetaoss/zetarun/pkg/run/lang/single"
	"github.com/zetaoss/zetarun/pkg/run/lang/types"
)

type ResponseObj struct {
	Status string       `json:"status"`
	Error  string       `json:"error,omitempty"`
	Data   types.Output `json:"-"`
}

func Single(c *gin.Context) {
	var input types.SingleInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	output, err := single.Run(input)
	if err != nil {
		switch err {
		case types.ErrInvalidLanguage:
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   output,
	})
}
