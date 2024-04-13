package lang

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun_simple(t *testing.T) {
	testcases := []struct {
		lang       string
		files      []File
		wantOutput *Output
		wantError  string
	}{
		{
			"bash",
			[]File{},
			nil,
			"ErrNoFiles",
		},
		{
			"bash",
			[]File{
				{Name: "greet.txt", Text: "hello", Main: false},
				{Name: "", Text: "cat greet.txt", Main: true},
			},
			&Output{Logs: []string{"0hello"}},
			"",
		},
		{
			"bash",
			[]File{
				{Name: "", Text: "cat greet.txt", Main: true},
				{Name: "greet.txt", Text: "hello", Main: false},
			},
			&Output{Logs: []string{"0hello"}},
			"",
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			input := Input{Lang: tc.lang, Files: tc.files}
			output, err := Run(input)
			if tc.wantError == "" {
				require.NoError(t, err)
				// ignore fields
				output.Time = ""
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.wantOutput, output)
		})
	}
}

func TestRun_invalid_language(t *testing.T) {
	testcases := []struct {
		lang      string
		file      File
		wantError string
	}{
		{"", File{}, "ErrInvalidLanguage"},
		{"", File{Text: `echo hello`}, "ErrInvalidLanguage"},
		{"X", File{}, "ErrInvalidLanguage"},
		{"X", File{Text: `echo hello`}, "ErrInvalidLanguage"},
	}
	for _, tc := range testcases {
		t.Run(tc.lang+"__", func(t *testing.T) {
			input := Input{Lang: tc.lang, Files: []File{tc.file}}
			output, err := Run(input)
			require.Nil(t, output)
			require.EqualError(t, err, tc.wantError)
		})
	}
}

