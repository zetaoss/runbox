package box

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/client"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/zetaoss/runbox/pkg/testutil"
)

var box1 *Box

func init() {
	d := testutil.NewDocker()
	box1 = New(d)
}

func equalResult(t *testing.T, want, got *Result) {
	t.Helper()

	assert.Greater(t, got.CPU, want.CPU/10, "cpu")
	assert.Greater(t, got.MEM, want.MEM/10, "mem")
	assert.Less(t, got.CPU, want.CPU*100, "cpu")
	assert.Less(t, got.MEM, want.MEM*10000, "mem")
	want.CPU = got.CPU
	want.MEM = got.MEM

	assert.Greater(t, got.Time, want.Time/10, "time")
	assert.Less(t, got.Time, want.Time*100, "time")
	want.Time = got.Time

	assert.Equal(t, want, got)
}

func equalStructSlices(a, b interface{}) bool {
	return cmp.Equal(a, b, cmpopts.SortSlices(func(x, y interface{}) bool {
		return fmt.Sprintf("%v", x) < fmt.Sprintf("%v", y)
	}))
}

func TestDockerEnv(t *testing.T) {
	assert.NotEmpty(t, os.Getenv("DOCKER_HOST"))
	assert.Equal(t, "1", os.Getenv("DOCKER_TLS_VERIFY"))
	assert.Equal(t, os.Getenv("HOME")+"/.docker", os.Getenv("DOCKER_CERT_PATH"))

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	assert.NoError(t, err)
	assert.NotEmpty(t, cli)
}

func TestRun_error(t *testing.T) {
	testCases := []struct {
		opts      *Opts
		wantError string
	}{
		{
			&Opts{},
			"checkImage err: invalid reference format",
		},
		{
			&Opts{Image: "."},
			"checkImage err: invalid reference format",
		},
		{
			&Opts{Image: "a"},
			"checkImage err: no image: 'a'",
		},
		{
			&Opts{Image: "xxx"},
			"checkImage err: no image: 'xxx'",
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.opts.Image), func(t *testing.T) {
			got, err := box1.Run(tc.opts)
			assert.EqualError(t, err, tc.wantError)
			assert.Nil(t, got)
		})
	}
}

