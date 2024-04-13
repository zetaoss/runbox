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

func TestRun_tex(t *testing.T) {
	testcases := []struct {
		reqBody      string
		wantCode     int
		wantResponse string
	}{
		// latex
		{
			reqBody:      `{"lang": "latex", "files": []}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrNoFiles","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "latex", "files": [{"text":""}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0This is pdfTeX, Version 3.141592653-2.6-1.40.22 (TeX Live 2021) (preloaded format=pdflatex)","0 restricted \\write18 enabled.","0entering extended mode","0(./runbox.tex","0LaTeX2e \u003c2020-10-01\u003e patch level 4","0L3 programming layer \u003c2021-02-18\u003e)","0*","0! Emergency stop.","0\u003c*\u003e runbox.tex","0              ","0!  ==\u003e Fatal error occurred, no output PDF file produced!","0Transcript written on runbox.log."],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		{
			reqBody:      `{"lang": "latex", "files": [{"text":"\\documentclass{minimal}\n\\begin{document}\nHello World!\n\\end{document}\n"}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":[],"images":["iVBORw0KGgoAAAANSUhEUgAAAlMAAANKCAQAAAAE5gOEAAAMUElEQVR42u3YwW3jRhiA0X/SgToI1IJaYAtuwS04JWRL2AZyiEuIS4hKiFtQCcwhtndtGLkEWX+C3zuQ0mg44Fw+iFz7AJT99NE3APDvZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIu9pMrdP6um7XYZ3W7+vm9fjz8Z1r/lindVh36+s6zqy7dffeui+fb9fdOqy/Pnqv8Lldbab28xzncb/s5zns92/Gn47vXHOZy36Z+9n2x5k571/eXffZ4xz2yzx+9F7hc7vaTL21tnWz3qRpndb2ZuxhbmbmOOd1ej3nn/O32eu4Tu+lDvjRrjtTp7WtbQ4z62aO+/28eoRb25z2h7lZ23eD/2Tq5fxtzn6er3OYl3jNzX6ey8zMnD96m/C5XXemzvvD/jCXmbmZy9qesvLsdu5n5jzfvbnaH2fWcWbuZ1vb/vBqzuN+fnl8vJ3zzNN6Dx+9TfjcrjtT31zmcX+YL2/Gjk/n793Pr/vDfnn69f05r7/7NwUf6moztU5zmG0d1mmO62Z+mW1ts82s0xyf3iv9Ntva5vQmXvcv/5EuM/Plec46zXGdnq+fL7Ot45zmuH6ePz96r/C5rf2j7+DHb/m4P86sw37572sB/79PmCngulztQx/wWcgUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUEPc3T4eA1oSjquEAAAAASUVORK5CYII="],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		// tex
		{
			reqBody:      `{"lang": "tex", "files": []}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrNoFiles","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "tex", "files": [{"text":""}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":["0This is pdfTeX, Version 3.141592653-2.6-1.40.22 (TeX Live 2021) (preloaded format=pdflatex)","0 restricted \\write18 enabled.","0entering extended mode","0(./runbox.tex","0LaTeX2e \u003c2020-10-01\u003e patch level 4","0L3 programming layer \u003c2021-02-18\u003e)","0*","0! Emergency stop.","0\u003c*\u003e runbox.tex","0              ","0!  ==\u003e Fatal error occurred, no output PDF file produced!","0Transcript written on runbox.log."],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
		},
		{
			reqBody:      `{"lang": "tex", "files": [{"text":"\\documentclass{minimal}\n\\begin{document}\nHello World!\n\\end{document}\n"}]}`,
			wantCode:     200,
			wantResponse: `{"data":{"logs":[],"images":["iVBORw0KGgoAAAANSUhEUgAAAlMAAANKCAQAAAAE5gOEAAAMUElEQVR42u3YwW3jRhiA0X/SgToI1IJaYAtuwS04JWRL2AZyiEuIS4hKiFtQCcwhtndtGLkEWX+C3zuQ0mg44Fw+iFz7AJT99NE3APDvZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIu9pMrdP6um7XYZ3W7+vm9fjz8Z1r/lindVh36+s6zqy7dffeui+fb9fdOqy/Pnqv8Lldbab28xzncb/s5zns92/Gn47vXHOZy36Z+9n2x5k571/eXffZ4xz2yzx+9F7hc7vaTL21tnWz3qRpndb2ZuxhbmbmOOd1ej3nn/O32eu4Tu+lDvjRrjtTp7WtbQ4z62aO+/28eoRb25z2h7lZ23eD/2Tq5fxtzn6er3OYl3jNzX6ey8zMnD96m/C5XXemzvvD/jCXmbmZy9qesvLsdu5n5jzfvbnaH2fWcWbuZ1vb/vBqzuN+fnl8vJ3zzNN6Dx+9TfjcrjtT31zmcX+YL2/Gjk/n793Pr/vDfnn69f05r7/7NwUf6moztU5zmG0d1mmO62Z+mW1ts82s0xyf3iv9Ntva5vQmXvcv/5EuM/Plec46zXGdnq+fL7Ot45zmuH6ePz96r/C5rf2j7+DHb/m4P86sw37572sB/79PmCngulztQx/wWcgUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUEPc3T4eA1oSjquEAAAAASUVORK5CYII="],"time":"0:00.00","cpu":0,"mem":0},"status":"success"}`,
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