func TestRun_part1(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	testcases := []struct {
		lang       string
		file       File
		wantOutput *Output
	}{
		// Bash
		{"bash", File{}, &Output{Logs: []string{}}},
		{"bash", File{Text: `echo hello`}, &Output{Logs: []string{"0hello"}}},
		{"bash", File{Text: `echo hello; echo world`}, &Output{Logs: []string{"0hello", "0world"}}},
		{"bash", File{Text: `echo hello; echo world; echo`}, &Output{Logs: []string{"0hello", "0world", "0"}}},
		{"bash", File{Text: `echo hello 1>&2; echo world 1>&2`}, &Output{Logs: []string{"1hello", "1world"}}},
		// C
		{"c", File{Text: "" +
			"\n" + `#include <stdio.h>` +
			"\n" + `int main() {` +
			"\n" + `	printf("Hello, World!");` +
			"\n" + `}` +
			"\n"}, &Output{Logs: []string{"0Hello, World!"}}},
		{"c", File{Text: "" +
			"\n" + `#include <stdio.h>` +
			"\n" + `int main() {` +
			"\n" + `	printf("Hello\nWorld!");` +
			"\n" + `}` +
			"\n"}, &Output{Logs: []string{"0Hello", "0World!"}}},
		// C++
		{"cpp", File{Text: "" +
			"\n" + `#include <iostream>` +
			"\n" + `int main() {` +
			"\n" + `	std::cout<<"hello";` +
			"\n" + `}` +
			"\n"}, &Output{Logs: []string{"0hello"}}},
		// C#
		{"csharp", File{Text: "" +
			"\n" + `using System;` +
			"\n" + `class Program` +
			"\n" + `{` +
			"\n" + `	static void Main() {` +
			"\n" + `		Console.Write("hello");` +
			"\n" + `	}` +
			"\n" + `}` +
			"\n"}, &Output{Logs: []string{"0hello"}}},
		// Java
		{"java", File{Text: "" +
			"\n" + `public class App {` +
			"\n" + `	public static void main(String args[]) {` +
			"\n" + `		System.out.println("hello");` +
			"\n" + `	}` +
			"\n" + `}` +
			"\n"}, &Output{Logs: []string{"0hello"}}},
		{"java", File{Text: "" +
			"\n" + `import java.awt.Graphics2D;` +
			"\n" + `import java.awt.image.BufferedImage;` +
			"\n" + `import java.io.File;` +
			"\n" + `import java.io.IOException;` +
			"\n" + `import javax.imageio.ImageIO;` +
			"\n" + `` +
			"\n" + `public class App {` +
			"\n" + `	public static void main(String[] args) throws IOException {` +
			"\n" + `		BufferedImage bufferedImage = new BufferedImage(300, 200, BufferedImage.TYPE_INT_RGB);` +
			"\n" + `		Graphics2D g = bufferedImage.createGraphics();` +
			"\n" + `		g.fillRect(0, 0, 300, 100);` +
			"\n" + `		g.drawString("Hello World", 120, 150);` +
			"\n" + `		g.dispose();` +
			"\n" + `		ImageIO.write(bufferedImage, "png", new File("myimage.png"));` +
			"\n" + `	}` +
			"\n" + `}` +
			"\n"}, &Output{
			Logs:   []string{},
			Images: []string{"iVBORw0KGgoAAAANSUhEUgAAASwAAADICAIAAADdvUsCAAACaUlEQVR4Xu3TQW4bMRAAQf3/08pBwGLD4VqyEbiloOpgkLMUfWHf7kDqtg6A3yVCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYjcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP6l+/2+LH7T1T+9msO7m293ThY/iHA5+fX2a1eHr+bw7ubbnZPFDyK8/X34av2Kq/NXc3h38+0uhRzmgacnz7Y3zPW85LE+hlfnz3P4JPPtHpPzy17Wc7Fdny2/ev2S5cLtfDkDn+TxfBfHp+Xk08V2e3jM59/z17n97hw+zHy751e+mAeWxXZ7eMzn3/PXuf3uHD7MfLtXr/wwD7zYw8zvfPLqku/O4cPMt3t+5dtItovtetpeeGy3l8xjx2K7hg8z3+558njcyxM/1k9PTtt7zpN5yXJs+bT9CQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADwv/sDs/dz4IQQs5EAAAAASUVORK5CYII="},
		}},
		// Kotlin
		{"kotlin", File{Text: "" +
			"\n" + `fun main() {` +
			"\n" + `	println("Hello, World!")` +
			"\n" + `}` +
			"\n"}, &Output{Logs: []string{"0Hello, World!"}}},
	}
	for _, tc := range testcases {
		t.Run(tc.lang+"__", func(t *testing.T) {
			input := Input{Lang: tc.lang, Files: []File{tc.file}}
			output, err := Run(input)
			require.NoError(t, err)

			// ignore fields
			output.Time = ""
			output.CPU = 0
			output.MEM = 0
			require.Equal(t, tc.wantOutput, output)
		})
	}
}

