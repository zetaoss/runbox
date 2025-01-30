package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zetaoss/runbox/pkg/apperror"
	"github.com/zetaoss/runbox/pkg/runner/notebook"
)

func (h *Handler) notebook(c *gin.Context) {
	var input notebook.Input
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.notebookRunner.Run(input)
	if err != nil {
		if apperror.IsAppError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, result)
}
