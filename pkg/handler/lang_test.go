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

func TestLang(t *testing.T) {
	testCases := []struct {
		data         map[string]any
		wantCode     int
		wantResponse string
	}{
		{
			data: map[string]any{
				"lang": "_",
				"files": []map[string]any{
					{"name": "greet.txt", "body": "hello"},
					{"name": "", "body": "cat greet.txt"},
				},
				"main": 1,
			},
			wantCode:     400,
			wantResponse: `{"error":"invalid language"}`,
		},
		{
			data: map[string]any{
				"lang":  "bash",
				"files": []map[string]any{},
			},
			wantCode:     400,
			wantResponse: `{"error":"no files"}`,
		},
		{
			data: map[string]any{
				"lang": "bash",
				"files": []map[string]any{
					{"name": "greet.txt", "body": "hello"},
					{"name": "", "body": "cat greet.txt"},
				},
				"main": 1,
			},
			wantCode:     200,
			wantResponse: `{"logs":["1hello"],"cpu":0,"mem":0,"time":0}`,
		},
		{
			data: map[string]any{
				"lang": "go",
				"files": []map[string]any{
					{
						"body": "" +
							"\n" + `package main` +
							"\n" +
							"\n" + `import "fmt"` +
							"\n" + `func main() {` +
							"\n" + `	fmt.Println("Hello, 世界")` +
							"\n" + `}` +
							"\n",
					},
				},
			},
			wantCode:     200,
			wantResponse: `{"logs":["1Hello, 世界"],"cpu":0,"mem":0,"time":0}`,
		},
		{
			data: map[string]any{
				"lang": "tex",
				"files": []map[string]any{
					{
						"body": "\\documentclass{article}\n\\usepackage[a6paper,landscape]{geometry}\n\\begin{document}\nHello world!\n\\end{document}",
					},
				},
			},
			wantCode:     200,
			wantResponse: `{"logs":["1This is pdfTeX, Version 3.141592653-2.6-1.40.22 (TeX Live 2021) (preloaded format=pdflatex)","1 restricted \\write18 enabled.","1entering extended mode","1(./runbox.tex","1LaTeX2e \u003c2020-10-01\u003e patch level 4","1L3 programming layer \u003c2021-02-18\u003e","1(/usr/local/texlive/2021/texmf-dist/tex/latex/base/article.cls","1Document Class: article 2020/04/10 v1.4m Standard LaTeX document class","1(/usr/local/texlive/2021/texmf-dist/tex/latex/base/size10.clo))","1(/usr/local/texlive/2021/texmf-dist/tex/latex/geometry/geometry.sty","1(/usr/local/texlive/2021/texmf-dist/tex/latex/graphics/keyval.sty)","1(/usr/local/texlive/2021/texmf-dist/tex/generic/iftex/ifvtex.sty","1(/usr/local/texlive/2021/texmf-dist/tex/generic/iftex/iftex.sty)))","1(/usr/local/texlive/2021/texmf-dist/tex/latex/l3backend/l3backend-pdftex.def)","1No file runbox.aux.","1*geometry* driver: auto-detecting","1*geometry* detected driver: pdftex","1[1{/usr/local/texlive/2021/texmf-var/fonts/map/pdftex/updmap/pdftex.map}]","1(./runbox.aux) )\u003c/usr/local/texlive/2021/texmf-dist/fonts/type1/public/amsfonts","1/cm/cmr10.pfb\u003e","1Output written on runbox.pdf (1 page, 11915 bytes).","1Transcript written on runbox.log."],"cpu":0,"mem":0,"time":0,"images":["iVBORw0KGgoAAAANSUhEUgAAAaQAAAEqCAQAAADw/+wWAAAE4klEQVR42u3X3Y0bRxpA0a8WTqBTYApMgSkwhUlBm4JCcAL74ElhQjBDMFNgCL0PkkY/gFcr6wLkCOc8NMkqdqMK6Ism1z7Az/rXvRcAvwIhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQeChQ1rH9ft6Wts6rj/W+evxT8cfv+Lr+6f1bmZt669775O376FD2i9zmOt+2y+z7c/fjH88/oMrfnKdbWa/zfXe++Tt++3eC/gR6zTbXPavbvx1nG2un8fWcba5zTbX2WbbXz7Nf3h9/dZhth/PEP7OQz+RZmbmuE7rNNvMOs9hf553X06u0xz3lzmv0+vQdZ72y5z369y+nN8v8/tsc5yZWcc575e5fTzjcu8t8vY9fkiX/WV/mdvMnOe2Tq+3/wdP8zwzl3n9B7XfZtY2t3Wew/7y1fx1v3z8gfg0l5nXK73ce4u8fY8f0me3ue4v8/6bscPH18+e5928n+NsfzP/7SdPJH7aQ4e0jrPNaW3rOId1nn/PaZ3mNLOOc1iHdZzD/GdO6zTHL/Pan2fbb/MhkPef5tdxDuv44cx5P6d1mOMc1ra2+fPe++TtW/u9VwC/gId+IsFbISQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkB7eerr3Cvg+IT24dZrzvdfA9wnpwe0v914B/w8hQUBIEBASBIT04NZpDuu8tnuvg/9t7fdeAfwCPJEgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAv8F6NR4np+BlBIAAAAASUVORK5CYII="]}`,
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.data), func(t *testing.T) {
			requestBody, err := json.Marshal(tc.data)
			if err != nil {
				panic("marshal request data error")
			}
			req := httptest.NewRequest("POST", "/run/lang", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler1.router.ServeHTTP(w, req)
			// require.Equal(t, tc.wantCode, w.Code)
			response := w.Body.String()
			// ignore json fields: [time, cpu, mem]
			re := regexp.MustCompile(`(,"cpu":)([0-9.]+)(,"mem":)([0-9.]+)(,"time":)([0-9.]+)`)
			response = re.ReplaceAllString(response, `${1}0${3}0${5}0`)
			require.Equal(t, tc.wantResponse, response)
		})
	}
}