func TestRun_part2(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	testcases := []struct {
		lang       string
		file       File
		wantOutput *Output
	}{
		// Go
		{"go", File{Text: "" +
			"\n" + `package main` +
			"\n" + `import "fmt"` +
			"\n" + `func main() {` +
			"\n" + `	fmt.Println("Hello, 世界")` +
			"\n" + `}` +
			"\n"}, &Output{Logs: []string{"0Hello, 世界"}}},
		// Lua
		{"lua", File{Text: `print("Hello, World!")`}, &Output{Logs: []string{"0Hello, World!"}}},
		// MySQL
		{"mysql", File{Text: `SELECT 'Hello, World!';`}, &Output{Logs: []string{
			"0+---------------+",
			"0| Hello, World! |",
			"0+---------------+",
			"0| Hello, World! |",
			"0+---------------+",
		}}},
		// Perl
		{"perl", File{Text: "" +
			"\n" + `use strict;` +
			"\n" + `use warnings;` +
			"\n" + `print "Hello, World!\n";`},
			&Output{Logs: []string{"0Hello, World!"}}},
		// PHP
		{"php", File{Main: false, Text: `echo "Hello, World!";`}, &Output{Logs: []string{"0Hello, World!"}}},
		{"php", File{Main: false, Text: "" +
			"\n" + `<?php` +
			"\n" + `echo "Hello, World!";`},
			&Output{Logs: []string{"0Hello, World!"}}},
		// PowerShell
		{"powershell", File{Text: `Write-Host "Hello, World!"`}, &Output{Logs: []string{"0Hello, World!"}}},
		// Python
		{"python", File{Text: `print("Hello, World!")`}, &Output{Logs: []string{"0Hello, World!"}}},
		// R
		{"r", File{Text: `print("Hello, World!")`}, &Output{Logs: []string{`0[1] "Hello, World!"`}}},
		// Ruby
		{"ruby", File{Text: `print("Hello, World!")`}, &Output{Logs: []string{"0Hello, World!"}}},
		// SQLite
		{"sqlite3", File{Text: `SELECT 'Hello, World!';`}, &Output{Logs: []string{
			"0+-----------------+",
			"0| 'Hello, World!' |",
			"0+-----------------+",
			"0| Hello, World!   |",
			"0+-----------------+",
		}}},
		{"sqlite3", File{Text: `.tables`}, &Output{Logs: []string{
			"0Category              EmployeeTerritory     Region              ",
			"0Customer              Order                 Shipper             ",
			"0CustomerCustomerDemo  OrderDetail           Supplier            ",
			"0CustomerDemographic   Product               Territory           ",
			"0Employee              ProductDetails_V    ",
		}}},
	}
	for _, tc := range testcases {
		t.Run(tc.lang+"__", func(t *testing.T) {
			input := Input{Lang: tc.lang, Files: []File{tc.file}}
			output, err := Run(input)
			require.NoError(t, err)

			// ignore fields
			output.Time = ""
			output.CPU = 0
			output.MEM = 0
			require.Equal(t, tc.wantOutput, output)
		})
	}
}

