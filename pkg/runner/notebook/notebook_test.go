package notebook

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zetaoss/runbox/pkg/runner/notebook/nbformat"
)

func TestRun(t *testing.T) {
	testCases := []struct {
		input      Input
		wantOutput *Output
		wantError  string
	}{
		{
			Input{},
			nil, "ErrInvalidLanguage",
		},
		{
			Input{Lang: "xxx"},
			nil, "ErrInvalidLanguage",
		},
		{
			Input{Lang: "python"},
			&Output{
				nbformat.Metadata{
					LanguageInfo: nbformat.LanguageInfo{Name: "python", CodemirrorMode: map[string]interface{}{"name": "ipython", "version": float64(3)}, FileExtension: ".py", Mimetype: "text/x-python", PygmentsLexer: "ipython3"},
				},
				nil,
			}, "",
		},
		{
			Input{Lang: "python", CellTexts: [][]string{
				{
					"print(\"hello1\")\n",
					"print(\"hello2\")\n",
				},
				{
					"print(\"world3\")\n",
					"print(\"world4\")\n",
				},
			}},
			&Output{
				Metadata: nbformat.Metadata{Kernelspec: nbformat.Kernelspec{}, LanguageInfo: nbformat.LanguageInfo{Name: "python", CodemirrorMode: map[string]interface{}{"name": "ipython", "version": float64(3)}, FileExtension: ".py", Mimetype: "text/x-python", PygmentsLexer: "ipython3"}},
				CellOutputs: [][]nbformat.Output{
					{{OutputType: "stream", Name: "stdout", Text: []string{"hello1\n", "hello2\n"}}},
					{{OutputType: "stream", Name: "stdout", Text: []string{"world3\n", "world4\n"}}},
				},
			}, "",
		},
		{
			Input{Lang: "r"},
			&Output{
				Metadata: nbformat.Metadata{
					Kernelspec:   nbformat.Kernelspec{Name: "ir"},
					LanguageInfo: nbformat.LanguageInfo{Name: "R", CodemirrorMode: "r", FileExtension: ".r", Mimetype: "text/x-r-source", PygmentsLexer: "r"},
				},
				CellOutputs: nil,
			}, "",
		},
		{
			Input{
				Lang: "r",
				CellTexts: [][]string{
					{
						`print("hello")`,
						`print("world")`,
					},
					{
						`head(iris)`,
					},
				},
			},
			&Output{
				Metadata: nbformat.Metadata{
					Kernelspec:   nbformat.Kernelspec{Name: "ir", DisplayName: "", Language: ""},
					LanguageInfo: nbformat.LanguageInfo{Name: "R", CodemirrorMode: "r", FileExtension: ".r", Mimetype: "text/x-r-source", PygmentsLexer: "r"},
				},
				CellOutputs: [][]nbformat.Output{
					{{
						OutputType:     "error",
						ExceptionName:  "ERROR",
						Exceptionvalue: "Error in parse(text = x, srcfile = src): <text>:1:15: unexpected symbol\n1: print(\"hello\")print\n                  ^\n",
						Traceback:      []string{"Error in parse(text = x, srcfile = src): <text>:1:15: unexpected symbol\n1: print(\"hello\")print\n                  ^\nTraceback:\n"},
					}},
					{{
						OutputType: "display_data",
						Metadata:   map[string]interface{}{},
						Data: map[string]interface{}{
							"text/html": []interface{}{"<table class=\"dataframe\">\n", "<caption>A data.frame: 6 × 5</caption>\n", "<thead>\n", "\t<tr><th></th><th scope=col>Sepal.Length</th><th scope=col>Sepal.Width</th><th scope=col>Petal.Length</th><th scope=col>Petal.Width</th><th scope=col>Species</th></tr>\n", "\t<tr><th></th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;fct&gt;</th></tr>\n", "</thead>\n", "<tbody>\n", "\t<tr><th scope=row>1</th><td>5.1</td><td>3.5</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>2</th><td>4.9</td><td>3.0</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>3</th><td>4.7</td><td>3.2</td><td>1.3</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>4</th><td>4.6</td><td>3.1</td><td>1.5</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>5</th><td>5.0</td><td>3.6</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>6</th><td>5.4</td><td>3.9</td><td>1.7</td><td>0.4</td><td>setosa</td></tr>\n", "</tbody>\n", "</table>\n"}, "text/latex": []interface{}{"A data.frame: 6 × 5\n", "\\begin{tabular}{r|lllll}\n", "  & Sepal.Length & Sepal.Width & Petal.Length & Petal.Width & Species\\\\\n", "  & <dbl> & <dbl> & <dbl> & <dbl> & <fct>\\\\\n", "\\hline\n", "\t1 & 5.1 & 3.5 & 1.4 & 0.2 & setosa\\\\\n", "\t2 & 4.9 & 3.0 & 1.4 & 0.2 & setosa\\\\\n", "\t3 & 4.7 & 3.2 & 1.3 & 0.2 & setosa\\\\\n", "\t4 & 4.6 & 3.1 & 1.5 & 0.2 & setosa\\\\\n", "\t5 & 5.0 & 3.6 & 1.4 & 0.2 & setosa\\\\\n", "\t6 & 5.4 & 3.9 & 1.7 & 0.4 & setosa\\\\\n", "\\end{tabular}\n"},
							"text/markdown": []interface{}{"\n", "A data.frame: 6 × 5\n", "\n", "| <!--/--> | Sepal.Length &lt;dbl&gt; | Sepal.Width &lt;dbl&gt; | Petal.Length &lt;dbl&gt; | Petal.Width &lt;dbl&gt; | Species &lt;fct&gt; |\n", "|---|---|---|---|---|---|\n", "| 1 | 5.1 | 3.5 | 1.4 | 0.2 | setosa |\n", "| 2 | 4.9 | 3.0 | 1.4 | 0.2 | setosa |\n", "| 3 | 4.7 | 3.2 | 1.3 | 0.2 | setosa |\n", "| 4 | 4.6 | 3.1 | 1.5 | 0.2 | setosa |\n", "| 5 | 5.0 | 3.6 | 1.4 | 0.2 | setosa |\n", "| 6 | 5.4 | 3.9 | 1.7 | 0.4 | setosa |\n", "\n"}, "text/plain": []interface{}{"  Sepal.Length Sepal.Width Petal.Length Petal.Width Species\n", "1 5.1          3.5         1.4          0.2         setosa \n", "2 4.9          3.0         1.4          0.2         setosa \n", "3 4.7          3.2         1.3          0.2         setosa \n", "4 4.6          3.1         1.5          0.2         setosa \n", "5 5.0          3.6         1.4          0.2         setosa \n", "6 5.4          3.9         1.7          0.4         setosa "}},
					}},
				},
			},
			"",
		},
		{
			Input{
				Lang: "r",
				CellTexts: [][]string{
					{
						`library(ggplot2)`,
						`table(mpg$class)`,
					},
					{
						`head(iris)`,
					},
				},
			},
			&Output{
				Metadata: nbformat.Metadata{
					Kernelspec:   nbformat.Kernelspec{Name: "ir", DisplayName: "", Language: ""},
					LanguageInfo: nbformat.LanguageInfo{Name: "R", CodemirrorMode: "r", FileExtension: ".r", Mimetype: "text/x-r-source", PygmentsLexer: "r"},
				},
				CellOutputs: [][]nbformat.Output{
					{{
						OutputType:     "error",
						ExceptionName:  "ERROR",
						Exceptionvalue: "Error in parse(text = x, srcfile = src): <text>:1:17: unexpected symbol\n1: library(ggplot2)table\n                    ^\n",
						Traceback:      []string{"Error in parse(text = x, srcfile = src): <text>:1:17: unexpected symbol\n1: library(ggplot2)table\n                    ^\nTraceback:\n"},
					}},
					{{
						OutputType: "display_data",
						Metadata:   map[string]interface{}{},
						Data: map[string]interface{}{
							"text/html":     []interface{}{"<table class=\"dataframe\">\n", "<caption>A data.frame: 6 × 5</caption>\n", "<thead>\n", "\t<tr><th></th><th scope=col>Sepal.Length</th><th scope=col>Sepal.Width</th><th scope=col>Petal.Length</th><th scope=col>Petal.Width</th><th scope=col>Species</th></tr>\n", "\t<tr><th></th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;fct&gt;</th></tr>\n", "</thead>\n", "<tbody>\n", "\t<tr><th scope=row>1</th><td>5.1</td><td>3.5</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>2</th><td>4.9</td><td>3.0</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>3</th><td>4.7</td><td>3.2</td><td>1.3</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>4</th><td>4.6</td><td>3.1</td><td>1.5</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>5</th><td>5.0</td><td>3.6</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n", "\t<tr><th scope=row>6</th><td>5.4</td><td>3.9</td><td>1.7</td><td>0.4</td><td>setosa</td></tr>\n", "</tbody>\n", "</table>\n"},
							"text/latex":    []interface{}{"A data.frame: 6 × 5\n", "\\begin{tabular}{r|lllll}\n", "  & Sepal.Length & Sepal.Width & Petal.Length & Petal.Width & Species\\\\\n", "  & <dbl> & <dbl> & <dbl> & <dbl> & <fct>\\\\\n", "\\hline\n", "\t1 & 5.1 & 3.5 & 1.4 & 0.2 & setosa\\\\\n", "\t2 & 4.9 & 3.0 & 1.4 & 0.2 & setosa\\\\\n", "\t3 & 4.7 & 3.2 & 1.3 & 0.2 & setosa\\\\\n", "\t4 & 4.6 & 3.1 & 1.5 & 0.2 & setosa\\\\\n", "\t5 & 5.0 & 3.6 & 1.4 & 0.2 & setosa\\\\\n", "\t6 & 5.4 & 3.9 & 1.7 & 0.4 & setosa\\\\\n", "\\end{tabular}\n"},
							"text/markdown": []interface{}{"\n", "A data.frame: 6 × 5\n", "\n", "| <!--/--> | Sepal.Length &lt;dbl&gt; | Sepal.Width &lt;dbl&gt; | Petal.Length &lt;dbl&gt; | Petal.Width &lt;dbl&gt; | Species &lt;fct&gt; |\n", "|---|---|---|---|---|---|\n", "| 1 | 5.1 | 3.5 | 1.4 | 0.2 | setosa |\n", "| 2 | 4.9 | 3.0 | 1.4 | 0.2 | setosa |\n", "| 3 | 4.7 | 3.2 | 1.3 | 0.2 | setosa |\n", "| 4 | 4.6 | 3.1 | 1.5 | 0.2 | setosa |\n", "| 5 | 5.0 | 3.6 | 1.4 | 0.2 | setosa |\n", "| 6 | 5.4 | 3.9 | 1.7 | 0.4 | setosa |\n", "\n"}, "text/plain": []interface{}{"  Sepal.Length Sepal.Width Petal.Length Petal.Width Species\n", "1 5.1          3.5         1.4          0.2         setosa \n", "2 4.9          3.0         1.4          0.2         setosa \n", "3 4.7          3.2         1.3          0.2         setosa \n", "4 4.6          3.1         1.5          0.2         setosa \n", "5 5.0          3.6         1.4          0.2         setosa \n", "6 5.4          3.9         1.7          0.4         setosa "},
						},
					}},
				},
			},
			"",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.input.Lang, func(t *testing.T) {
			output, err := Run(tc.input)
			if tc.wantError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.wantError)
			}
			require.Equal(t, tc.wantOutput, output)
		})
	}
}
