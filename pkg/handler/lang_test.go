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
			wantResponse: `{"logs":["1This is pdfTeX, Version 3.141592653-2.6-1.40.25 (TeX Live 2023/Debian) (preloaded format=pdflatex)","1 restricted \\write18 enabled.","1entering extended mode","1(./runbox.tex","1LaTeX2e \u003c2023-11-01\u003e patch level 1","1L3 programming layer \u003c2024-01-22\u003e","1(/usr/share/texlive/texmf-dist/tex/latex/base/article.cls","1Document Class: article 2023/05/17 v1.4n Standard LaTeX document class","1(/usr/share/texlive/texmf-dist/tex/latex/base/size10.clo))","1(/usr/share/texlive/texmf-dist/tex/latex/geometry/geometry.sty","1(/usr/share/texlive/texmf-dist/tex/latex/graphics/keyval.sty)","1(/usr/share/texlive/texmf-dist/tex/generic/iftex/ifvtex.sty","1(/usr/share/texlive/texmf-dist/tex/generic/iftex/iftex.sty)))","1(/usr/share/texlive/texmf-dist/tex/latex/l3backend/l3backend-pdftex.def)","1No file runbox.aux.","1*geometry* driver: auto-detecting","1*geometry* detected driver: pdftex","1[1{/var/lib/texmf/fonts/map/pdftex/updmap/pdftex.map}] (./runbox.aux) )\u003c/usr/sh","1are/texlive/texmf-dist/fonts/type1/public/amsfonts/cm/cmr10.pfb\u003e","1Output written on runbox.pdf (1 page, 12754 bytes).","1Transcript written on runbox.log."],"cpu":0,"mem":0,"time":0,"images":["iVBORw0KGgoAAAANSUhEUgAAAaQAAAEqCAQAAADw/+wWAAAE9UlEQVR42u3Y0Y0bRxZA0VeLTaBT6BSYAlNgCkpBzmCtELwB7IcnhQlhmcKkwBDaH/JoZGPXtuALkCOc80GCVSjiNcCLAriOAf6uf9x7APgeCAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICGnWdu8JeP8eOqR1Wj+tD2tbp/XTurz94Ndp/fz6+s3f+OXM+rh+nFnb/Pfez8n799AhHdfZ5+W4HdfZj6fj9tX69vr6zd/4duY6M3Pc5uXez8n79897D/At1nm2uR6/+eGv02zz8ra2TrPNbbZ5mW224/l1f51m/5zOzMzaZ5v93s/D9+Ohb6SZmTmt8zrPNrMusx9P8/HrzXWe/Xie8zp/WXqZD8d1LnOb29f7x3V+nG1OMzPrNJfjOq833PWvDwP/2+OHdD2ej+e5zcxlbus8t/X1TfJhnmfmZT68Lhy3mbXNbS6zH8+/2X85rsfTr6euM19Cer73I/L+PX5Ib27zcjzPpy8BfF7bf31/8zQf59+zz/5/9n//yY3E3/bQIa3TbHNe2zrNvi7zrzmv81yO2zrNvvZ1mn3+M+d1ntP88HbqeJrteJnPN82n1/11mn2dPp+cT3Ne+5xmX5t/7Sis494TwHfgoW8keC+EBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBDSg1vb+njvGfhzQnp0pznfewT+nJAe3PF87wn4K4QEASFBQEgQENKDW+fZ12Vt956DP7aOe08A3wE3EgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUHgF1IQfK0vXiM4AAAAAElFTkSuQmCC"]}`,
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.data), func(t *testing.T) {
			requestBody, err := json.Marshal(tc.data)
			if err != nil {
				panic("marshal request data error")
			}
			req := httptest.NewRequest("POST", "/lang", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			handler1.router.ServeHTTP(w, req)
			require.Equal(t, tc.wantCode, w.Code)
			response := w.Body.String()
			// ignore json fields: [time, cpu, mem]
			re := regexp.MustCompile(`(,"cpu":)([0-9.]+)(,"mem":)([0-9.]+)(,"time":)([0-9.]+)`)
			response = re.ReplaceAllString(response, `${1}0${3}0${5}0`)
			require.Equal(t, tc.wantResponse, response)
		})
	}
}
