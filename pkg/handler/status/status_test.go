package status

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router1 *gin.Engine

func init() {
	gin.SetMode(gin.TestMode)
	router1 = gin.Default()
	router1.GET("/-/healthy", Healthy)
}

func TestHealthy(t *testing.T) {
	req := httptest.NewRequest("GET", "/-/healthy", nil)
	w := httptest.NewRecorder()
	router1.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Healthy.\n", w.Body.String())
}
