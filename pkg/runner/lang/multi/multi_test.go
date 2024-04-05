package multi

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/runner/lang/types"
)

func TestRun_Simple(t *testing.T) {
	testcases := []struct {
		lang       string
		files      []types.File
		wantOutput *types.Output
		wantError  string
	}{
		{
			"bash",
			[]types.File{},
			nil,
			"getRunOpts err: no files",
		},
		{
			"bash",
			[]types.File{
				{Name: "greet.txt", Content: "hello", IsMain: false},
				{Name: "", Content: "cat greet.txt", IsMain: true},
			},
			&types.Output{Logs: []string{"0hello"}},
			"",
		},
	}
	for _, tc := range testcases {
		t.Run("", func(t *testing.T) {
			input := types.MultiInput{Lang: tc.lang, Files: tc.files}
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
		file      types.File
		wantError string
	}{
		{"", types.File{}, "invalid language"},
		{"", types.File{Content: `echo hello`, IsMain: false}, "invalid language"},
		{"X", types.File{}, "invalid language"},
		{"X", types.File{Content: `echo hello`, IsMain: false}, "invalid language"},
	}
	for _, tc := range testcases {
		t.Run(tc.lang+"__", func(t *testing.T) {
			input := types.MultiInput{Lang: tc.lang, Files: []types.File{tc.file}}
			output, err := Run(input)
			require.Nil(t, output)
			require.EqualError(t, err, tc.wantError)
		})
	}
}

func TestRun_SingleFile(t *testing.T) {
	testcases := []struct {
		lang       string
		file       types.File
		wantOutput *types.Output
	}{
		// Bash
		{"bash", types.File{}, &types.Output{Logs: []string{}}},
		{"bash", types.File{Name: "", Content: `echo hello`, IsMain: false}, &types.Output{Logs: []string{"0hello"}}},
		{"bash", types.File{Name: "", Content: `echo hello; echo world`, IsMain: false}, &types.Output{Logs: []string{"0hello", "0world"}}},
		{"bash", types.File{Name: "", Content: `echo hello; echo world; echo`, IsMain: false}, &types.Output{Logs: []string{"0hello", "0world", "0"}}},
		{"bash", types.File{Name: "", Content: `echo hello 1>&2; echo world 1>&2`, IsMain: false}, &types.Output{Logs: []string{"1hello", "1world"}}},
		// C
		{"c", types.File{Name: "", Content: `
#include <stdio.h>
int main() {
	printf("Hello, World!");
}
`, IsMain: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		{"c", types.File{Name: "", Content: `
#include <stdio.h>
int main() {
	printf("Hello\nWorld!");
}
`, IsMain: false}, &types.Output{Logs: []string{"0Hello", "0World!"}}},
		// C++
		{"cpp", types.File{Name: "", Content: `
#include <iostream>
int main() {
	std::cout<<"hello";
}
`, IsMain: false}, &types.Output{Logs: []string{"0hello"}}},
		// C#
		{"csharp", types.File{Name: "", Content: `
using System;
class Program
{
	static void Main() {
		Console.Write("hello");
	}
}
`, IsMain: false}, &types.Output{Logs: []string{"0hello"}}},
		// Java
		{"java", types.File{Name: "", Content: `
public class App {
	public static void main(String args[]) {
		System.out.println("hello");
	}
}
`, IsMain: false}, &types.Output{Logs: []string{"0hello"}}},
		{"java", types.File{Name: "", Content: `
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
`, IsMain: false}, &types.Output{
			Logs:   []string{},
			Images: []string{"iVBORw0KGgoAAAANSUhEUgAAASwAAADICAIAAADdvUsCAAACaUlEQVR4Xu3TQW4bMRAAQf3/08pBwGLD4VqyEbiloOpgkLMUfWHf7kDqtg6A3yVCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYjcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP6l+/2+LH7T1T+9msO7m293ThY/iHA5+fX2a1eHr+bw7ubbnZPFDyK8/X34av2Kq/NXc3h38+0uhRzmgacnz7Y3zPW85LE+hlfnz3P4JPPtHpPzy17Wc7Fdny2/ev2S5cLtfDkDn+TxfBfHp+Xk08V2e3jM59/z17n97hw+zHy751e+mAeWxXZ7eMzn3/PXuf3uHD7MfLtXr/wwD7zYw8zvfPLqku/O4cPMt3t+5dtItovtetpeeGy3l8xjx2K7hg8z3+558njcyxM/1k9PTtt7zpN5yXJs+bT9CQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADwv/sDs/dz4IQQs5EAAAAASUVORK5CYII="},
		}},
		// Kotlin
		{"kotlin", types.File{Name: "", Content: `
fun main() {
	println("Hello, World!")
}
`, IsMain: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		// Go
		{"go", types.File{Name: "", Content: `
package main
import "fmt"
func main() {
	fmt.Println("Hello, 世界")
}
`, IsMain: false}, &types.Output{Logs: []string{"0Hello, 世界"}}},
		// Lua
		{"lua", types.File{Name: "", Content: `print("Hello, World!")`, IsMain: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		// MySQL
		{"mysql", types.File{Name: "", Content: `SELECT 'Hello, World!';`, IsMain: false}, &types.Output{Logs: []string{
			"0+---------------+",
			"0| Hello, World! |",
			"0+---------------+",
			"0| Hello, World! |",
			"0+---------------+",
		}}},
		// Perl
		{"perl", types.File{Name: "", Content: `` +
			"\n" + `use strict;` +
			"\n" + `use warnings;` +
			"\n" + `print "Hello, World!\n";`},
			&types.Output{Logs: []string{"0Hello, World!"}}},
		// PHP
		{"php", types.File{Name: "", IsMain: false, Content: `echo "Hello, World!";`}, &types.Output{Logs: []string{"0Hello, World!"}}},
		{"php", types.File{Name: "", IsMain: false, Content: `` +
			"\n" + `<?php` +
			"\n" + `echo "Hello, World!";`},
			&types.Output{Logs: []string{"0Hello, World!"}}},
		// PowerShell
		{"powershell", types.File{Name: "", Content: `Write-Host "Hello, World!"`, IsMain: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		// Python
		{"python", types.File{Name: "", Content: `print("Hello, World!")`, IsMain: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		// R
		{"r", types.File{Name: "", Content: `print("Hello, World!")`, IsMain: false}, &types.Output{Logs: []string{`0[1] "Hello, World!"`}}},
		// Ruby
		{"ruby", types.File{Name: "", Content: `print("Hello, World!")`, IsMain: false}, &types.Output{Logs: []string{"0Hello, World!"}}},
		// SQLite
		{"sqlite3", types.File{Name: "", Content: `SELECT 'Hello, World!';`, IsMain: false}, &types.Output{Logs: []string{
			"0+-----------------+",
			"0| 'Hello, World!' |",
			"0+-----------------+",
			"0| Hello, World!   |",
			"0+-----------------+",
		}}},
		{"sqlite3", types.File{Name: "", Content: `.tables`, IsMain: false}, &types.Output{Logs: []string{
			"0Category              EmployeeTerritory     Region              ",
			"0Customer              Order                 Shipper             ",
			"0CustomerCustomerDemo  OrderDetail           Supplier            ",
			"0CustomerDemographic   Product               Territory           ",
			"0Employee              ProductDetails_V    ",
		}}},
	}
	for _, tc := range testcases {
		t.Run(tc.lang+"__", func(t *testing.T) {
			input := types.MultiInput{Lang: tc.lang, Files: []types.File{tc.file}}
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
		file        types.File
		wantOutput  *types.Output
		wantLengths []int
	}{
		{"python", types.File{Name: "", Content: `print(100*"X")`, IsMain: false}, &types.Output{}, []int{101}},
		{"python", types.File{Name: "", Content: `print(1000*"X")`, IsMain: false}, &types.Output{}, []int{1001}},
		{"python", types.File{Name: "", Content: `print(10000*"X")`, IsMain: false}, &types.Output{Warning: types.WarnOutputLimitReached}, []int{8001}},
		{"python", types.File{Name: "", Content: `[print(1000*"X") for _ in range(1)]`, IsMain: false}, &types.Output{}, []int{1001}},
		{"python", types.File{Name: "", Content: `[print(1000*"X") for _ in range(10)]`, IsMain: false}, &types.Output{Warning: types.WarnOutputLimitReached}, []int{1001, 1001, 1001, 1001, 1001, 1001, 1001, 994}},
	}
	for _, tc := range testcases {
		t.Run(tc.lang+"__", func(t *testing.T) {
			input := types.MultiInput{Lang: tc.lang, Files: []types.File{tc.file}}
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
		file       types.File
		wantOutput *types.Output
		wantError  string
	}{
		// Bash
		{
			"bash", types.File{Name: "", Content: `echo hello`, IsMain: false},
			&types.Output{Logs: []string{"0hello"}}, "",
		},
		{
			"bash", types.File{Name: "", Content: `echo hello; sleep 3`, IsMain: false},
			&types.Output{Logs: []string{"0hello"}, Warning: types.WarnTimeout}, "",
		},
		{
			"bash", types.File{Name: "", Content: `sleep 3; echo hello`, IsMain: false},
			&types.Output{Logs: []string{}, Warning: types.WarnTimeout}, "",
		},
		{
			"bash", types.File{Name: "", Content: `echo hello; echo world`, IsMain: false},
			&types.Output{Logs: []string{"0hello", "0world"}}, "",
		},
		{
			"bash", types.File{Name: "", Content: `echo hello; echo world; sleep 3`, IsMain: false},
			&types.Output{Logs: []string{"0hello", "0world"}, Warning: types.WarnTimeout}, "",
		},
		{
			"bash", types.File{Name: "", Content: `echo hello; sleep 3; echo world`, IsMain: false},
			&types.Output{Logs: []string{"0hello"}, Warning: types.WarnTimeout}, "",
		},
		{
			"bash", types.File{Name: "", Content: `sleep 3; echo hello; echo world`, IsMain: false},
			&types.Output{Logs: []string{}, Warning: types.WarnTimeout}, "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.lang, func(t *testing.T) {
			input := types.MultiInput{Lang: tc.lang, Files: []types.File{tc.file}}
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
