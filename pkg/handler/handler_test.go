package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var router1 = NewRouter()

func TestHealthy(t *testing.T) {
	req := httptest.NewRequest("GET", "/-/healthy", nil)
	w := httptest.NewRecorder()
	router1.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Healthy.\n", w.Body.String())
}
