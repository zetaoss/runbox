package run

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSingle(t *testing.T) {
	testcases := []struct {
		lang         string
		source       string
		wantCode     int
		wantResponse string
	}{
		{
			lang:         "_",
			source:       "echo hello",
			wantCode:     400,
			wantResponse: `{"error":"invalid language","status":"error"}`,
		},
		{
			lang:         "bash",
			source:       "echo hello",
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		{
			lang: "go",
			source: "" +
				"\n" + `package main` +
				"\n" +
				"\n" + `import "fmt"` +
				"\n" + `func main() {` +
				"\n" + `	fmt.Println("Hello, 世界")` +
				"\n" + `}` +
				"\n",
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0Hello, 世界"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.lang, func(t *testing.T) {
			requestData := map[string]string{
				"lang":   tc.lang,
				"source": tc.source,
			}
			requestBody, err := json.Marshal(requestData)
			if err != nil {
				panic("marshal request data error")
			}
			req := httptest.NewRequest("POST", "/run/single", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router1.ServeHTTP(w, req)
			require.Equal(t, tc.wantCode, w.Code)
			response := w.Body.String()
			// ignore json fields: [time, cpu, mem]
			re := regexp.MustCompile(`("time":)([^,]+)(,"cpu":)([0-9.]+)(,"mem":)([0-9.]+)`)
			response = re.ReplaceAllString(response, `${1}"0:00.00"${3}0${5}0`)
			require.Equal(t, tc.wantResponse, response)
		})
	}
}
