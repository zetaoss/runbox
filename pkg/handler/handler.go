package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zetaoss/runbox/pkg/runner/lang"
	"github.com/zetaoss/runbox/pkg/runner/notebook"
)

type Handler struct {
	langRunner     *lang.Lang
	notebookRunner *notebook.Notebook
	router         *gin.Engine
}

func New(langRunner *lang.Lang, notebookRunner *notebook.Notebook) *Handler {
	h := &Handler{
		langRunner:     langRunner,
		notebookRunner: notebookRunner,
	}
	h.setupRouter()
	return h
}

func (h *Handler) setupRouter() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/-/healthy", healthy)
	r.POST("/lang", h.lang)
	r.POST("/notebook", h.notebook)
	h.router = r
}

func (h *Handler) Run(addr ...string) error {
	return h.router.Run(addr...)
}

func healthy(c *gin.Context) {
	c.String(http.StatusOK, "Healthy.\n")
}
