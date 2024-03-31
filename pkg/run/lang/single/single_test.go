package single

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetaoss/zetarun/pkg/run/lang/types"
)

type TestCaseOK struct {
	lang       string
	source     string
	wantOutput *types.Output
}

func TestRun_error(t *testing.T) {
	testCases := []struct {
		lang      string
		source    string
		wantError string
	}{
		{"", "", "invalid language"},
		{"X", "", "invalid language"},
	}
	for _, tc := range testCases {
		t.Run(tc.lang, func(t *testing.T) {
			input := types.SingleInput{Lang: tc.lang, Source: tc.source}
			output, err := Run(input)
			require.EqualError(t, err, tc.wantError)
			require.Nil(t, output)
		})
	}
}

func TestRun_simple(t *testing.T) {
	testCases := []TestCaseOK{
		// Bash
		{"bash", `echo hello`, &types.Output{Logs: []string{"0hello"}}},
		{"bash", `echo hello; echo world`, &types.Output{Logs: []string{"0hello", "0world"}}},
		{"bash", `echo hello; echo world; echo`, &types.Output{Logs: []string{"0hello", "0world", "0"}}},
		{"bash", `echo hello; echo world; echo`, &types.Output{Logs: []string{"0hello", "0world", "0"}}},
		{"bash", `echo hello 1>&2; echo world 1>&2`, &types.Output{Logs: []string{"1hello", "1world"}}},
		// Lua
		{"lua", `print("Hello, World!")`, &types.Output{Logs: []string{"0Hello, World!"}}},
		// PowerShell
		{"powershell", `Write-Host "Hello, World!"`, &types.Output{Logs: []string{"0Hello, World!"}}},
		// Python
		{"python", `print("Hello, World!")`, &types.Output{Logs: []string{"0Hello, World!"}}},
		// R
		{"r", `print("Hello, World!")`, &types.Output{Logs: []string{`0[1] "Hello, World!"`}}},
		// Ruby
		{"ruby", `print("Hello, World!")`, &types.Output{Logs: []string{"0Hello, World!"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.lang, func(t *testing.T) {
			input := types.SingleInput{Lang: tc.lang, Source: tc.source}
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

func TestRun_db(t *testing.T) {
	testCases := []TestCaseOK{
		// MySQL
		{"mysql", `SELECT 'Hello, World!';`, &types.Output{Logs: []string{
			"0+---------------+",
			"0| Hello, World! |",
			"0+---------------+",
			"0| Hello, World! |",
			"0+---------------+"}}},
		// SQLite
		{"sqlite3", `SELECT 'Hello, World!';`, &types.Output{Logs: []string{
			"0+-----------------+",
			"0| 'Hello, World!' |",
			"0+-----------------+",
			"0| Hello, World!   |",
			"0+-----------------+"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.lang, func(t *testing.T) {
			input := types.SingleInput{Lang: tc.lang, Source: tc.source}
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

func TestRun_multiline(t *testing.T) {
	testCases := []TestCaseOK{
		// C
		{"c", `
#include <stdio.h>
int main() {
	printf("Hello, World!");
}`, &types.Output{Logs: []string{"0Hello, World!"}}},
		{"c", `
#include <stdio.h>
int main() {
	printf("Hello\nWorld!");
}`, &types.Output{Logs: []string{"0Hello", "0World!"}}},
		// C++
		{"cpp", `
#include <iostream>
int main() {
	std::cout<<"hello";
}`, &types.Output{Logs: []string{"0hello"}}},
		// C#
		{"csharp", `
using System;
class Program
{
	static void Main() {
		Console.Write("hello");
	}
}`, &types.Output{Logs: []string{"0hello"}}},
		// Java
		{"java", `
public class App {
	public static void main(String args[]) {
		System.out.println("hello");
	}
}`, &types.Output{Logs: []string{"0hello"}}},
		{"java", `
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
}`, &types.Output{
			Logs:   []string{},
			Images: []string{"iVBORw0KGgoAAAANSUhEUgAAASwAAADICAIAAADdvUsCAAACaUlEQVR4Xu3TQW4bMRAAQf3/08pBwGLD4VqyEbiloOpgkLMUfWHf7kDqtg6A3yVCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYjcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP6l+/2+LH7T1T+9msO7m293ThY/iHA5+fX2a1eHr+bw7ubbnZPFDyK8/X34av2Kq/NXc3h38+0uhRzmgacnz7Y3zPW85LE+hlfnz3P4JPPtHpPzy17Wc7Fdny2/ev2S5cLtfDkDn+TxfBfHp+Xk08V2e3jM59/z17n97hw+zHy751e+mAeWxXZ7eMzn3/PXuf3uHD7MfLtXr/wwD7zYw8zvfPLqku/O4cPMt3t+5dtItovtetpeeGy3l8xjx2K7hg8z3+558njcyxM/1k9PTtt7zpN5yXJs+bT9CQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADwv/sDs/dz4IQQs5EAAAAASUVORK5CYII="},
		}},
		// Kotlin
		{"kotlin", `
fun main() {
	println("Hello, World!")
}`, &types.Output{Logs: []string{"0Hello, World!"}}},
		// Go
		{"go", `
package main
import "fmt"
func main() {
	fmt.Println("Hello, World!")
}`, &types.Output{Logs: []string{"0Hello, World!"}}},
		// Perl
		{"perl", `
use strict;
use warnings;

print "Hello, World!\n";
`, &types.Output{Logs: []string{"0Hello, World!"}}},
		// PHP
		{"php", `
<?php
echo "Hello, World!";
`, &types.Output{Logs: []string{"0Hello, World!"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.lang, func(t *testing.T) {
			input := types.SingleInput{Lang: tc.lang, Source: tc.source}
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

func TestRun_timeout(t *testing.T) {
	testCases := []TestCaseOK{
		{"bash", `echo hello`, &types.Output{Logs: []string{"0hello"}, Timeout: false}},
		{"bash", `echo hello; echo world`, &types.Output{Logs: []string{"0hello", "0world"}, Timeout: false}},

		{"bash", `echo hello; sleep 3`, &types.Output{Logs: []string{"0hello"}, Timeout: true}},
		{"bash", `sleep 3; echo hello`, &types.Output{Logs: []string{}, Timeout: true}},
		{"bash", `echo hello; echo world; sleep 3`, &types.Output{Logs: []string{"0hello", "0world"}, Timeout: true}},
		{"bash", `echo hello; sleep 3; echo world`, &types.Output{Logs: []string{"0hello"}, Timeout: true}},
		{"bash", `sleep 3; echo hello; echo world`, &types.Output{Logs: []string{}, Timeout: true}},
	}
	for _, tc := range testCases {
		t.Run(tc.lang, func(t *testing.T) {
			input := types.SingleInput{Lang: tc.lang, Source: tc.source}
			output, err := Run(input, map[string]int{"timeoutSeconds": 1})
			require.NoError(t, err)
			// ignore fields
			output.Time = ""
			output.CPU = 0
			output.MEM = 0
			require.Equal(t, tc.wantOutput, output)
		})
	}
}
