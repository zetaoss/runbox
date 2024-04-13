package notebook

import (
	"fmt"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var router1 *gin.Engine

func init() {
	gin.SetMode(gin.TestMode)
	router1 = gin.Default()
	router1.POST("/notebook", Run)
}

func TestRun(t *testing.T) {
	testcases := []struct {
		reqBody      string
		wantCode     int
		wantResponse string
	}{
		{
			reqBody:      ``,
			wantCode:     400,
			wantResponse: `{"error":"ErrBindJSON","status":"error"}`,
		},
		{
			reqBody:      `{}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrInvalidLanguage","status":"error"}`,
		},
		{
			reqBody:      `{"asdfasdf": ""}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrInvalidLanguage","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "bash"}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrInvalidLanguage","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "", "source": "echo hello"}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrInvalidLanguage","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "_", "source": "echo hello"}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrInvalidLanguage","status":"error"}`,
		},
		{
			reqBody:      `{"lang": "bash", "source": "echo hello"}`,
			wantCode:     400,
			wantResponse: `{"error":"ErrInvalidLanguage","status":"error"}`,
		},
		// python
		{
			reqBody:      `{"lang": "python"}`,
			wantCode:     200,
			wantResponse: `{"data":{"metadata":{"kernelspec":{"name":"","display_name":""},"language_info":{"name":"python","codemirror_mode":{"name":"ipython","version":3},"file_extension":".py","mimetype":"text/x-python","pygments_lexer":"ipython3"}},"cellOutputs":null},"status":"success"}`,
		},
		{
			reqBody:  `{"lang": "python","cellTexts":[["print(\"hello1\")\n","print(\"hello2\")\n"],["print(\"world3\")\n","print(\"world4\")\n"]]}`,
			wantCode: 200,
			wantResponse: `{
				"data": {
				  "metadata": {
					"kernelspec": {
					  "name": "",
					  "display_name": ""
					},
					"language_info": {
					  "name": "python",
					  "codemirror_mode": {
						"name": "ipython",
						"version": 3
					  },
					  "file_extension": ".py",
					  "mimetype": "text/x-python",
					  "pygments_lexer": "ipython3"
					}
				  },
				  "cellOutputs": [
					[
					  {
						"output_type": "stream",
						"name": "stdout",
						"text": [
						  "hello1\n",
						  "hello2\n"
						]
					  }
					],
					[
					  {
						"output_type": "stream",
						"name": "stdout",
						"text": [
						  "world3\n",
						  "world4\n"
						]
					  }
					]
				  ]
				},
				"status": "success"
			  }`,
		},
		// r
		{
			reqBody:      `{"lang": "r"}`,
			wantCode:     200,
			wantResponse: `{"data":{"metadata":{"kernelspec":{"name":"ir","display_name":""},"language_info":{"name":"R","codemirror_mode":"r","file_extension":".r","mimetype":"text/x-r-source","pygments_lexer":"r"}},"cellOutputs":null},"status":"success"}`,
		},
		{
			reqBody:  `{"lang":"r","cellTexts":[["print(\"hello\")","print(\"world\")"],["head(iris)"]]}`,
			wantCode: 200,
			wantResponse: `{
				"data": {
				  "metadata": {
					"kernelspec": {
					  "name": "ir",
					  "display_name": ""
					},
					"language_info": {
					  "name": "R",
					  "codemirror_mode": "r",
					  "file_extension": ".r",
					  "mimetype": "text/x-r-source",
					  "pygments_lexer": "r"
					}
				  },
				  "cellOutputs": [
					[
					  {
						"output_type": "error",
						"ename": "ERROR",
						"evalue": "Error in parse(text = x, srcfile = src): <text>:1:15: unexpected symbol\n1: print(\"hello\")print\n                  ^\n",
						"traceback": [
						  "Error in parse(text = x, srcfile = src): <text>:1:15: unexpected symbol\n1: print(\"hello\")print\n                  ^\nTraceback:\n"
						]
					  }
					],
					[
					  {
						"output_type": "display_data",
						"data": {
						  "text/html": [
							"<table class=\"dataframe\">\n",
							"<caption>A data.frame: 6 × 5</caption>\n",
							"<thead>\n",
							"\t<tr><th></th><th scope=col>Sepal.Length</th><th scope=col>Sepal.Width</th><th scope=col>Petal.Length</th><th scope=col>Petal.Width</th><th scope=col>Species</th></tr>\n",
							"\t<tr><th></th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;fct&gt;</th></tr>\n",
							"</thead>\n",
							"<tbody>\n",
							"\t<tr><th scope=row>1</th><td>5.1</td><td>3.5</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n",
							"\t<tr><th scope=row>2</th><td>4.9</td><td>3.0</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n",
							"\t<tr><th scope=row>3</th><td>4.7</td><td>3.2</td><td>1.3</td><td>0.2</td><td>setosa</td></tr>\n",
							"\t<tr><th scope=row>4</th><td>4.6</td><td>3.1</td><td>1.5</td><td>0.2</td><td>setosa</td></tr>\n",
							"\t<tr><th scope=row>5</th><td>5.0</td><td>3.6</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n",
							"\t<tr><th scope=row>6</th><td>5.4</td><td>3.9</td><td>1.7</td><td>0.4</td><td>setosa</td></tr>\n",
							"</tbody>\n",
							"</table>\n"
						  ],
						  "text/latex": [
							"A data.frame: 6 × 5\n",
							"\\begin{tabular}{r|lllll}\n",
							"  & Sepal.Length & Sepal.Width & Petal.Length & Petal.Width & Species\\\\\n",
							"  & <dbl> & <dbl> & <dbl> & <dbl> & <fct>\\\\\n",
							"\\hline\n",
							"\t1 & 5.1 & 3.5 & 1.4 & 0.2 & setosa\\\\\n",
							"\t2 & 4.9 & 3.0 & 1.4 & 0.2 & setosa\\\\\n",
							"\t3 & 4.7 & 3.2 & 1.3 & 0.2 & setosa\\\\\n",
							"\t4 & 4.6 & 3.1 & 1.5 & 0.2 & setosa\\\\\n",
							"\t5 & 5.0 & 3.6 & 1.4 & 0.2 & setosa\\\\\n",
							"\t6 & 5.4 & 3.9 & 1.7 & 0.4 & setosa\\\\\n",
							"\\end{tabular}\n"
						  ],
						  "text/markdown": [
							"\n",
							"A data.frame: 6 × 5\n",
							"\n",
							"| <!--/--> | Sepal.Length &lt;dbl&gt; | Sepal.Width &lt;dbl&gt; | Petal.Length &lt;dbl&gt; | Petal.Width &lt;dbl&gt; | Species &lt;fct&gt; |\n",
							"|---|---|---|---|---|---|\n",
							"| 1 | 5.1 | 3.5 | 1.4 | 0.2 | setosa |\n",
							"| 2 | 4.9 | 3.0 | 1.4 | 0.2 | setosa |\n",
							"| 3 | 4.7 | 3.2 | 1.3 | 0.2 | setosa |\n",
							"| 4 | 4.6 | 3.1 | 1.5 | 0.2 | setosa |\n",
							"| 5 | 5.0 | 3.6 | 1.4 | 0.2 | setosa |\n",
							"| 6 | 5.4 | 3.9 | 1.7 | 0.4 | setosa |\n",
							"\n"
						  ],
						  "text/plain": [
							"  Sepal.Length Sepal.Width Petal.Length Petal.Width Species\n",
							"1 5.1          3.5         1.4          0.2         setosa \n",
							"2 4.9          3.0         1.4          0.2         setosa \n",
							"3 4.7          3.2         1.3          0.2         setosa \n",
							"4 4.6          3.1         1.5          0.2         setosa \n",
							"5 5.0          3.6         1.4          0.2         setosa \n",
							"6 5.4          3.9         1.7          0.4         setosa "
						  ]
						}
					  }
					]
				  ]
				},
				"status": "success"
			  }`,
		},
		{
			reqBody:  `{"lang":"r","cellTexts":[["library(ggplot2)","table(mpg$class)"],["head(iris)"]]}`,
			wantCode: 200,
			wantResponse: `{
				"data": {
				  "metadata": {
					"kernelspec": {
					  "name": "ir",
					  "display_name": ""
					},
					"language_info": {
					  "name": "R",
					  "codemirror_mode": "r",
					  "file_extension": ".r",
					  "mimetype": "text/x-r-source",
					  "pygments_lexer": "r"
					}
				  },
				  "cellOutputs": [
					[
					  {
						"output_type": "error",
						"ename": "ERROR",
						"evalue": "Error in parse(text = x, srcfile = src): <text>:1:17: unexpected symbol\n1: library(ggplot2)table\n                    ^\n",
						"traceback": [
						  "Error in parse(text = x, srcfile = src): <text>:1:17: unexpected symbol\n1: library(ggplot2)table\n                    ^\nTraceback:\n"
						]
					  }
					],
					[
					  {
						"output_type": "display_data",
						"data": {
						  "text/html": [
							"<table class=\"dataframe\">\n",
							"<caption>A data.frame: 6 × 5</caption>\n",
							"<thead>\n",
							"\t<tr><th></th><th scope=col>Sepal.Length</th><th scope=col>Sepal.Width</th><th scope=col>Petal.Length</th><th scope=col>Petal.Width</th><th scope=col>Species</th></tr>\n",
							"\t<tr><th></th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;dbl&gt;</th><th scope=col>&lt;fct&gt;</th></tr>\n",
							"</thead>\n",
							"<tbody>\n",
							"\t<tr><th scope=row>1</th><td>5.1</td><td>3.5</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n",
							"\t<tr><th scope=row>2</th><td>4.9</td><td>3.0</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n",
							"\t<tr><th scope=row>3</th><td>4.7</td><td>3.2</td><td>1.3</td><td>0.2</td><td>setosa</td></tr>\n",
							"\t<tr><th scope=row>4</th><td>4.6</td><td>3.1</td><td>1.5</td><td>0.2</td><td>setosa</td></tr>\n",
							"\t<tr><th scope=row>5</th><td>5.0</td><td>3.6</td><td>1.4</td><td>0.2</td><td>setosa</td></tr>\n",
							"\t<tr><th scope=row>6</th><td>5.4</td><td>3.9</td><td>1.7</td><td>0.4</td><td>setosa</td></tr>\n",
							"</tbody>\n",
							"</table>\n"
						  ],
						  "text/latex": [
							"A data.frame: 6 × 5\n",
							"\\begin{tabular}{r|lllll}\n",
							"  & Sepal.Length & Sepal.Width & Petal.Length & Petal.Width & Species\\\\\n",
							"  & <dbl> & <dbl> & <dbl> & <dbl> & <fct>\\\\\n",
							"\\hline\n",
							"\t1 & 5.1 & 3.5 & 1.4 & 0.2 & setosa\\\\\n",
							"\t2 & 4.9 & 3.0 & 1.4 & 0.2 & setosa\\\\\n",
							"\t3 & 4.7 & 3.2 & 1.3 & 0.2 & setosa\\\\\n",
							"\t4 & 4.6 & 3.1 & 1.5 & 0.2 & setosa\\\\\n",
							"\t5 & 5.0 & 3.6 & 1.4 & 0.2 & setosa\\\\\n",
							"\t6 & 5.4 & 3.9 & 1.7 & 0.4 & setosa\\\\\n",
							"\\end{tabular}\n"
						  ],
						  "text/markdown": [
							"\n",
							"A data.frame: 6 × 5\n",
							"\n",
							"| <!--/--> | Sepal.Length &lt;dbl&gt; | Sepal.Width &lt;dbl&gt; | Petal.Length &lt;dbl&gt; | Petal.Width &lt;dbl&gt; | Species &lt;fct&gt; |\n",
							"|---|---|---|---|---|---|\n",
							"| 1 | 5.1 | 3.5 | 1.4 | 0.2 | setosa |\n",
							"| 2 | 4.9 | 3.0 | 1.4 | 0.2 | setosa |\n",
							"| 3 | 4.7 | 3.2 | 1.3 | 0.2 | setosa |\n",
							"| 4 | 4.6 | 3.1 | 1.5 | 0.2 | setosa |\n",
							"| 5 | 5.0 | 3.6 | 1.4 | 0.2 | setosa |\n",
							"| 6 | 5.4 | 3.9 | 1.7 | 0.4 | setosa |\n",
							"\n"
						  ],
						  "text/plain": [
							"  Sepal.Length Sepal.Width Petal.Length Petal.Width Species\n",
							"1 5.1          3.5         1.4          0.2         setosa \n",
							"2 4.9          3.0         1.4          0.2         setosa \n",
							"3 4.7          3.2         1.3          0.2         setosa \n",
							"4 4.6          3.1         1.5          0.2         setosa \n",
							"5 5.0          3.6         1.4          0.2         setosa \n",
							"6 5.4          3.9         1.7          0.4         setosa "
						  ]
						}
					  }
					]
				  ]
				},
				"status": "success"
			  }`,
		},
	}
	for i, tc := range testcases {
		t.Run(fmt.Sprintf("#%d_%v", i, tc.reqBody), func(t *testing.T) {
			req := httptest.NewRequest("POST", "/notebook", strings.NewReader(tc.reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router1.ServeHTTP(w, req)
			require.Equal(t, tc.wantCode, w.Code)
			response := w.Body.String()
			// ignore json fields: [time, cpu, mem]
			re := regexp.MustCompile(`("time":)([^,]+)(,"cpu":)([0-9.]+)(,"mem":)([0-9.]+)`)
			response = re.ReplaceAllString(response, `${1}"0:00.00"${3}0${5}0`)
			require.JSONEq(t, tc.wantResponse, response)
		})
	}
}
