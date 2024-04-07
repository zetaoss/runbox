package multi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type Params map[string]string

func TestRun_error(t *testing.T) {
	testcases := []struct {
		reqBody      string
		wantCode     int
		wantResponse string
	}{
		{
			reqBody:      ``,
			wantCode:     400,
			wantResponse: `{"error":"ErrBindJSON","status":"error"}`,
		},
		{
			reqBody:      `{}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrNoFiles","status":"error"}`,
		},
		{
			reqBody:      `{"asdfasdf": ""}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrNoFiles","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "bash"}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrNoFiles","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "_", "files": ""}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrBindJSON","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "_", "files": {}}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrBindJSON","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "_", "files": []}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrNoFiles","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "_", "files": "{}"}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrBindJSON","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "_", "files": "[]"}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrBindJSON","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "_", "files": ["echo hello"]}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrBindJSON","status":"error"}`,
		},
		{
			reqBody: `{"lang": "go", "source": ` +
				"\n" + `package main` +
				"\n" +
				"\n" + `import "fmt"` +
				"\n" + `func main() {` +
				"\n" + `	fmt.Println("Hello, 世界")` +
				"\n" + `}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrBindJSON","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "_", "files": [{}]}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrInvalidLanguage","status":"error"}`,
		},
	}
	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d %s", i, tc.reqBody), func(t *testing.T) {
			req := httptest.NewRequest("POST", "/multi", strings.NewReader(tc.reqBody))
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

func TestRun_ok(t *testing.T) {
	testcases := []struct {
		reqBody      string
		wantCode     int
		wantResponse string
	}{
		{
			reqBody:      `{"lang": "bash", "files": [{"text":"echo hello"}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		// no main
		{
			reqBody:      `{"lang": "bash", "files": [{"text":"cat greet.txt"},{"name":"greet.txt","text":"hello"}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		{
			reqBody:      `{"lang": "bash", "files": [{"name":"greet.txt","text":"hello"},{"text":"cat greet.txt"}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		// one main
		{
			reqBody:      `{"lang": "bash", "files": [{"text":"cat greet.txt","main":true},{"name":"greet.txt","text":"hello"}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		{
			reqBody:      `{"lang": "bash", "files": [{"name":"greet.txt","text":"hello"},{"text":"cat greet.txt","main":true}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		// two mains
		{
			reqBody:      `{"lang": "bash", "files": [{"text":"cat greet.txt","main":true},{"name":"greet.txt","text":"hello","main":true}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		{
			reqBody:      `{"lang": "bash", "files": [{"name":"greet.txt","text":"hello","main":true},{"text":"cat greet.txt","main":true}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		{
			reqBody:      `{"lang": "go", "files": [{"text":"package main\nimport \"fmt\"\nfunc main() {\n    fmt.Println(\"hello\")\n    fmt.Println(\"world\")\n}"}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello","0world"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		{
			reqBody:      `{"lang": "python", "files": [{"text":"print('hello')\nprint('world')"}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello","0world"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
	}
	for i, tc := range testcases {
		t.Run(fmt.Sprintf("%d %s", i, tc.reqBody), func(t *testing.T) {
			fmt.Println("tc.reqBody=", tc.reqBody)
			req := httptest.NewRequest("POST", "/multi", strings.NewReader(tc.reqBody))
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

func TestRun_warnings(t *testing.T) {
	testcases := []struct {
		reqBody      string
		wantCode     int
		wantResponse string
	}{
		{
			reqBody:      `{"lang": "python", "files": [{"text":"print(1000*\"X\")"}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0` + strings.Repeat("X", 1000) + `"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		{
			reqBody:      `{"lang": "python", "files": [{"text":"print(10000*\"X\")"}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0` + strings.Repeat("X", 8000) + `"],"warnings":["WarnOutputLimitReached"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/multi", strings.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router1.ServeHTTP(w, req)
			require.Equal(t, tc.wantCode, w.Code)
			response := w.Body.String()
			// ignore json fields: [time, cpu, mem]
			re := regexp.MustCompile(`("time":)([^,]+)(,"cpu":)([0-9.]+)(,"mem":)([0-9.]+)`)
			response = re.ReplaceAllString(response, `${1}"0:00.00"${3}0${5}0`)
			require.JSONEq(t, tc.wantResponse, response)
		})
	}
}

func TestRun_fakeErr(t *testing.T) {
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
			wantCode:     400,
			wantResponse: `{"error":"ErrNoFiles","status":"error"}`,
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
			req := httptest.NewRequest("POST", "/multi", bytes.NewBuffer(reqBody))
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
