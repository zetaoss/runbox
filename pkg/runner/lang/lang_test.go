package lang

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"github.com/zetaoss/runbox/pkg/testutil"
)

func TestRun_bash(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang:  "bash",
				Files: []box.File{{Body: `echo hello`}},
			},
			&box.Result{
				Logs: []box.Log{{Stream: 1, Log: "hello"}},
				CPU:  9461,
				MEM:  560,
				Time: 143,
			},
		},
		{
			Input{
				Lang:  "bash",
				Files: []box.File{{Body: `echo hello; echo world`}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "hello"},
					{Stream: 1, Log: "world"},
				},
				CPU:  8935,
				MEM:  564,
				Time: 192,
			},
		},
		{
			Input{
				Lang:  "bash",
				Files: []box.File{{Body: `echo hello; echo; echo world; echo`}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "hello"},
					{Stream: 1, Log: ""},
					{Stream: 1, Log: "world"},
					{Stream: 1, Log: ""},
				},
				CPU:  9512,
				MEM:  568,
				Time: 192,
			},
		},
		{
			Input{
				Lang:  "bash",
				Files: []box.File{{Body: `echo hello 1>&2; echo world 1>&2`}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 2, Log: "hello"},
					{Stream: 2, Log: "world"},
				},
				CPU:  9052,
				MEM:  548,
				Time: 16,
			},
		},
		{
			Input{
				Lang: "bash",
				Files: []box.File{
					{Body: `echo -n hello`},
					{Body: `echo -n world`},
				},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "helloworld"},
				},
				CPU:  9052,
				MEM:  548,
				Time: 16,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_c(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		// C
		{
			Input{
				Lang: "c",
				Files: []box.File{{Body: `
				#include <stdio.h>
				int main() {
					printf("Hello, World!");
				}`}},
			},
			&box.Result{
				Logs: []box.Log{{Stream: 1, Log: "Hello, World!"}},
				CPU:  25931,
				MEM:  484,
				Time: 393,
			},
		},
		{
			Input{
				Lang: "c",
				Files: []box.File{{Body: `
				#include <stdio.h>
				int main() {
					printf("Hello\nWorld!");
				}
			`}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "Hello"},
					{Stream: 1, Log: "World!"},
				},
				CPU:  26210,
				MEM:  788,
				Time: 35,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_cpp(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang: "cpp",
				Files: []box.File{{Body: `
				#include <iostream>
				int main() {
					std::cout<<"hello";
				}
			`}},
			},
			&box.Result{
				Logs: []box.Log{{Stream: 1, Log: "hello"}},
				CPU:  141223,
				MEM:  30496,
				Time: 1600,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_csharp(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang: "csharp",
				Files: []box.File{{Body: `
				using System;
				class Program
				{
					static void Main() {
						Console.Write("hello");
					}
				}
			`}},
			},
			&box.Result{
				Logs: []box.Log{{Stream: 1, Log: "hello"}},
				CPU:  107935,
				MEM:  25340,
				Time: 119,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_java(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang: "java",
				Files: []box.File{{Body: `
				public class App {
					public static void main(String args[]) {
						System.out.println("hello");
					}
				}`}},
			},
			&box.Result{
				Logs: []box.Log{{Stream: 1, Log: "hello"}},
				CPU:  734089,
				MEM:  772,
				Time: 752,
			},
		},
		{
			Input{
				Lang: "java",
				Files: []box.File{{Body: `
				import java.awt.Graphics2D;
				import java.awt.image.BufferedImage;
				import java.io.File;
				import java.io.IOException;
				import javax.imageio.ImageIO;

				public class App {
					public static void main(String[] args) throws IOException {
						BufferedImage bufferedImage = new BufferedImage(300, 200, BufferedImage.TYPE_INT_RGB);
						Graphics2D g = bufferedImage.createGraphics();
						g.fillRect(0, 0, 300, 100);
						g.drawString("Hello World", 120, 150);
						g.dispose();
						ImageIO.write(bufferedImage, "png", new File("myimage.png"));
					}
				}`}},
			},
			&box.Result{
				CPU:    997467,
				MEM:    972,
				Time:   1038,
				Images: []string{"iVBORw0KGgoAAAANSUhEUgAAASwAAADICAIAAADdvUsCAAACaUlEQVR4Xu3TQW4bMRAAQf3/08pBwGLD4VqyEbiloOpgkLMUfWHf7kDqtg6A3yVCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYjcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP6l+/2+LH7T1T+9msO7m293ThY/iHA5+fX2a1eHr+bw7ubbnZPFDyK8/X34av2Kq/NXc3h38+0uhRzmgacnz7Y3zPW85LE+hlfnz3P4JPPtHpPzy17Wc7Fdny2/ev2S5cLtfDkDn+TxfBfHp+Xk08V2e3jM59/z17n97hw+zHy751e+mAeWxXZ7eMzn3/PXuf3uHD7MfLtXr/wwD7zYw8zvfPLqku/O4cPMt3t+5dtItovtetpeeGy3l8xjx2K7hg8z3+558njcyxM/1k9PTtt7zpN5yXJs+bT9CQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADwv/sDs/dz4IQQs5EAAAAASUVORK5CYII="},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_kotlin(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang: "kotlin",
				Files: []box.File{{Body: `
				fun main() {
					println("Hello, Kotlin!")
				}`}},
			},
			&box.Result{
				Logs: []box.Log{{Stream: 1, Log: "Hello, Kotlin!"}},
				CPU:  5050983,
				MEM:  11516,
				Time: 14444,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_go(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang: "go",
				Files: []box.File{{Body: `
				package main
				import "fmt"
				func main() {
					fmt.Println("Hello, Go!")
				}`}},
			},
			&box.Result{
				Logs: []box.Log{{Stream: 1, Log: "Hello, Go!"}},
				CPU:  163465,
				MEM:  105064,
				Time: 608,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_latex(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang:  "latex",
				Files: []box.File{{Body: "\\documentclass{article}\n\\usepackage[a6paper,landscape]{geometry}\n\\begin{document}\nHello world!\n\\end{document}"}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "This is pdfTeX, Version 3.141592653-2.6-1.40.25 (TeX Live 2023/Debian) (preloaded format=pdflatex)"},
					{Stream: 1, Log: " restricted \\write18 enabled."},
					{Stream: 1, Log: "entering extended mode"},
					{Stream: 1, Log: "(./runbox.tex"},
					{Stream: 1, Log: "LaTeX2e <2023-11-01> patch level 1"},
					{Stream: 1, Log: "L3 programming layer <2024-01-22>"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/base/article.cls"},
					{Stream: 1, Log: "Document Class: article 2023/05/17 v1.4n Standard LaTeX document class"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/base/size10.clo))"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/geometry/geometry.sty"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/graphics/keyval.sty)"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/generic/iftex/ifvtex.sty"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/generic/iftex/iftex.sty)))"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/l3backend/l3backend-pdftex.def)"},
					{Stream: 1, Log: "No file runbox.aux."},
					{Stream: 1, Log: "*geometry* driver: auto-detecting"},
					{Stream: 1, Log: "*geometry* detected driver: pdftex"},
					{Stream: 1, Log: "[1{/var/lib/texmf/fonts/map/pdftex/updmap/pdftex.map}] (./runbox.aux) )</usr/sh"},
					{Stream: 1, Log: "are/texlive/texmf-dist/fonts/type1/public/amsfonts/cm/cmr10.pfb>"},
					{Stream: 1, Log: "Output written on runbox.pdf (1 page, 12754 bytes)."},
					{Stream: 1, Log: "Transcript written on runbox.log."},
				},
				Code:     0,
				CPU:      2552283,
				MEM:      147288,
				Time:     2714,
				Timedout: false,
				Images: []string{
					"iVBORw0KGgoAAAANSUhEUgAAAaQAAAEqCAQAAADw/+wWAAAE9UlEQVR42u3Y0Y0bRxZA0VeLTaBT6BSYAlNgCkpBzmCtELwB7IcnhQlhmcKkwBDaH/JoZGPXtuALkCOc80GCVSjiNcCLAriOAf6uf9x7APgeCAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICGnWdu8JeP8eOqR1Wj+tD2tbp/XTurz94Ndp/fz6+s3f+OXM+rh+nFnb/Pfez8n799AhHdfZ5+W4HdfZj6fj9tX69vr6zd/4duY6M3Pc5uXez8n79897D/At1nm2uR6/+eGv02zz8ra2TrPNbbZ5mW224/l1f51m/5zOzMzaZ5v93s/D9+Ohb6SZmTmt8zrPNrMusx9P8/HrzXWe/Xie8zp/WXqZD8d1LnOb29f7x3V+nG1OMzPrNJfjOq833PWvDwP/2+OHdD2ej+e5zcxlbus8t/X1TfJhnmfmZT68Lhy3mbXNbS6zH8+/2X85rsfTr6euM19Cer73I/L+PX5Ib27zcjzPpy8BfF7bf31/8zQf59+zz/5/9n//yY3E3/bQIa3TbHNe2zrNvi7zrzmv81yO2zrNvvZ1mn3+M+d1ntP88HbqeJrteJnPN82n1/11mn2dPp+cT3Ne+5xmX5t/7Sis494TwHfgoW8keC+EBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBDSg1vb+njvGfhzQnp0pznfewT+nJAe3PF87wn4K4QEASFBQEgQENKDW+fZ12Vt956DP7aOe08A3wE3EgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUHgF1IQfK0vXiM4AAAAAElFTkSuQmCC",
				},
			},
		},
		{
			Input{
				Lang: "latex",
				Files: []box.File{{Body: "\\documentclass{article}\\usepackage[a6paper,landscape]{geometry}" +
					"\\begin{document}\nLorem Ipsum 1\n" +
					"\\newpage\nLorem Ipsum 2\n" +
					"\\newpage\nLorem Ipsum 3\n" +
					"\\newpage\nLorem Ipsum 4\n" +
					"\\newpage\nLorem Ipsum 5\n" +
					"\\newpage\nLorem Ipsum 6\n" +
					"\\newpage\nLorem Ipsum 7\n" +
					"\\newpage\nLorem Ipsum 8\n" +
					"\\newpage\nLorem Ipsum 9\n" +
					"\\newpage\nLorem Ipsum 10\n" +
					"\\newpage\nLorem Ipsum 11\n" +
					"\\newpage\nLorem Ipsum 12\n" +
					"\\end{document}"}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "This is pdfTeX, Version 3.141592653-2.6-1.40.25 (TeX Live 2023/Debian) (preloaded format=pdflatex)"},
					{Stream: 1, Log: " restricted \\write18 enabled."},
					{Stream: 1, Log: "entering extended mode"},
					{Stream: 1, Log: "(./runbox.tex"},
					{Stream: 1, Log: "LaTeX2e <2023-11-01> patch level 1"},
					{Stream: 1, Log: "L3 programming layer <2024-01-22>"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/base/article.cls"},
					{Stream: 1, Log: "Document Class: article 2023/05/17 v1.4n Standard LaTeX document class"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/base/size10.clo))"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/geometry/geometry.sty"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/graphics/keyval.sty)"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/generic/iftex/ifvtex.sty"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/generic/iftex/iftex.sty)))"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/l3backend/l3backend-pdftex.def)"},
					{Stream: 1, Log: "No file runbox.aux."},
					{Stream: 1, Log: "*geometry* driver: auto-detecting"},
					{Stream: 1, Log: "*geometry* detected driver: pdftex"},
					{Stream: 1, Log: "[1{/var/lib/texmf/fonts/map/pdftex/updmap/pdftex.map}] [2] [3] [4] [5] [6]"},
					{Stream: 1, Log: "[7] [8] [9] [10] [11] [12] (./runbox.aux) )</usr/share/texlive/texmf-dist/fonts"},
					{Stream: 1, Log: "/type1/public/amsfonts/cm/cmr10.pfb>"},
					{Stream: 1, Log: "Output written on runbox.pdf (12 pages, 17111 bytes)."},
					{Stream: 1, Log: "Transcript written on runbox.log."},
				},
				Code:     0,
				CPU:      5966753,
				MEM:      28480,
				Time:     3609,
				Timedout: false,
				Images: []string{
					"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",
				},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_lua(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang:  "lua",
				Files: []box.File{{Body: `print("Hello, World!")`}},
			},
			&box.Result{
				Logs: []box.Log{{Stream: 1, Log: "Hello, World!"}},
				CPU:  9177,
				MEM:  3660,
				Time: 231,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_mysql(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang:  "mysql",
				Files: []box.File{{Body: `SELECT 'Hello, World!';`}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "+---------------+"},
					{Stream: 1, Log: "| Hello, World! |"},
					{Stream: 1, Log: "+---------------+"},
					{Stream: 1, Log: "| Hello, World! |"},
					{Stream: 1, Log: "+---------------+"},
				},
				CPU:  814171,
				MEM:  842256,
				Time: 2669,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_perl(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang: "perl",
				Files: []box.File{{Body: `
				use strict;
				use warnings;

				print "Hello, World!\n";`}},
			},
			&box.Result{
				Logs: []box.Log{{Stream: 1, Log: "Hello, World!"}},
				CPU:  13291,
				MEM:  7172,
				Time: 40,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_php(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang: "php",
				Files: []box.File{{Body: `<?php
						echo "Hello, World!";`}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "Hello, World!"},
				},
				CPU:  25300,
				MEM:  1456,
				Time: 114,
			},
		},
		{
			Input{
				Lang:  "php",
				Files: []box.File{{Body: `echo "Hello, World!";`}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "Hello, World!"},
				},
				CPU:  25300,
				MEM:  604,
				Time: 114,
			},
		},
		{
			Input{
				Lang: "php",
				Files: []box.File{
					{Body: `echo "hello";`},
					{Body: `echo "world";`},
				},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: `helloworld`},
				},
				CPU:  25300,
				MEM:  604,
				Time: 114,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_powershell(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang:  "powershell",
				Files: []box.File{{Body: `Write-Host "Hello, World!"`}},
			},
			&box.Result{
				Logs: []box.Log{{Stream: 1, Log: "Hello, World!"}},
				CPU:  279337,
				MEM:  75040,
				Time: 461,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_python(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang:  "python",
				Files: []box.File{{Body: `print("Hello, World!")`}},
			},
			&box.Result{
				Logs:     []box.Log{{Stream: 1, Log: "Hello, World!"}},
				Code:     0,
				CPU:      22896,
				MEM:      9316,
				Time:     60,
				Timedout: false,
				Images:   []string(nil),
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_r(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang:  "r",
				Files: []box.File{{Body: `print("Hello, World!")`}},
			},
			&box.Result{
				Logs:     []box.Log{{Stream: 1, Log: "[1] \"Hello, World!\""}},
				Code:     0,
				CPU:      102922,
				MEM:      28924,
				Time:     190,
				Timedout: false,
				Images:   []string(nil),
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_ruby(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang:  "ruby",
				Files: []box.File{{Body: `print("Hello, World!")`}},
			},
			&box.Result{
				Logs:     []box.Log{{Stream: 1, Log: "Hello, World!"}},
				Code:     0,
				CPU:      51197,
				MEM:      7228,
				Time:     88,
				Timedout: false,
				Images:   []string(nil),
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_sqlite3(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang:  "sqlite3",
				Files: []box.File{{Body: `SELECT 'Hello, World!';`}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "+-----------------+"},
					{Stream: 1, Log: "| 'Hello, World!' |"},
					{Stream: 1, Log: "+-----------------+"},
					{Stream: 1, Log: "| Hello, World!   |"},
					{Stream: 1, Log: "+-----------------+"},
				},
				Code:     0,
				CPU:      10842,
				MEM:      6504,
				Time:     437,
				Timedout: false,
				Images:   []string(nil),
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_tex(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{
				Lang:  "tex",
				Files: []box.File{{Body: "\\documentclass{article}\n\\usepackage[a6paper,landscape]{geometry}\n\\begin{document}\nHello world!\n\\end{document}"}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "This is pdfTeX, Version 3.141592653-2.6-1.40.25 (TeX Live 2023/Debian) (preloaded format=pdflatex)"},
					{Stream: 1, Log: " restricted \\write18 enabled."},
					{Stream: 1, Log: "entering extended mode"},
					{Stream: 1, Log: "(./runbox.tex"},
					{Stream: 1, Log: "LaTeX2e <2023-11-01> patch level 1"},
					{Stream: 1, Log: "L3 programming layer <2024-01-22>"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/base/article.cls"},
					{Stream: 1, Log: "Document Class: article 2023/05/17 v1.4n Standard LaTeX document class"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/base/size10.clo))"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/geometry/geometry.sty"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/graphics/keyval.sty)"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/generic/iftex/ifvtex.sty"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/generic/iftex/iftex.sty)))"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/l3backend/l3backend-pdftex.def)"},
					{Stream: 1, Log: "No file runbox.aux."},
					{Stream: 1, Log: "*geometry* driver: auto-detecting"},
					{Stream: 1, Log: "*geometry* detected driver: pdftex"},
					{Stream: 1, Log: "[1{/var/lib/texmf/fonts/map/pdftex/updmap/pdftex.map}] (./runbox.aux) )</usr/sh"},
					{Stream: 1, Log: "are/texlive/texmf-dist/fonts/type1/public/amsfonts/cm/cmr10.pfb>"},
					{Stream: 1, Log: "Output written on runbox.pdf (1 page, 12754 bytes)."},
					{Stream: 1, Log: "Transcript written on runbox.log."},
				},
				Code:     0,
				CPU:      2552283,
				MEM:      147288,
				Time:     2714,
				Timedout: false,
				Images: []string{
					"iVBORw0KGgoAAAANSUhEUgAAAaQAAAEqCAQAAADw/+wWAAAE9UlEQVR42u3Y0Y0bRxZA0VeLTaBT6BSYAlNgCkpBzmCtELwB7IcnhQlhmcKkwBDaH/JoZGPXtuALkCOc80GCVSjiNcCLAriOAf6uf9x7APgeCAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICGnWdu8JeP8eOqR1Wj+tD2tbp/XTurz94Ndp/fz6+s3f+OXM+rh+nFnb/Pfez8n799AhHdfZ5+W4HdfZj6fj9tX69vr6zd/4duY6M3Pc5uXez8n79897D/At1nm2uR6/+eGv02zz8ra2TrPNbbZ5mW224/l1f51m/5zOzMzaZ5v93s/D9+Ohb6SZmTmt8zrPNrMusx9P8/HrzXWe/Xie8zp/WXqZD8d1LnOb29f7x3V+nG1OMzPrNJfjOq833PWvDwP/2+OHdD2ej+e5zcxlbus8t/X1TfJhnmfmZT68Lhy3mbXNbS6zH8+/2X85rsfTr6euM19Cer73I/L+PX5Ib27zcjzPpy8BfF7bf31/8zQf59+zz/5/9n//yY3E3/bQIa3TbHNe2zrNvi7zrzmv81yO2zrNvvZ1mn3+M+d1ntP88HbqeJrteJnPN82n1/11mn2dPp+cT3Ne+5xmX5t/7Sis494TwHfgoW8keC+EBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBDSg1vb+njvGfhzQnp0pznfewT+nJAe3PF87wn4K4QEASFBQEgQENKDW+fZ12Vt956DP7aOe08A3wE3EgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUHgF1IQfK0vXiM4AAAAAElFTkSuQmCC",
				},
			},
		},
		{
			Input{
				Lang: "tex",
				Files: []box.File{{Body: "\\documentclass{article}\\usepackage[a6paper,landscape]{geometry}" +
					"\\begin{document}\nLorem Ipsum 1\n" +
					"\\newpage\nLorem Ipsum 2\n" +
					"\\newpage\nLorem Ipsum 3\n" +
					"\\newpage\nLorem Ipsum 4\n" +
					"\\newpage\nLorem Ipsum 5\n" +
					"\\newpage\nLorem Ipsum 6\n" +
					"\\newpage\nLorem Ipsum 7\n" +
					"\\newpage\nLorem Ipsum 8\n" +
					"\\newpage\nLorem Ipsum 9\n" +
					"\\newpage\nLorem Ipsum 10\n" +
					"\\newpage\nLorem Ipsum 11\n" +
					"\\newpage\nLorem Ipsum 12\n" +
					"\\end{document}"}},
			},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "This is pdfTeX, Version 3.141592653-2.6-1.40.25 (TeX Live 2023/Debian) (preloaded format=pdflatex)"},
					{Stream: 1, Log: " restricted \\write18 enabled."},
					{Stream: 1, Log: "entering extended mode"},
					{Stream: 1, Log: "(./runbox.tex"},
					{Stream: 1, Log: "LaTeX2e <2023-11-01> patch level 1"},
					{Stream: 1, Log: "L3 programming layer <2024-01-22>"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/base/article.cls"},
					{Stream: 1, Log: "Document Class: article 2023/05/17 v1.4n Standard LaTeX document class"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/base/size10.clo))"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/geometry/geometry.sty"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/graphics/keyval.sty)"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/generic/iftex/ifvtex.sty"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/generic/iftex/iftex.sty)))"},
					{Stream: 1, Log: "(/usr/share/texlive/texmf-dist/tex/latex/l3backend/l3backend-pdftex.def)"},
					{Stream: 1, Log: "No file runbox.aux."},
					{Stream: 1, Log: "*geometry* driver: auto-detecting"},
					{Stream: 1, Log: "*geometry* detected driver: pdftex"},
					{Stream: 1, Log: "[1{/var/lib/texmf/fonts/map/pdftex/updmap/pdftex.map}] [2] [3] [4] [5] [6]"},
					{Stream: 1, Log: "[7] [8] [9] [10] [11] [12] (./runbox.aux) )</usr/share/texlive/texmf-dist/fonts"},
					{Stream: 1, Log: "/type1/public/amsfonts/cm/cmr10.pfb>"},
					{Stream: 1, Log: "Output written on runbox.pdf (12 pages, 17111 bytes)."},
					{Stream: 1, Log: "Transcript written on runbox.log."},
				},
				Code:     0,
				CPU:      5966753,
				MEM:      28480,
				Time:     3609,
				Timedout: false,
				Images: []string{
					"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",
				},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}
