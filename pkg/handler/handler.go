package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zetaoss/runbox/pkg/runner/lang"
)

type Handler struct {
	langRunner *lang.Lang
	router     *gin.Engine
}

func New(langRunner *lang.Lang) *Handler {
	h := &Handler{
		langRunner: langRunner,
	}
	h.setupRouter()
	return h
}

func (h *Handler) setupRouter() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/-/healthy", healthy)
	r.POST("/run/lang", h.lang)
	h.router = r
}

func (h *Handler) Run(addr ...string) error {
	return h.router.Run(addr...)
}

func healthy(c *gin.Context) {
	c.String(http.StatusOK, "Healthy.\n")
}
