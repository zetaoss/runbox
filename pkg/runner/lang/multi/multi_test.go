package multi

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/runner/lang/types"
)

func TestRun_Simple(t *testing.T) {
	testcases := []struct {
		lang       string
		files      []File
		wantOutput *types.Output
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
			&types.Output{Logs: []string{"0hello"}},
			"",
		},
		{
			"bash",
			[]File{
				{Name: "", Text: "cat greet.txt", Main: true},
				{Name: "greet.txt", Text: "hello", Main: false},
			},
			&types.Output{Logs: []string{"0hello"}},
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
		{"", File{Text: `echo hello`, Main: false}, "ErrInvalidLanguage"},
		{"X", File{}, "ErrInvalidLanguage"},
		{"X", File{Text: `echo hello`, Main: false}, "ErrInvalidLanguage"},
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
		wantOutput *types.Output
	}{
		// Bash
		{"bash", File{}, &types.Output{Logs: []string{}}},
		{"bash", File{Name: "", Text: `echo hello`, Main: false}, &types.Output{Logs: []string{"0hello"}}},
		{"bash", File{Name: "", Text: `echo hello; echo world`, Main: false}, &types.Output{Logs: []string{"0hello", "0world"}}},
		{"bash", File{Name: "", Text: `echo hello; echo world; echo`, Main: false}, &types.Output{Logs: []string{"0hello", "0world", "0"}}},
		{"bash", File{Name: "", Text: `echo hello 1>&2; echo world 1>&2`, Main: false}, &types.Output{Logs: []string{"1hello", "1world"}}},
		// C
		{"c", File{Name: "", Text: `
#include <stdio.h>
int main() {
	printf("Hello, World!");
}
`, Main: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		{"c", File{Name: "", Text: `
#include <stdio.h>
int main() {
	printf("Hello\nWorld!");
}
`, Main: false}, &types.Output{Logs: []string{"0Hello", "0World!"}}},
		// C++
		{"cpp", File{Name: "", Text: `
#include <iostream>
int main() {
	std::cout<<"hello";
}
`, Main: false}, &types.Output{Logs: []string{"0hello"}}},
		// C#
		{"csharp", File{Name: "", Text: `
using System;
class Program
{
	static void Main() {
		Console.Write("hello");
	}
}
`, Main: false}, &types.Output{Logs: []string{"0hello"}}},
		// Java
		{"java", File{Name: "", Text: `
public class App {
	public static void main(String args[]) {
		System.out.println("hello");
	}
}
`, Main: false}, &types.Output{Logs: []string{"0hello"}}},
		{"java", File{Name: "", Text: `
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
}
`, Main: false}, &types.Output{
			Logs:   []string{},
			Images: []string{"iVBORw0KGgoAAAANSUhEUgAAASwAAADICAIAAADdvUsCAAACaUlEQVR4Xu3TQW4bMRAAQf3/08pBwGLD4VqyEbiloOpgkLMUfWHf7kDqtg6A3yVCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYjcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP6l+/2+LH7T1T+9msO7m293ThY/iHA5+fX2a1eHr+bw7ubbnZPFDyK8/X34av2Kq/NXc3h38+0uhRzmgacnz7Y3zPW85LE+hlfnz3P4JPPtHpPzy17Wc7Fdny2/ev2S5cLtfDkDn+TxfBfHp+Xk08V2e3jM59/z17n97hw+zHy751e+mAeWxXZ7eMzn3/PXuf3uHD7MfLtXr/wwD7zYw8zvfPLqku/O4cPMt3t+5dtItovtetpeeGy3l8xjx2K7hg8z3+558njcyxM/1k9PTtt7zpN5yXJs+bT9CQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADwv/sDs/dz4IQQs5EAAAAASUVORK5CYII="},
		}},
		// Kotlin
		{"kotlin", File{Name: "", Text: `
fun main() {
	println("Hello, World!")
}
`, Main: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
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
		wantOutput *types.Output
	}{

		// Go
		{"go", File{Name: "", Text: `
package main
import "fmt"
func main() {
	fmt.Println("Hello, 世界")
}
`, Main: false}, &types.Output{Logs: []string{"0Hello, 世界"}}},
		// Lua
		{"lua", File{Name: "", Text: `print("Hello, World!")`, Main: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		// MySQL
		{"mysql", File{Name: "", Text: `SELECT 'Hello, World!';`, Main: false}, &types.Output{Logs: []string{
			"0+---------------+",
			"0| Hello, World! |",
			"0+---------------+",
			"0| Hello, World! |",
			"0+---------------+",
		}}},
		// Perl
		{"perl", File{Name: "", Text: `` +
			"\n" + `use strict;` +
			"\n" + `use warnings;` +
			"\n" + `print "Hello, World!\n";`},
			&types.Output{Logs: []string{"0Hello, World!"}}},
		// PHP
		{"php", File{Name: "", Main: false, Text: `echo "Hello, World!";`}, &types.Output{Logs: []string{"0Hello, World!"}}},
		{"php", File{Name: "", Main: false, Text: `` +
			"\n" + `<?php` +
			"\n" + `echo "Hello, World!";`},
			&types.Output{Logs: []string{"0Hello, World!"}}},
		// PowerShell
		{"powershell", File{Name: "", Text: `Write-Host "Hello, World!"`, Main: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		// Python
		{"python", File{Name: "", Text: `print("Hello, World!")`, Main: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		// R
		{"r", File{Name: "", Text: `print("Hello, World!")`, Main: false}, &types.Output{Logs: []string{`0[1] "Hello, World!"`}}},
		// Ruby
		{"ruby", File{Name: "", Text: `print("Hello, World!")`, Main: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		// SQLite
		{"sqlite3", File{Name: "", Text: `SELECT 'Hello, World!';`, Main: false}, &types.Output{Logs: []string{
			"0+-----------------+",
			"0| 'Hello, World!' |",
			"0+-----------------+",
			"0| Hello, World!   |",
			"0+-----------------+",
		}}},
		{"sqlite3", File{Name: "", Text: `.tables`, Main: false}, &types.Output{Logs: []string{
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
		wantOutput *types.Output
	}{
		// Tex
		{"tex", File{Name: "", Text: `` +
			"\n" + `\documentclass{article}` +
			"\n" + `\usepackage[a6paper,landscape]{geometry}` +
			"\n" + `\begin{document}` +
			"\n" + `Hello world!` +
			"\n" + `\end{document}`},
			&types.Output{Logs: []string{"0Hello, World!"}},
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

func TestRun_OutputLimitReached(t *testing.T) {
	testcases := []struct {
		lang        string
		file        File
		wantOutput  *types.Output
		wantLengths []int
	}{
		{"python", File{Name: "", Text: `print(100*"X")`, Main: false}, &types.Output{}, []int{101}},
		{"python", File{Name: "", Text: `print(1000*"X")`, Main: false}, &types.Output{}, []int{1001}},
		{"python", File{Name: "", Text: `print(10000*"X")`, Main: false}, &types.Output{Warnings: []string{types.WarnOutputLimitReached}}, []int{8001}},
		{"python", File{Name: "", Text: `[print(1000*"X") for _ in range(1)]`, Main: false}, &types.Output{}, []int{1001}},
		{"python", File{Name: "", Text: `[print(1000*"X") for _ in range(10)]`, Main: false}, &types.Output{Warnings: []string{types.WarnOutputLimitReached}}, []int{1001, 1001, 1001, 1001, 1001, 1001, 1001, 994}},
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

func TestRunWithTimeout(t *testing.T) {
	testCases := []struct {
		lang       string
		file       File
		wantOutput *types.Output
		wantError  string
	}{
		// Bash
		{
			"bash", File{Name: "", Text: `echo hello`, Main: false},
			&types.Output{Logs: []string{"0hello"}}, "",
		},
		{
			"bash", File{Name: "", Text: `echo hello; sleep 3`, Main: false},
			&types.Output{Logs: []string{"0hello"}, Warnings: []string{types.WarnTimeout}}, "",
		},
		{
			"bash", File{Name: "", Text: `sleep 3; echo hello`, Main: false},
			&types.Output{Logs: []string{}, Warnings: []string{types.WarnTimeout}}, "",
		},
		{
			"bash", File{Name: "", Text: `echo hello; echo world`, Main: false},
			&types.Output{Logs: []string{"0hello", "0world"}}, "",
		},
		{
			"bash", File{Name: "", Text: `echo hello; echo world; sleep 3`, Main: false},
			&types.Output{Logs: []string{"0hello", "0world"}, Warnings: []string{types.WarnTimeout}}, "",
		},
		{
			"bash", File{Name: "", Text: `echo hello; sleep 3; echo world`, Main: false},
			&types.Output{Logs: []string{"0hello"}, Warnings: []string{types.WarnTimeout}}, "",
		},
		{
			"bash", File{Name: "", Text: `sleep 3; echo hello; echo world`, Main: false},
			&types.Output{Logs: []string{}, Warnings: []string{types.WarnTimeout}}, "",
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
