package run

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type Params map[string]string

func TestSingle(t *testing.T) {
	testcases := []struct {
		params       Params
		wantCode     int
		wantResponse string
	}{
		{
			params:       nil,
			wantCode:     400,
			wantResponse: `{"error":"no source","status":"error"}`,
		},
		{
			params:       nil,
			wantCode:     400,
			wantResponse: `{"error":"no source","status":"error"}`,
		},
		{
			params:       Params{},
			wantCode:     400,
			wantResponse: `{"error":"no source","status":"error"}`,
		},
		{
			params:       Params{"asdfasdf": ""},
			wantCode:     400,
			wantResponse: `{"error":"no source","status":"error"}`,
		},
		{
			params:       Params{"lang": "bash"},
			wantCode:     400,
			wantResponse: `{"error":"no source","status":"error"}`,
		},
		{
			params:       Params{"lang": "_", "source": "echo hello"},
			wantCode:     400,
			wantResponse: `{"error":"invalid language","status":"error"}`,
		},
		{
			params:       Params{"lang": "bash", "source": "echo hello"},
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		{
			params: Params{"lang": "go", "source": "" +
				"\n" + `package main` +
				"\n" +
				"\n" + `import "fmt"` +
				"\n" + `func main() {` +
				"\n" + `	fmt.Println("Hello, 世界")` +
				"\n" + `}` +
				"\n"},
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0Hello, 世界"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			reqBody, err := json.Marshal(tc.params)
			if err != nil {
				panic("marshal request data error")
			}
			req := httptest.NewRequest("POST", "/run/single", bytes.NewBuffer(reqBody))
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

func TestSingle_warning(t *testing.T) {
	testcases := []struct {
		params       Params
		wantCode     int
		wantResponse string
	}{
		{
			params:       Params{"lang": "python", "source": `print(10000*"X")`},
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0` + strings.Repeat("X", 8000) + `"],"warning":1,"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			reqBody, err := json.Marshal(tc.params)
			if err != nil {
				panic("marshal request data error")
			}
			req := httptest.NewRequest("POST", "/run/single", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router1.ServeHTTP(w, req)
			require.Equal(t, tc.wantCode, w.Code)
			response := w.Body.String()
			// ignore json fields: [time, cpu, mem]
			re := regexp.MustCompile(`("time":)([^,]+)(,"cpu":)([0-9.]+)(,"mem":)([0-9.]+)`)
			response = re.ReplaceAllString(response, `${1}"0:00.00"${3}0${5}0`)
			require.Equal(t, len(tc.wantResponse), len(response))
		})
	}
}

func TestSingle_fakeErr(t *testing.T) {
	testcases := []struct {
		fakeErr      Error
		params       Params
		wantCode     int
		wantResponse string
	}{
		{
			fakeErr:      ErrBindJSON,
			params:       nil,
			wantCode:     400,
			wantResponse: `{"error":"ErrBindJSON","status":"error"}`,
		},
		{
			fakeErr:      ErrUnknown,
			params:       Params{"lang": "bash", "source": "echo hello"},
			wantCode:     500,
			wantResponse: `{"error":"ErrUnknown","status":"error"}`,
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			fakeErr = tc.fakeErr
			defer func() {
				fakeErr = NoError
			}()
			reqBody, err := json.Marshal(tc.params)
			if err != nil {
				panic("marshal request data error")
			}
			req := httptest.NewRequest("POST", "/run/single", bytes.NewBuffer(reqBody))
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