func TestRun_alpine(t *testing.T) {
	testCases := []struct {
		opts *Opts
		want *Result
	}{
		{
			&Opts{Image: "alpine", Command: "echo hello"},
			&Result{
				Logs: []Log{{Stream: 1, Log: "hello"}},
				CPU:  20240,
				MEM:  462,
				Time: 18,
			},
		},
		{
			&Opts{Image: "alpine", Command: "echo hello; echo; echo world; echo"},
			&Result{
				Logs: []Log{
					{Stream: 1, Log: "hello"},
					{Stream: 1, Log: ""},
					{Stream: 1, Log: "world"},
					{Stream: 1, Log: ""},
				},
				CPU:  20816,
				MEM:  462,
				Time: 18,
			},
		},
		{
			&Opts{Image: "alpine", Command: "cat /etc/os-release | head -2"},
			&Result{
				Logs: []Log{
					{Stream: 1, Log: "NAME=\"Alpine Linux\""},
					{Stream: 1, Log: "ID=alpine"},
				},
				CPU:  20174,
				MEM:  512,
				Time: 15,
			},
		},
		{
			&Opts{Image: "alpine", Command: "foo"},
			&Result{
				Logs: []Log{{Stream: 2, Log: "sh: foo: not found"}},
				Code: 127,
				CPU:  20787,
				MEM:  2466,
				Time: 20,
			},
		},
		{
			&Opts{Image: "alpine", Command: "echo hello; exit 42"},
			&Result{
				Logs: []Log{{Stream: 1, Log: "hello"}},
				Code: 42,
				CPU:  20510,
				MEM:  466,
				Time: 16,
			},
		},
		{
			&Opts{Image: "alpine", Command: "echo a"},
			&Result{
				Logs: []Log{{Stream: 1, Log: "a"}},
				CPU:  19352,
				MEM:  462,
				Time: 17,
			},
		},
		{
			&Opts{Image: "jmnote/runbox:python", Shell: "python", Command: "print('hello')"},
			&Result{
				Logs: []Log{{Stream: 1, Log: "hello"}},
				CPU:  31282,
				MEM:  524,
				Time: 33,
			},
		},
		{
			&Opts{
				Image:      "jmnote/runbox:java",
				Command:    `javac -d bin -cp "lib/*" src/*; java -cp "bin:lib/*" App`,
				WorkingDir: "/demo",
				Files:      []File{{"/demo/src/App.java", `public class App{public static void main(String args[]){System.out.println("hello");}}`}},
			},
			&Result{
				Logs: []Log{{Stream: 1, Log: "hello"}},
				CPU:  837273,
				MEM:  772,
				Time: 992,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.opts.Image, tc.opts.Command), func(t *testing.T) {
			got, err := box1.Run(tc.opts)
			assert.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_language(t *testing.T) {
	testCases := []struct {
		opts *Opts
		want *Result
	}{
		{
			&Opts{Image: "jmnote/runbox:python", Shell: "python", Command: "print('hello')"},
			&Result{
				Logs: []Log{{Stream: 1, Log: "hello"}},
				CPU:  31282,
				MEM:  536,
				Time: 33,
			},
		},
		{
			&Opts{
				Image:      "jmnote/runbox:java",
				Command:    `javac -d bin -cp "lib/*" src/*; java -cp "bin:lib/*" App`,
				WorkingDir: "/demo",
				Files:      []File{{"/demo/src/App.java", `public class App{public static void main(String args[]){System.out.println("hello");}}`}},
			},
			&Result{
				Logs: []Log{{Stream: 1, Log: "hello"}},
				CPU:  837273,
				MEM:  772,
				Time: 992,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.opts.Image, tc.opts.Command), func(t *testing.T) {
			got, err := box1.Run(tc.opts)
			assert.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_stderr(t *testing.T) {
	testCases := []struct {
		opts *Opts
		want *Result
	}{
		{
			&Opts{Image: "alpine", Command: "echo stdout; echo stderr >&2; echo stdout; echo stderr >&2; echo stdout; echo stderr >&2; echo stdout; echo stderr >&2; echo stdout; echo stderr >&2"},
			&Result{
				Logs: []Log{
					{Stream: 1, Log: "stdout"},
					{Stream: 2, Log: "stderr"},
					{Stream: 1, Log: "stdout"},
					{Stream: 2, Log: "stderr"},
					{Stream: 1, Log: "stdout"},
					{Stream: 2, Log: "stderr"},
					{Stream: 1, Log: "stdout"},
					{Stream: 2, Log: "stderr"},
					{Stream: 1, Log: "stdout"},
					{Stream: 2, Log: "stderr"},
				},
				CPU:      21435,
				MEM:      452,
				Time:     18,
				Timedout: false,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.opts.Image, tc.opts.Command), func(t *testing.T) {
			got, err := box1.Run(tc.opts)
			assert.NoError(t, err)
			assert.True(t, equalStructSlices(tc.want.Logs, got.Logs))
			tc.want.Logs = got.Logs
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_files(t *testing.T) {
	testCases := []struct {
		opts *Opts
		want *Result
	}{
		{
			&Opts{Image: "alpine", Command: "cat /tmp/hello.txt", Files: []File{
				{"/tmp/hello.txt", "world"},
			}},
			&Result{
				Logs:     []Log{{Stream: 1, Log: "world"}},
				CPU:      20637,
				MEM:      508,
				Time:     19,
				Timedout: false,
			},
		},
		{
			&Opts{Image: "jmnote/runbox:python", Shell: "python", Command: "print(open('/tmp/hello.txt').read())", Files: []File{
				{"/tmp/hello.txt", "world"},
			}},
			&Result{
				Logs: []Log{{Stream: 1, Log: "world"}},

				CPU:      32381,
				MEM:      524,
				Time:     30,
				Timedout: false,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.opts.Image, tc.opts.Command), func(t *testing.T) {
			got, err := box1.Run(tc.opts)
			assert.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_forkbomb(t *testing.T) {
	testCases := []struct {
		opts *Opts
		want *Result
	}{
		{
			&Opts{
				Image:   "bash",
				Shell:   "bash",
				Command: ":(){ :|:& };:",
			},
			&Result{
				Logs:     []Log{{Stream: 2, Log: "Resource temporarily unavailable"}},
				CPU:      82547,
				MEM:      28544,
				Time:     2016,
				Timedout: false,
			},
		},
		{
			&Opts{
				Image:   "jmnote/runbox:python",
				Shell:   "python",
				Command: "import os; [os.fork() for _ in range(1000)]",
			},
			&Result{
				Logs: []Log{
					{Stream: 2, Log: "Traceback (most recent call last):"},
					{Stream: 2, Log: "  File \"<string>\", line 1, in <module>"},
				},
				Code:     1,
				CPU:      773947,
				MEM:      4336,
				Time:     908,
				Timedout: false,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.opts.Image, tc.opts.Command), func(t *testing.T) {
			got, err := box1.Run(tc.opts)
			assert.NoError(t, err)

			var logsString string
			for _, l := range got.Logs {
				logsString += l.Log
			}
			for _, l := range tc.want.Logs {
				assert.Contains(t, logsString, l.Log)
			}
			tc.want.Logs = got.Logs
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_timeout(t *testing.T) {
	testCases := []struct {
		opts *Opts
		want *Result
	}{
		{
			&Opts{
				Image:   "alpine",
				Command: "sleep 10",
				Timeout: 500,
			},
			&Result{
				CPU:      20195,
				MEM:      356,
				Time:     500,
				Timedout: true,
			},
		},
		{
			&Opts{
				Image:   "alpine",
				Command: "echo hello; sleep 10",
				Timeout: 500,
			},
			&Result{
				Logs:     []Log{{Stream: 1, Log: "hello"}},
				CPU:      20154,
				MEM:      368,
				Time:     500,
				Timedout: true,
			},
		},
	}

	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.opts.Image, tc.opts.Command), func(t *testing.T) {
			got, err := box1.Run(tc.opts)
			assert.NoError(t, err)

			assert.Greater(t, got.Time, tc.want.Time*10/12, "time")
			assert.Less(t, got.Time, tc.want.Time*12/10, "time")
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_images(t *testing.T) {
	testCases := []struct {
		opts *Opts
		want *Result
	}{
		{
			&Opts{
				CollectImages: true,
				Image:         "jmnote/runbox:java",
				Command:       `javac -d bin -cp "lib/*" src/*; java -cp "bin:lib/*" App`,
				WorkingDir:    "/demo",
				Files: []File{{
					"/demo/src/App.java", `
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
			&Result{
				CPU:  988464,
				MEM:  968,
				Time: 1016,
				Images: []string{
					"iVBORw0KGgoAAAANSUhEUgAAASwAAADICAIAAADdvUsCAAACaUlEQVR4Xu3TQW4bMRAAQf3/08pBwGLD4VqyEbiloOpgkLMUfWHf7kDqtg6A3yVCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYjcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP6l+/2+LH7T1T+9msO7m293ThY/iHA5+fX2a1eHr+bw7ubbnZPFDyK8/X34av2Kq/NXc3h38+0uhRzmgacnz7Y3zPW85LE+hlfnz3P4JPPtHpPzy17Wc7Fdny2/ev2S5cLtfDkDn+TxfBfHp+Xk08V2e3jM59/z17n97hw+zHy751e+mAeWxXZ7eMzn3/PXuf3uHD7MfLtXr/wwD7zYw8zvfPLqku/O4cPMt3t+5dtItovtetpeeGy3l8xjx2K7hg8z3+558njcyxM/1k9PTtt7zpN5yXJs+bT9CQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADwv/sDs/dz4IQQs5EAAAAASUVORK5CYII=",
				},
			},
		},
		{
			&Opts{
				CollectImages: true,
				Image:         "jmnote/runbox:java",
				Command:       `javac -d bin -cp "lib/*" src/*; java -cp "bin:lib/*" App`,
				WorkingDir:    "/demo",
				Files: []File{{
					"/demo/src/App.java", `
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
							ImageIO.write(bufferedImage, "png", new File("image1.png"));
							ImageIO.write(bufferedImage, "png", new File("image2.png"));
							ImageIO.write(bufferedImage, "png", new File("image3.png"));
						}
					}`}},
			},
			&Result{
				Code:     0,
				CPU:      988464,
				MEM:      968,
				Time:     1016,
				Timedout: false,
				Images: []string{
					"iVBORw0KGgoAAAANSUhEUgAAASwAAADICAIAAADdvUsCAAACaUlEQVR4Xu3TQW4bMRAAQf3/08pBwGLD4VqyEbiloOpgkLMUfWHf7kDqtg6A3yVCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYjcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP6l+/2+LH7T1T+9msO7m293ThY/iHA5+fX2a1eHr+bw7ubbnZPFDyK8/X34av2Kq/NXc3h38+0uhRzmgacnz7Y3zPW85LE+hlfnz3P4JPPtHpPzy17Wc7Fdny2/ev2S5cLtfDkDn+TxfBfHp+Xk08V2e3jM59/z17n97hw+zHy751e+mAeWxXZ7eMzn3/PXuf3uHD7MfLtXr/wwD7zYw8zvfPLqku/O4cPMt3t+5dtItovtetpeeGy3l8xjx2K7hg8z3+558njcyxM/1k9PTtt7zpN5yXJs+bT9CQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADwv/sDs/dz4IQQs5EAAAAASUVORK5CYII=",
					"iVBORw0KGgoAAAANSUhEUgAAASwAAADICAIAAADdvUsCAAACaUlEQVR4Xu3TQW4bMRAAQf3/08pBwGLD4VqyEbiloOpgkLMUfWHf7kDqtg6A3yVCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYiKEmAghJkKIiRBiIoSYCCEmQoiJEGIihJgIISZCiIkQYjcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAP6l+/2+LH7T1T+9msO7m293ThY/iHA5+fX2a1eHr+bw7ubbnZPFDyK8/X34av2Kq/NXc3h38+0uhRzmgacnz7Y3zPW85LE+hlfnz3P4JPPtHpPzy17Wc7Fdny2/ev2S5cLtfDkDn+TxfBfHp+Xk08V2e3jM59/z17n97hw+zHy751e+mAeWxXZ7eMzn3/PXuf3uHD7MfLtXr/wwD7zYw8zvfPLqku/O4cPMt3t+5dtItovtetpeeGy3l8xjx2K7hg8z3+558njcyxM/1k9PTtt7zpN5yXJs+bT9CQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADwv/sDs/dz4IQQs5EAAAAASUVORK5CYII=",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.opts.Image, tc.opts.Command), func(t *testing.T) {
			got, err := box1.Run(tc.opts)
			assert.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}

func TestRun_tex(t *testing.T) {
	testCases := []struct {
		opts *Opts
		want *Result
	}{
		{
			&Opts{
				CollectImages: true,
				Image:         "jmnote/runbox:tex",
				Command:       `touch oblivoir.sty && pdflatex -halt-on-error runbox.tex && convert runbox.pdf -strip p%d.png`,
				WorkingDir:    "/home/user01",
				User:          "root",
				Files: []File{{
					"/home/user01/runbox.tex",
					"\\documentclass{article}\n\\usepackage[a6paper,landscape]{geometry}\n\\begin{document}\nHello world!\n\\end{document}",
				}},
			},
			&Result{Logs: []Log{
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
			}, Code: 0, CPU: 735518, MEM: 724, Time: 686, Timedout: false, Images: []string{"iVBORw0KGgoAAAANSUhEUgAAAaQAAAEqCAQAAADw/+wWAAAE4klEQVR42u3X3Y0bRxpA0a8WTqBTYApMgSkwhUlBm4JCcAL74ElhQjBDMFNgCL0PkkY/gFcr6wLkCOc8NMkqdqMK6Ism1z7Az/rXvRcAvwIhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQUBIEBASBIQEASFBQEgQEBIEhAQBIUFASBAQEgSEBAEhQeChQ1rH9ft6Wts6rj/W+evxT8cfv+Lr+6f1bmZt669775O376FD2i9zmOt+2y+z7c/fjH88/oMrfnKdbWa/zfXe++Tt++3eC/gR6zTbXPavbvx1nG2un8fWcba5zTbX2WbbXz7Nf3h9/dZhth/PEP7OQz+RZmbmuE7rNNvMOs9hf553X06u0xz3lzmv0+vQdZ72y5z369y+nN8v8/tsc5yZWcc575e5fTzjcu8t8vY9fkiX/WV/mdvMnOe2Tq+3/wdP8zwzl3n9B7XfZtY2t3Wew/7y1fx1v3z8gfg0l5nXK73ce4u8fY8f0me3ue4v8/6bscPH18+e5928n+NsfzP/7SdPJH7aQ4e0jrPNaW3rOId1nn/PaZ3mNLOOc1iHdZzD/GdO6zTHL/Pan2fbb/MhkPef5tdxDuv44cx5P6d1mOMc1ra2+fPe++TtW/u9VwC/gId+IsFbISQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkB7eerr3Cvg+IT24dZrzvdfA9wnpwe0v914B/w8hQUBIEBASBIT04NZpDuu8tnuvg/9t7fdeAfwCPJEgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAkKCgJAgICQICAkCQoKAkCAgJAgICQJCgoCQICAkCAgJAv8F6NR4np+BlBIAAAAASUVORK5CYII="}},
		},
		{
			&Opts{
				CollectImages:      true,
				CollectImagesCount: 8,
				Image:              "jmnote/runbox:tex",
				Command:            `touch oblivoir.sty && pdflatex -halt-on-error runbox.tex && convert runbox.pdf -strip p%d.png`,
				WorkingDir:         "/home/user01",
				User:               "root",
				Files: []File{{
					"/home/user01/runbox.tex",
					"\\documentclass{article}\\usepackage[a6paper,landscape]{geometry}" +
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
						"\\end{document}",
				}},
			},
			&Result{
				Logs: []Log{
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
				CPU:      1417094,
				MEM:      848,
				Time:     968,
				Timedout: false,
				Images: []string{
					"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",
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.opts.Image, tc.opts.Command), func(t *testing.T) {
			got, err := box1.Run(tc.opts)
			assert.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}
