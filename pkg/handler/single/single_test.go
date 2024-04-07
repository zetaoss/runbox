package single

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

func TestRun(t *testing.T) {
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
			wantResponse: `{"error":"ErrNoSource","status":"error"}`,
		},
		{
			reqBody:      `{"asdfasdf": ""}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrNoSource","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "bash"}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrNoSource","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "", "source": "echo hello"}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrInvalidLanguage","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "_", "source": "echo hello"}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrInvalidLanguage","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "bash", "source": "echo hello"}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0hello"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/single", strings.NewReader(tc.reqBody))
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

func TestRun_warning(t *testing.T) {
	testcases := []struct {
		params       Params
		wantCode     int
		wantResponse string
	}{
		{
			params:       Params{"lang": "python", "source": `print(10000*"X")`},
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0` + strings.Repeat("X", 8000) + `"],"warnings":["WarnOutputLimitReached"],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			reqBody, err := json.Marshal(tc.params)
			if err != nil {
				panic("marshal request data error")
			}
			req := httptest.NewRequest("POST", "/single", bytes.NewBuffer(reqBody))
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
			req := httptest.NewRequest("POST", "/single", bytes.NewBuffer(reqBody))
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
