package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"github.com/zetaoss/runbox/pkg/runner/lang"
	"github.com/zetaoss/runbox/pkg/testutil"
)

var handler1 *Handler

func init() {
	d := testutil.NewDocker()
	handler1 = New(lang.New(box.New(d)))
}

func TestHealthy(t *testing.T) {
	req := httptest.NewRequest("GET", "/-/healthy", nil)
	w := httptest.NewRecorder()
	handler1.router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Healthy.\n", w.Body.String())
}
