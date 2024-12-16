package lang

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"github.com/zetaoss/runbox/pkg/testutil"
	"k8s.io/utils/ptr"
)

var lang1 *Lang

func init() {
	d := testutil.NewDocker()
	lang1 = New(box.New(d))
}

func equalResult(t *testing.T, want, got *box.Result) {
	t.Helper()

	assert.Greater(t, got.CPU, want.CPU/100, "want.CPU", want.CPU)
	assert.Greater(t, got.MEM, want.MEM/1000, "want.MEM", want.MEM)
	assert.Less(t, got.CPU, want.CPU*100, "want.CPU", want.CPU)
	assert.Less(t, got.MEM, want.MEM*1000, "want.MEM", want.MEM)
	want.CPU = got.CPU
	want.MEM = got.MEM

	assert.Greater(t, got.Time, want.Time/10, "want.Time", want.Time)
	assert.Less(t, got.Time, want.Time*10, "want.Time", want.Time)
	want.Time = got.Time

	assert.Equal(t, want, got)
}

func TestToLangOpts(t *testing.T) {
	testcases := []struct {
		input     Input
		want      *LangOpts
		wantError string
	}{
		{
			Input{
				Lang: "bash",
				Files: []box.File{
					{Name: "greet.txt", Body: "hello"},
					{Body: "cat greet.txt"},
				},
				Main: 1,
			},
			&LangOpts{
				Input:          Input{Lang: "bash", Files: []box.File{{Name: "greet.txt", Body: "hello"}, {Body: "cat greet.txt"}}, Main: 1},
				Command:        "/bin/bash runbox.sh",
				FileName:       "runbox",
				FileExt:        "sh",
				Shell:          "bash",
				TimeoutSeconds: 10,
				WorkingDir:     "/home/user01",
			},
			"",
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			got, err := toLangOpts(tc.input)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestToBoxOpts(t *testing.T) {
	testcases := []struct {
		langOpts LangOpts
		want     box.Opts
	}{
		{
			LangOpts{
				Input:          Input{Lang: "bash", Files: []box.File{{Name: "greet.txt", Body: "hello"}, {Body: "cat greet.txt"}}, Main: 1},
				Command:        "/bin/bash runbox.sh",
				FileName:       "runbox",
				FileExt:        "sh",
				Shell:          "bash",
				TimeoutSeconds: 10,
				WorkingDir:     "/home/user01",
			},
			box.Opts{
				CollectStats:  ptr.To(true),
				CollectImages: true,
				Command:       "/bin/bash runbox.sh",
				Env:           nil,
				Files: []box.File{
					{Name: "/home/user01/greet.txt", Body: "hello"},
					{Name: "/home/user01/runbox.sh", Body: "cat greet.txt"},
				},
				Image:      "jmnote/runbox:bash",
				Shell:      "bash",
				Timeout:    10000,
				WorkingDir: "/home/user01",
			},
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			got := toBoxOpts(tc.langOpts)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestRun_simple(t *testing.T) {
	testcases := []struct {
		input     Input
		want      *box.Result
		wantError string
	}{
		{
			Input{
				Lang: "bash",
				Files: []box.File{
					{Name: "greet.txt", Body: "hello"},
					{Body: "cat greet.txt"},
				},
				Main: 1,
			},
			&box.Result{
				Logs:     []box.Log{{Stream: 1, Log: "hello"}},
				CPU:      9183,
				MEM:      676,
				Time:     196,
				Timedout: false,
			},
			"",
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_error(t *testing.T) {
	testCases := []struct {
		langInput Input
		wantError string
	}{
		{Input{Lang: "", Files: []box.File{}}, "no files"},
		{Input{Lang: "x", Files: []box.File{}}, "no files"},
		{Input{Lang: "go", Files: []box.File{}}, "no files"},
		{Input{Lang: "", Files: []box.File{{Body: `echo hello`}}}, "invalid language"},
		{Input{Lang: "x", Files: []box.File{{Body: `echo hello`}}}, "invalid language"},
	}
	for _, tc := range testCases {
		t.Run(testutil.Name(tc.langInput), func(t *testing.T) {
			output, err := lang1.Run(tc.langInput)
			require.Nil(t, output)
			require.EqualError(t, err, tc.wantError)
		})
	}
}

func TestRun_singleFile_bash(t *testing.T) {
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
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_singleFile_cxx(t *testing.T) {
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
		// C++
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
		// C#
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

func TestRun_singleFile_php(t *testing.T) {
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
		// PHP
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
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_singleFile_java_kotlin(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		// Java
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
		// Kotlin
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

func TestRun_singleFile_other_languages(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		// Go
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
		// Lua
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
		// Perl
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
		// PowerShell
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
		// Python
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
		// R
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
		// Ruby
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

func TestRun_singleFile_DB(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		// MySQL
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
		// SQLite
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

func TestRun_singleFile_TeX(t *testing.T) {
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
					{Stream: 1, Log: "This is pdfTeX, Version 3.141592653-2.6-1.40.22 (TeX Live 2021) (preloaded format=pdflatex)"},
					{Stream: 1, Log: " restricted \\write18 enabled."},
					{Stream: 1, Log: "entering extended mode"},
					{Stream: 1, Log: "(./runbox.tex"},
					{Stream: 1, Log: "LaTeX2e <2020-10-01> patch level 4"},
					{Stream: 1, Log: "L3 programming layer <2021-02-18>"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/latex/base/article.cls"},
					{Stream: 1, Log: "Document Class: article 2020/04/10 v1.4m Standard LaTeX document class"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/latex/base/size10.clo))"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/latex/geometry/geometry.sty"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/latex/graphics/keyval.sty)"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/generic/iftex/ifvtex.sty"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/generic/iftex/iftex.sty)))"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/latex/l3backend/l3backend-pdftex.def)"},
					{Stream: 1, Log: "No file runbox.aux."},
					{Stream: 1, Log: "*geometry* driver: auto-detecting"},
					{Stream: 1, Log: "*geometry* detected driver: pdftex"},
					{Stream: 1, Log: "[1{/usr/local/texlive/2021/texmf-var/fonts/map/pdftex/updmap/pdftex.map}]"},
					{Stream: 1, Log: "(./runbox.aux) )</usr/local/texlive/2021/texmf-dist/fonts/type1/public/amsfonts"},
					{Stream: 1, Log: "/cm/cmr10.pfb>"},
					{Stream: 1, Log: "Output written on runbox.pdf (1 page, 11915 bytes)."},
					{Stream: 1, Log: "Transcript written on runbox.log."},
				},
				Code:     0,
				CPU:      721795,
				MEM:      81156,
				Time:     739,
				Timedout: false,
				Images: []string{
					"iVBORw0KGgoAAAANSUhEUgAAAaQAAAEqCAQAAADw/+wWAAAE4klEQVR42u3X3Y0bRxpA0a8WTqBTYApMgSkwhUlBm4JCcAL74ElhQjBDMFNgCL0PkkY/gFcr6wLkCOc8NMkqdqMK6Ism1z7Az/rXvRcAvwIhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQeChQ1rH9ft6Wts6rj/W+evxT8cfv+Lr+6f1bmZt669775O376FD2i9zmOt+2y+z7c/fjH88/oMrfnKdbWa/zfXe++Tt++3eC/gR6zTbXPavbvx1nG2un8fWcba5zTbX2WbbXz7Nf3h9/dZhth/PEP7OQz+RZmbmuE7rNNvMOs9hf553X06u0xz3lzmv0+vQdZ72y5z369y+nN8v8/tsc5yZWcc575e5fTzjcu8t8vY9fkiX/WV/mdvMnOe2Tq+3/wdP8zwzl3n9B7XfZtY2t3Wew/7y1fx1v3z8gfg0l5nXK73ce4u8fY8f0me3ue4v8/6bscPH18+e5928n+NsfzP/7SdPJH7aQ4e0jrPNaW3rOId1nn/PaZ3mNLOOc1iHdZzD/GdO6zTHL/Pan2fbb/MhkPef5tdxDuv44cx5P6d1mOMc1ra2+fPe++TtW/u9VwC/gId+IsFbISQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkB7eerr3Cvg+IT24dZrzvdfA9wnpwe0v914B/w8hQUBIEBASBIT04NZpDuu8tnuvg/9t7fdeAfwCPJEgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAv8F6NR4np+BlBIAAAAASUVORK5CYII=",
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
					{Stream: 1, Log: "This is pdfTeX, Version 3.141592653-2.6-1.40.22 (TeX Live 2021) (preloaded format=pdflatex)"},
					{Stream: 1, Log: " restricted \\write18 enabled."},
					{Stream: 1, Log: "entering extended mode"},
					{Stream: 1, Log: "(./runbox.tex"},
					{Stream: 1, Log: "LaTeX2e <2020-10-01> patch level 4"},
					{Stream: 1, Log: "L3 programming layer <2021-02-18>"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/latex/base/article.cls"},
					{Stream: 1, Log: "Document Class: article 2020/04/10 v1.4m Standard LaTeX document class"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/latex/base/size10.clo))"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/latex/geometry/geometry.sty"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/latex/graphics/keyval.sty)"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/generic/iftex/ifvtex.sty"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/generic/iftex/iftex.sty)))"},
					{Stream: 1, Log: "(/usr/local/texlive/2021/texmf-dist/tex/latex/l3backend/l3backend-pdftex.def)"},
					{Stream: 1, Log: "No file runbox.aux."},
					{Stream: 1, Log: "*geometry* driver: auto-detecting"},
					{Stream: 1, Log: "*geometry* detected driver: pdftex"},
					{Stream: 1, Log: "[1{/usr/local/texlive/2021/texmf-var/fonts/map/pdftex/updmap/pdftex.map}]"},
					{Stream: 1, Log: "[2] [3] [4] [5] [6] [7] [8] [9] [10] [11] [12] (./runbox.aux) )</usr/local/texl"},
					{Stream: 1, Log: "ive/2021/texmf-dist/fonts/type1/public/amsfonts/cm/cmr10.pfb>"},
					{Stream: 1, Log: "Output written on runbox.pdf (12 pages, 16268 bytes)."},
					{Stream: 1, Log: "Transcript written on runbox.log."},
				},
				Code:     0,
				CPU:      1650328,
				MEM:      832,
				Time:     1160,
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

func TestRun_timeout(t *testing.T) {
	testCases := []struct {
		input Input
		want  *box.Result
	}{
		{
			Input{Lang: "bash", Files: []box.File{{Body: `echo hello; sleep 3`}}},
			&box.Result{
				Logs:     []box.Log{{Stream: 1, Log: "hello"}},
				Code:     0,
				CPU:      11040,
				MEM:      4648,
				Time:     2000,
				Timedout: true,
			},
		},
		{
			Input{Lang: "bash", Files: []box.File{{Body: `sleep 3; echo hello`}}},
			&box.Result{
				CPU:      10097,
				MEM:      788,
				Time:     2001,
				Timedout: true,
			},
		},
		{
			Input{Lang: "bash", Files: []box.File{{Body: `echo hello; echo world; sleep 3`}}},
			&box.Result{
				Logs: []box.Log{
					{Stream: 1, Log: "hello"},
					{Stream: 1, Log: "world"},
				},
				CPU:      9588,
				MEM:      796,
				Time:     2001,
				Timedout: true,
			},
		},
		{
			Input{Lang: "bash", Files: []box.File{{Body: `echo hello; sleep 3; echo world`}}},
			&box.Result{
				Logs:     []box.Log{{Stream: 1, Log: "hello"}},
				CPU:      9681,
				MEM:      804,
				Time:     2001,
				Timedout: true,
			},
		},
		{
			Input{Lang: "bash", Files: []box.File{{Body: `sleep 3; echo hello; echo world`}}},
			&box.Result{
				CPU:      9853,
				MEM:      800,
				Time:     2001,
				Timedout: true,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := lang1.Run(tc.input, map[string]int{"timeoutSeconds": 1})
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}