func TestRun_tex(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	testcases := []struct {
		lang       string
		file       File
		wantOutput *Output
	}{
		// latex
		{
			"latex", File{Text: "" +
				"\n" + `\begin{document}` +
				"\n" + `Hello World!` +
				"\n" + `\end{document}`,
			},
			&Output{
				Logs: []string{
					"0This is pdfTeX, Version 3.141592653-2.6-1.40.22 (TeX Live 2021) (preloaded format=pdflatex)",
					"0 restricted \\write18 enabled.",
					"0entering extended mode",
					"0(./runbox.tex",
					"0LaTeX2e <2020-10-01> patch level 4",
					"0L3 programming layer <2021-02-18>",
					"0(/usr/local/texlive/2021/texmf-dist/tex/latex/l3backend/l3backend-pdftex.def)",
					"0No file runbox.aux.",
					"0",
					"0! LaTeX Error: The font size command \\normalsize is not defined:",
					"0               there is probably something wrong with the class file.",
					"0",
					"0See the LaTeX manual or LaTeX Companion for explanation.",
					"0Type  H <return>  for immediate help.",
					"0 ...                                              ",
					"0                                                  ",
					"0l.2 \\begin{document}",
					"0                    ",
					"0!  ==> Fatal error occurred, no output PDF file produced!",
					"0Transcript written on runbox.log.",
				},
			},
		},
		{
			"latex", File{Text: "" +
				"\n" + `\documentclass{minimal}` +
				"\n" + `\begin{document}` +
				"\n" + `Hello World!` +
				"\n" + `\end{document}`,
			},
			&Output{
				Logs:   []string{},
				Images: []string{"iVBORw0KGgoAAAANSUhEUgAAAlMAAANKCAQAAAAE5gOEAAAMUElEQVR42u3YwW3jRhiA0X/SgToI1IJaYAtuwS04JWRL2AZyiEuIS4hKiFtQCcwhtndtGLkEWX+C3zuQ0mg44Fw+iFz7AJT99NE3APDvZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIu9pMrdP6um7XYZ3W7+vm9fjz8Z1r/lindVh36+s6zqy7dffeui+fb9fdOqy/Pnqv8Lldbab28xzncb/s5zns92/Gn47vXHOZy36Z+9n2x5k571/eXffZ4xz2yzx+9F7hc7vaTL21tnWz3qRpndb2ZuxhbmbmOOd1ej3nn/O32eu4Tu+lDvjRrjtTp7WtbQ4z62aO+/28eoRb25z2h7lZ23eD/2Tq5fxtzn6er3OYl3jNzX6ey8zMnD96m/C5XXemzvvD/jCXmbmZy9qesvLsdu5n5jzfvbnaH2fWcWbuZ1vb/vBqzuN+fnl8vJ3zzNN6Dx+9TfjcrjtT31zmcX+YL2/Gjk/n793Pr/vDfnn69f05r7/7NwUf6moztU5zmG0d1mmO62Z+mW1ts82s0xyf3iv9Ntva5vQmXvcv/5EuM/Plec46zXGdnq+fL7Ot45zmuH6ePz96r/C5rf2j7+DHb/m4P86sw37572sB/79PmCngulztQx/wWcgUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUEPc3T4eA1oSjquEAAAAASUVORK5CYII="},
			},
		},
		{
			"latex", File{Text: "" +
				"\n" + `\documentclass{article}` +
				"\n" + `\usepackage[a6paper,landscape]{geometry}` +
				"\n" + `\begin{document}` +
				"\n" + `Hello world!` +
				"\n" + `\end{document}`,
			},
			&Output{
				Logs:   []string{},
				Images: []string{"iVBORw0KGgoAAAANSUhEUgAAAaQAAAEqCAQAAADw/+wWAAAE4klEQVR42u3X3Y0bRxpA0a8WTqBTYApMgSkwhUlBm4JCcAL74ElhQjBDMFNgCL0PkkY/gFcr6wLkCOc8NMkqdqMK6Ism1z7Az/rXvRcAvwIhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQeChQ1rH9ft6Wts6rj/W+evxT8cfv+Lr+6f1bmZt669775O376FD2i9zmOt+2y+z7c/fjH88/oMrfnKdbWa/zfXe++Tt++3eC/gR6zTbXPavbvx1nG2un8fWcba5zTbX2WbbXz7Nf3h9/dZhth/PEP7OQz+RZmbmuE7rNNvMOs9hf553X06u0xz3lzmv0+vQdZ72y5z369y+nN8v8/tsc5yZWcc575e5fTzjcu8t8vY9fkiX/WV/mdvMnOe2Tq+3/wdP8zwzl3n9B7XfZtY2t3Wew/7y1fx1v3z8gfg0l5nXK73ce4u8fY8f0me3ue4v8/6bscPH18+e5928n+NsfzP/7SdPJH7aQ4e0jrPNaW3rOId1nn/PaZ3mNLOOc1iHdZzD/GdO6zTHL/Pan2fbb/MhkPef5tdxDuv44cx5P6d1mOMc1ra2+fPe++TtW/u9VwC/gId+IsFbISQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkB7eerr3Cvg+IT24dZrzvdfA9wnpwe0v914B/w8hQUBIEBASBIT04NZpDuu8tnuvg/9t7fdeAfwCPJEgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAv8F6NR4np+BlBIAAAAASUVORK5CYII="},
			},
		},
		// tex
		{
			"tex", File{Text: "" +
				"\n" + `\begin{document}` +
				"\n" + `Hello World!` +
				"\n" + `\end{document}`,
			},
			&Output{
				Logs: []string{
					"0This is pdfTeX, Version 3.141592653-2.6-1.40.22 (TeX Live 2021) (preloaded format=pdflatex)",
					"0 restricted \\write18 enabled.",
					"0entering extended mode",
					"0(./runbox.tex",
					"0LaTeX2e <2020-10-01> patch level 4",
					"0L3 programming layer <2021-02-18>",
					"0(/usr/local/texlive/2021/texmf-dist/tex/latex/l3backend/l3backend-pdftex.def)",
					"0No file runbox.aux.",
					"0",
					"0! LaTeX Error: The font size command \\normalsize is not defined:",
					"0               there is probably something wrong with the class file.",
					"0",
					"0See the LaTeX manual or LaTeX Companion for explanation.",
					"0Type  H <return>  for immediate help.",
					"0 ...                                              ",
					"0                                                  ",
					"0l.2 \\begin{document}",
					"0                    ",
					"0!  ==> Fatal error occurred, no output PDF file produced!",
					"0Transcript written on runbox.log.",
				},
			},
		},
		{
			"tex", File{Text: "" +
				"\n" + `\documentclass{minimal}` +
				"\n" + `\begin{document}` +
				"\n" + `Hello World!` +
				"\n" + `\end{document}`,
			},
			&Output{
				Logs:   []string{},
				Images: []string{"iVBORw0KGgoAAAANSUhEUgAAAlMAAANKCAQAAAAE5gOEAAAMUElEQVR42u3YwW3jRhiA0X/SgToI1IJaYAtuwS04JWRL2AZyiEuIS4hKiFtQCcwhtndtGLkEWX+C3zuQ0mg44Fw+iFz7AJT99NE3APDvZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIkykgTqaAOJkC4mQKiJMpIE6mgDiZAuJkCoiTKSBOpoA4mQLiZAqIu9pMrdP6um7XYZ3W7+vm9fjz8Z1r/lindVh36+s6zqy7dffeui+fb9fdOqy/Pnqv8Lldbab28xzncb/s5zns92/Gn47vXHOZy36Z+9n2x5k571/eXffZ4xz2yzx+9F7hc7vaTL21tnWz3qRpndb2ZuxhbmbmOOd1ej3nn/O32eu4Tu+lDvjRrjtTp7WtbQ4z62aO+/28eoRb25z2h7lZ23eD/2Tq5fxtzn6er3OYl3jNzX6ey8zMnD96m/C5XXemzvvD/jCXmbmZy9qesvLsdu5n5jzfvbnaH2fWcWbuZ1vb/vBqzuN+fnl8vJ3zzNN6Dx+9TfjcrjtT31zmcX+YL2/Gjk/n793Pr/vDfnn69f05r7/7NwUf6moztU5zmG0d1mmO62Z+mW1ts82s0xyf3iv9Ntva5vQmXvcv/5EuM/Plec46zXGdnq+fL7Ot45zmuH6ePz96r/C5rf2j7+DHb/m4P86sw37572sB/79PmCngulztQx/wWcgUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUECdTQJxMAXEyBcTJFBAnU0CcTAFxMgXEyRQQJ1NAnEwBcTIFxMkUEPc3T4eA1oSjquEAAAAASUVORK5CYII="},
			},
		},
		{
			"tex", File{Text: "" +
				"\n" + `\documentclass{article}` +
				"\n" + `\usepackage[a6paper,landscape]{geometry}` +
				"\n" + `\begin{document}` +
				"\n" + `Hello world!` +
				"\n" + `\end{document}`,
			},
			&Output{
				Logs:   []string{},
				Images: []string{"iVBORw0KGgoAAAANSUhEUgAAAaQAAAEqCAQAAADw/+wWAAAE4klEQVR42u3X3Y0bRxpA0a8WTqBTYApMgSkwhUlBm4JCcAL74ElhQjBDMFNgCL0PkkY/gFcr6wLkCOc8NMkqdqMK6Ism1z7Az/rXvRcAvwIhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQeChQ1rH9ft6Wts6rj/W+evxT8cfv+Lr+6f1bmZt669775O376FD2i9zmOt+2y+z7c/fjH88/oMrfnKdbWa/zfXe++Tt++3eC/gR6zTbXPavbvx1nG2un8fWcba5zTbX2WbbXz7Nf3h9/dZhth/PEP7OQz+RZmbmuE7rNNvMOs9hf553X06u0xz3lzmv0+vQdZ72y5z369y+nN8v8/tsc5yZWcc575e5fTzjcu8t8vY9fkiX/WV/mdvMnOe2Tq+3/wdP8zwzl3n9B7XfZtY2t3Wew/7y1fx1v3z8gfg0l5nXK73ce4u8fY8f0me3ue4v8/6bscPH18+e5928n+NsfzP/7SdPJH7aQ4e0jrPNaW3rOId1nn/PaZ3mNLOOc1iHdZzD/GdO6zTHL/Pan2fbb/MhkPef5tdxDuv44cx5P6d1mOMc1ra2+fPe++TtW/u9VwC/gId+IsFbISQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkB7eerr3Cvg+IT24dZrzvdfA9wnpwe0v914B/w8hQUBIEBASBIT04NZpDuu8tnuvg/9t7fdeAfwCPJEgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAv8F6NR4np+BlBIAAAAASUVORK5CYII="},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.lang+"__", func(t *testing.T) {
			input := Input{Lang: tc.lang, Files: []File{tc.file}}
			output, err := Run(input)
			require.NoError(t, err)
			// ignore fields
			output.Time = ""
			output.CPU = 0
			output.MEM = 0
			require.Equal(t, tc.wantOutput, output)
		})
	}
}

func TestRun_outputLimitReached(t *testing.T) {
	testcases := []struct {
		lang        string
		file        File
		wantOutput  *Output
		wantLengths []int
	}{
		{"python", File{Text: `print(100*"X")`}, &Output{}, []int{101}},
		{"python", File{Text: `print(1000*"X")`}, &Output{}, []int{1001}},
		{"python", File{Text: `print(10000*"X")`}, &Output{Warnings: []string{WarnOutputLimitReached}}, []int{8001}},
		{"python", File{Text: `[print(1000*"X") for _ in range(1)]`}, &Output{}, []int{1001}},
		{"python", File{Text: `[print(1000*"X") for _ in range(10)]`}, &Output{Warnings: []string{WarnOutputLimitReached}}, []int{1001, 1001, 1001, 1001, 1001, 1001, 1001, 994}},
	}
	for _, tc := range testcases {
		t.Run(tc.lang+"__", func(t *testing.T) {
			input := Input{Lang: tc.lang, Files: []File{tc.file}}
			output, err := Run(input)
			require.NoError(t, err)

			lengths := []int{}
			for _, l := range output.Logs {
				lengths = append(lengths, len(l))
			}

			// ignore fields
			output.Logs = nil
			output.Time = ""
			output.CPU = 0
			output.MEM = 0

			require.Equal(t, tc.wantOutput, output)
			require.Equal(t, tc.wantLengths, lengths)
		})
	}
}

func TestRun_timeout(t *testing.T) {
	testCases := []struct {
		lang       string
		file       File
		wantOutput *Output
		wantError  string
	}{
		// Bash
		{
			"bash", File{Text: `echo hello`},
			&Output{Logs: []string{"0hello"}}, "",
		},
		{
			"bash", File{Text: `echo hello; sleep 3`},
			&Output{Logs: []string{"0hello"}, Warnings: []string{WarnTimeout}}, "",
		},
		{
			"bash", File{Text: `sleep 3; echo hello`},
			&Output{Logs: []string{}, Warnings: []string{WarnTimeout}}, "",
		},
		{
			"bash", File{Text: `echo hello; echo world`},
			&Output{Logs: []string{"0hello", "0world"}}, "",
		},
		{
			"bash", File{Text: `echo hello; echo world; sleep 3`},
			&Output{Logs: []string{"0hello", "0world"}, Warnings: []string{WarnTimeout}}, "",
		},
		{
			"bash", File{Text: `echo hello; sleep 3; echo world`},
			&Output{Logs: []string{"0hello"}, Warnings: []string{WarnTimeout}}, "",
		},
		{
			"bash", File{Text: `sleep 3; echo hello; echo world`},
			&Output{Logs: []string{}, Warnings: []string{WarnTimeout}}, "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.lang, func(t *testing.T) {
			input := Input{Lang: tc.lang, Files: []File{tc.file}}
			output, err := Run(input, map[string]int{"timeoutSeconds": 1})
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			// ignore fields
			output.Time = ""
			output.CPU = 0
			output.MEM = 0
			require.Equal(t, tc.wantOutput, output)
		})
	}
}
