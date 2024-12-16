package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zetaoss/runbox/pkg/errors"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"github.com/zetaoss/runbox/pkg/runner/lang"
)

type Result struct {
	Logs     []string `json:"logs,omitempty"`
	Code     int      `json:"code,omitempty"`
	CPU      int      `json:"cpu,omitempty"`
	MEM      int      `json:"mem,omitempty"`
	Time     int      `json:"time,omitempty"`
	Timedout bool     `json:"timedout,omitempty"`
	Images   []string `json:"images,omitempty"`
}

func (h *Handler) lang(c *gin.Context) {
	var input lang.Input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.langRunner.Run(input)
	if err != nil {
		switch err {
		case errors.ErrInvalidLanguage:
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		case errors.ErrNoFiles:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, toResult(result))
}

func toResult(boxResult *box.Result) *Result {
	logs := make([]string, len(boxResult.Logs))
	for i, l := range boxResult.Logs {
		logs[i] = fmt.Sprintf("%d", l.Stream) + l.Log
	}
	return &Result{
		Logs:     logs,
		Code:     boxResult.Code,
		CPU:      boxResult.CPU,
		MEM:      boxResult.MEM,
		Time:     boxResult.Time,
		Timedout: boxResult.Timedout,
		Images:   boxResult.Images,
	}
}
