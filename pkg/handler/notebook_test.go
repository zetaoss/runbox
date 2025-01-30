package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/testutil"
)

func TestNotebook(t *testing.T) {
	testCases := []struct {
		data         map[string]any
		wantCode     int
		wantResponse string
	}{
		{
			data: map[string]any{
				"lang":    "bash",
				"sources": []string{},
			},
			wantCode:     400,
			wantResponse: `{"error":"invalid language"}`,
		},
		{
			data: map[string]any{
				"lang":    "python",
				"sources": []string{},
			},
			wantCode:     400,
			wantResponse: `{"error":"no sources"}`,
		},
		{
			data: map[string]any{
				"lang": "python",
				"sources": []string{
					`print("Hello, Python!")`,
				},
			},
			wantCode:     200,
			wantResponse: `{"outputsList":[[{"output_type":"stream","name":"stdout","text":["Hello, Python!\n"]}]],"cpu":0,"mem":0,"time":0,"timedout":false}`,
		},
		{
			data: map[string]any{
				"lang": "python",
				"sources": []string{
					`msg = "Hello, Python!!"`,
					`print(msg)`,
				},
			},
			wantCode:     200,
			wantResponse: `{"outputsList":[[],[{"output_type":"stream","name":"stdout","text":["Hello, Python!!\n"]}]],"cpu":0,"mem":0,"time":0,"timedout":false}`,
		},
		{
			data: map[string]any{
				"lang": "r",
				"sources": []string{
					`print("Hello, R!")`,
				},
			},
			wantCode:     200,
			wantResponse: `{"outputsList":[[{"output_type":"stream","name":"stdout","text":["[1] \"Hello, R!\"\n"]}]],"cpu":0,"mem":0,"time":0,"timedout":false}`,
		},
		{
			data: map[string]any{
				"lang": "r",
				"sources": []string{
					`msg <- "Hello, R!!"`,
					`print(msg)`,
				},
			},
			wantCode:     200,
			wantResponse: `{"outputsList":[[],[{"output_type":"stream","name":"stdout","text":["[1] \"Hello, R!!\"\n"]}]],"cpu":0,"mem":0,"time":0,"timedout":false}`,
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.data), func(t *testing.T) {
			requestBody, err := json.Marshal(tc.data)
			require.NoError(t, err)
			req := httptest.NewRequest("POST", "/notebook", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler1.router.ServeHTTP(w, req)
			require.Equal(t, tc.wantCode, w.Code)
			response := w.Body.String()
			re := regexp.MustCompile(`(,"cpu":)([0-9.]+)(,"mem":)([0-9.]+)(,"time":)([0-9.]+)`)
			response = re.ReplaceAllString(response, `${1}0${3}0${5}0`)
			require.Equal(t, tc.wantResponse, response)
		})
	}
}
