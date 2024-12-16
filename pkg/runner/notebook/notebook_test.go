package notebook

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"github.com/zetaoss/runbox/pkg/testutil"
)

var notebook1 *Notebook

func init() {
	d := testutil.NewDocker()
	notebook1 = New(box.New(d))
}

func equalResult(t *testing.T, want, got *Result) {
	t.Helper()

	assert.Greater(t, got.CPU, want.CPU/32, "want.CPU", want.CPU)
	assert.Greater(t, got.MEM, want.MEM/32, "want.MEM", want.MEM)
	assert.Less(t, got.CPU, want.CPU*32, "want.CPU", want.CPU)
	assert.Less(t, got.MEM, want.MEM*32, "want.MEM", want.MEM)
	want.CPU = got.CPU
	want.MEM = got.MEM

	assert.Greater(t, got.Time, want.Time/8, "want.Time", want.Time)
	assert.Less(t, got.Time, want.Time*8, "want.Time", want.Time)
	want.Time = got.Time

	assert.Equal(t, want, got)
}

func TestRun_error(t *testing.T) {
	testCases := []struct {
		input     Input
		wantError string
	}{
		{
			Input{},
			"toFileBody err: invalid language",
		},
		{
			Input{Lang: "xxx"},
			"toFileBody err: invalid language",
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := notebook1.Run(tc.input)
			require.EqualError(t, err, tc.wantError)
			require.Nil(t, got)
		})
	}
}

func TestRun_ok(t *testing.T) {
	testCases := []struct {
		input Input
		want  *Result
	}{
		{
			Input{Lang: "python"},
			&Result{
				OutputsList: []Outputs{},
				CPU:         1195288,
				MEM:         3244,
				Time:        1712,
				Stderr:      "[NbConvertApp] Converting notebook /tmp/runbox.ipynb to notebook\n",
			},
		},
		{
			Input{
				Lang: "python",
				Sources: []string{
					`print("hello1")`,
					`print("world2")`,
				},
			},
			&Result{
				OutputsList: []Outputs{
					{
						Output{OutputType: "stream", Name: "stdout", Text: []string{"hello1\n"}},
					},
					{
						Output{OutputType: "stream", Name: "stdout", Text: []string{"world2\n"}},
					},
				},
				CPU:      1248567,
				MEM:      3256,
				Time:     1759,
				Timedout: false,
				Stderr:   "[NbConvertApp] Converting notebook /tmp/runbox.ipynb to notebook\n",
			},
		},
		{
			Input{Lang: "r"},
			&Result{
				OutputsList: []Outputs{},
				CPU:         1056951,
				MEM:         86544,
				Time:        1903,
				Stderr:      "[NbConvertApp] Converting notebook /tmp/runbox.ipynb to notebook\n\nExecution halted\n",
			},
		},
		{
			Input{
				Lang: "r",
				Sources: []string{
					`print("hello")`,
					`print("world")`,
					`head(iris)`,
				},
			},
			&Result{
				OutputsList: []Outputs{
					{Output{
						OutputType: "stream",
						Name:       "stdout",
						Text:       []string{"[1] \"hello\"\n"},
					}},
					{Output{
						OutputType: "stream",
						Name:       "stdout",
						Text:       []string{"[1] \"world\"\n"},
					}},
					{Output{
						OutputType: "display_data",
						Data: map[string]any{
							"text/html":     []any{"<table class=\"dataframe\">\n", "<caption>A data.frame: 6 × 5</caption>\n", "<thead>\n", "\t<tr><th></th><th scope=col>Sepal.Length</th><th scope=col>Sepal.Width</th><th scope=col>Petal.Length</th><th scope=col>Petal.Width</th><th scope=col>Species</th></tr>\n", "\t<tr><th></th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;fct&gt;</th></tr>\n", "</thead>\n", "<tbody>\n", "\t<tr><th scope=row>1</th><td>5.1</td><td>3.5</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>2</th><td>4.9</td><td>3.0</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>3</th><td>4.7</td><td>3.2</td><td>1.3</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>4</th><td>4.6</td><td>3.1</td><td>1.5</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>5</th><td>5.0</td><td>3.6</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>6</th><td>5.4</td><td>3.9</td><td>1.7</td><td>0.4</td><td>setosa</td></tr>\n", "</tbody>\n", "</table>\n"},
							"text/latex":    []any{"A data.frame: 6 × 5\n", "\\begin{tabular}{r|lllll}\n", "  & Sepal.Length & Sepal.Width & Petal.Length & Petal.Width & Species\\\\\n", "  & <dbl> & <dbl> & <dbl> & <dbl> & <fct>\\\\\n", "\\hline\n", "\t1 & 5.1 & 3.5 & 1.4 & 0.2 & setosa\\\\\n", "\t2 & 4.9 & 3.0 & 1.4 & 0.2 & setosa\\\\\n", "\t3 & 4.7 & 3.2 & 1.3 & 0.2 & setosa\\\\\n", "\t4 & 4.6 & 3.1 & 1.5 & 0.2 & setosa\\\\\n", "\t5 & 5.0 & 3.6 & 1.4 & 0.2 & setosa\\\\\n", "\t6 & 5.4 & 3.9 & 1.7 & 0.4 & setosa\\\\\n", "\\end{tabular}\n"},
							"text/markdown": []any{"\n", "A data.frame: 6 × 5\n", "\n", "| <!--/--> | Sepal.Length &lt;dbl&gt; | Sepal.Width &lt;dbl&gt; | Petal.Length &lt;dbl&gt; | Petal.Width &lt;dbl&gt; | Species &lt;fct&gt; |\n", "|---|---|---|---|---|---|\n", "| 1 | 5.1 | 3.5 | 1.4 | 0.2 | setosa |\n", "| 2 | 4.9 | 3.0 | 1.4 | 0.2 | setosa |\n", "| 3 | 4.7 | 3.2 | 1.3 | 0.2 | setosa |\n", "| 4 | 4.6 | 3.1 | 1.5 | 0.2 | setosa |\n", "| 5 | 5.0 | 3.6 | 1.4 | 0.2 | setosa |\n", "| 6 | 5.4 | 3.9 | 1.7 | 0.4 | setosa |\n", "\n"},
							"text/plain":    []any{"  Sepal.Length Sepal.Width Petal.Length Petal.Width Species\n", "1 5.1          3.5         1.4          0.2         setosa \n", "2 4.9          3.0         1.4          0.2         setosa \n", "3 4.7          3.2         1.3          0.2         setosa \n", "4 4.6          3.1         1.5          0.2         setosa \n", "5 5.0          3.6         1.4          0.2         setosa \n", "6 5.4          3.9         1.7          0.4         setosa "},
						},
						Metadata: &map[string]any{},
					}},
				},
				CPU:    1129615,
				MEM:    12620,
				Time:   1702,
				Stderr: "[NbConvertApp] Converting notebook /tmp/runbox.ipynb to notebook\n",
			},
		},
		{
			Input{
				Lang: "r",
				Sources: []string{
					`library(ggplot2)`,
					`table(mpg$class)`,
					`head(iris)`,
				},
			},
			&Result{
				OutputsList: []Outputs{
					{},
					{Output{
						OutputType: "display_data",
						Data:       map[string]any{"text/plain": []any{"\n", "   2seater    compact    midsize    minivan     pickup subcompact        suv \n", "         5         47         41         11         33         35         62 "}},
						Metadata:   &map[string]any{},
					}},
					{Output{
						OutputType: "display_data",
						Data:       map[string]any{"text/html": []any{"<table class=\"dataframe\">\n", "<caption>A data.frame: 6 × 5</caption>\n", "<thead>\n", "\t<tr><th></th><th scope=col>Sepal.Length</th><th scope=col>Sepal.Width</th><th scope=col>Petal.Length</th><th scope=col>Petal.Width</th><th scope=col>Species</th></tr>\n", "\t<tr><th></th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;fct&gt;</th></tr>\n", "</thead>\n", "<tbody>\n", "\t<tr><th scope=row>1</th><td>5.1</td><td>3.5</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>2</th><td>4.9</td><td>3.0</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>3</th><td>4.7</td><td>3.2</td><td>1.3</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>4</th><td>4.6</td><td>3.1</td><td>1.5</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>5</th><td>5.0</td><td>3.6</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>6</th><td>5.4</td><td>3.9</td><td>1.7</td><td>0.4</td><td>setosa</td></tr>\n", "</tbody>\n", "</table>\n"}, "text/latex": []any{"A data.frame: 6 × 5\n", "\\begin{tabular}{r|lllll}\n", "  & Sepal.Length & Sepal.Width & Petal.Length & Petal.Width & Species\\\\\n", "  & <dbl> & <dbl> & <dbl> & <dbl> & <fct>\\\\\n", "\\hline\n", "\t1 & 5.1 & 3.5 & 1.4 & 0.2 & setosa\\\\\n", "\t2 & 4.9 & 3.0 & 1.4 & 0.2 & setosa\\\\\n", "\t3 & 4.7 & 3.2 & 1.3 & 0.2 & setosa\\\\\n", "\t4 & 4.6 & 3.1 & 1.5 & 0.2 & setosa\\\\\n", "\t5 & 5.0 & 3.6 & 1.4 & 0.2 & setosa\\\\\n", "\t6 & 5.4 & 3.9 & 1.7 & 0.4 & setosa\\\\\n", "\\end{tabular}\n"}, "text/markdown": []any{"\n", "A data.frame: 6 × 5\n", "\n", "| <!--/--> | Sepal.Length &lt;dbl&gt; | Sepal.Width &lt;dbl&gt; | Petal.Length &lt;dbl&gt; | Petal.Width &lt;dbl&gt; | Species &lt;fct&gt; |\n", "|---|---|---|---|---|---|\n", "| 1 | 5.1 | 3.5 | 1.4 | 0.2 | setosa |\n", "| 2 | 4.9 | 3.0 | 1.4 | 0.2 | setosa |\n", "| 3 | 4.7 | 3.2 | 1.3 | 0.2 | setosa |\n", "| 4 | 4.6 | 3.1 | 1.5 | 0.2 | setosa |\n", "| 5 | 5.0 | 3.6 | 1.4 | 0.2 | setosa |\n", "| 6 | 5.4 | 3.9 | 1.7 | 0.4 | setosa |\n", "\n"}, "text/plain": []any{"  Sepal.Length Sepal.Width Petal.Length Petal.Width Species\n", "1 5.1          3.5         1.4          0.2         setosa \n", "2 4.9          3.0         1.4          0.2         setosa \n", "3 4.7          3.2         1.3          0.2         setosa \n", "4 4.6          3.1         1.5          0.2         setosa \n", "5 5.0          3.6         1.4          0.2         setosa \n", "6 5.4          3.9         1.7          0.4         setosa "}},
						Metadata:   &map[string]any{},
					}},
				},
				CPU:    1302084,
				MEM:    3644,
				Time:   1841,
				Stderr: "[NbConvertApp] Converting notebook /tmp/runbox.ipynb to notebook\n",
			},
		},
	}
	for i, tc := range testCases {
		t.Run(testutil.Name(i, tc.input), func(t *testing.T) {
			got, err := notebook1.Run(tc.input)
			require.NoError(t, err)
			equalResult(t, tc.want, got)
		})
	}
}
