package multi

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zetaoss/runbox/pkg/docker"
	"github.com/zetaoss/runbox/pkg/runner/lang/types"
	"github.com/zetaoss/runbox/pkg/util/runid"
	"k8s.io/klog/v2"
)

func getOutputAndLogs(runResult *docker.Result, runID string) (*types.Output, []docker.Log) {
	output := &types.Output{}
	if runResult.ExitCode == 124 {
		output.Warnings = append(output.Warnings, types.WarnTimeout)
	}
	if runResult.OutputLimitReached {
		output.Warnings = append(output.Warnings, types.WarnOutputLimitReached)
	}
	logs := runResult.Logs
	var timeIndex int = -1
	for i, log := range logs {
		if strings.HasPrefix(log.Log, "@@"+runID+"@@") {
			timeIndex = i
			pieces := strings.Split(log.Log, "@@")
			if len(pieces) < 3 {
				return output, logs
			}
			pieces = strings.Split(pieces[2], ",")
			if len(pieces) < 3 {
				return output, logs
			}
			output.Time = pieces[0]
			var floatValue float64
			floatValue, _ = strconv.ParseFloat(pieces[1], 32)
			output.CPU = float32(floatValue)
			floatValue, _ = strconv.ParseFloat(pieces[2], 32)
			output.MEM = float32(floatValue)
			break
		}
	}
	if timeIndex != -1 {
		logs = append(logs[:timeIndex], logs[timeIndex+1:]...)
	}
	return output, logs
}

func toOutput(runResult *docker.Result, runID string) *types.Output {
	output, logs := getOutputAndLogs(runResult, runID)

	var sepIndex int = -1
	for i, log := range logs {
		if log.Log == "=="+runID+"==\n" {
			sepIndex = i
		}
	}
	if sepIndex != -1 {
		if sepIndex < len(logs)-1 { // one or more images exist
			for _, log := range logs[sepIndex+2:] {
				output.Images = append(output.Images, log.Log)
			}
		}
		logs = logs[:sepIndex-1]
	}
	last := len(logs) - 1
	output.Logs = []string{}
	for i, log := range logs {
		line := log.Log
		if i != last {
			line = line[:len(line)-1]
		} else if line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
		if log.Stream == "stderr" {
			line = "1" + line
		} else {
			line = "0" + line
		}
		output.Logs = append(output.Logs, line)
	}
	return output
}

//gocyclo:ignore
func getRunOpts(input Input) (*types.RunOpts, error) {
	var opts = &types.RunOpts{
		Command:          "",
		Env:              []string{},
		FileName:         "runbox",
		FileExt:          input.Lang,
		ModifySourceFunc: nil,
		Postflight:       nil,
		Shell:            "sh",
		TimeoutCommand:   "timeout --kill-after=1",
		TimeoutSeconds:   10,
		VolSubPath:       "",
		WorkingDir:       "/home/user01",
	}
	var langFunc = map[string]func(*types.RunOpts){
		"bash": func(*types.RunOpts) {
			opts.Command = "/bin/bash runbox.sh"
			opts.FileExt = "sh"
			opts.Shell = "bash"
		},
		"c": func(*types.RunOpts) {
			opts.Command = "gcc runbox.c; ./a.out"
		},
		"cpp": func(*types.RunOpts) {
			opts.Command = "gcc -lstdc++ runbox.cpp; ./a.out"
		},
		"csharp": func(*types.RunOpts) {
			opts.Command = "mcs runbox.cs; mono runbox.exe"
			opts.FileExt = "cs"
		},
		"java": func(*types.RunOpts) {
			opts.Command = `javac -d bin -cp "lib/*" src/*; java -cp "bin:lib/*" App && echo && echo ==` + input.RunID + `== && ls *.png 2>/dev/null | head -2 | xargs -i sh -c "echo; base64 -w0 {}"`
			opts.FileName = "App"
			opts.VolSubPath = "/src"
			opts.WorkingDir = "/demo"
		},
		"latex": func(*types.RunOpts) {
			opts.Command = `touch oblivoir.sty && pdflatex -halt-on-error runbox.tex && convert -strip runbox.pdf p%d.png && echo ==` + input.RunID + `== && ls *.png 2>/dev/null | head -10 | xargs -i sh -c "echo; base64 -w0 {}"`
			opts.Image = fmt.Sprintf("jmnote/runbox:tex")
			opts.FileExt = "tex"
			opts.Postflight = func(o *types.Output) {
				if len(o.Images) > 0 {
					o.Logs = []string{}
				}
			}
		},
		"kotlin": func(*types.RunOpts) {
			opts.Command = "kotlinc runbox.kt -d runbox.jar && TIME java -jar runbox.jar" // 13.6s 17.2s
			opts.FileExt = "kt"
			opts.TimeoutCommand = "timeout -s KILL"
			opts.TimeoutSeconds = 30
		},
		"go": func(*types.RunOpts) {
			opts.Command = "go mod tidy 2>/dev/null && TIME go run runbox.go"
			opts.Env = []string{"TINI_SUBREAPER=1"}
			opts.TimeoutCommand = "tini timeout"
			opts.TimeoutSeconds = 30
		},
		"lua": func(*types.RunOpts) {
			opts.Command = "lua runbox.lua"
		},
		"mysql": func(*types.RunOpts) {
			opts.Command = "bash /tmp/entrypoint.sh"
			opts.FileExt = "sql"
			opts.TimeoutSeconds = 30
		},
		"perl": func(*types.RunOpts) {
			opts.Command = "perl runbox.pl"
			opts.FileExt = "pl"
		},
		"php": func(*types.RunOpts) {
			opts.Command = "php runbox.php"
			opts.ModifySourceFunc = func(source string) string {
				source = strings.TrimLeft(source, " \t\n")
				if source[:5] != "<?php" {
					source = "<?php\nrequire_once('vendor/autoload.php');\n" + source
				}
				return source
			}
		},
		"powershell": func(*types.RunOpts) {
			opts.Command = "pwsh runbox.ps"
			opts.FileExt = "ps"
		},
		"python": func(*types.RunOpts) {
			opts.Command = "python runbox.py"
			opts.FileExt = "py"
		},
		"r": func(*types.RunOpts) {
			opts.Command = "Rscript runbox.r"
			opts.ModifySourceFunc = func(source string) string {
				return "png(width=500,height=400);\n" + source + "\n" +
					`cat('\n==` + input.RunID + `==\n'); options(echo=F); invisible(dev.off());` +
					`system('find . -name "*.pdf" -exec mogrify -density 80 -format png {} \\;',ignore.stdout=T,ignore.stderr=F);` +
					`system('ls Rplot00?.png 2>/dev/null | head -2 | xargs -i sh -c "echo; base64 -w0 {}"')`
			}
		},
		"ruby": func(*types.RunOpts) {
			opts.Command = "ruby runbox.rb"
			opts.FileExt = "rb"
		},
		"sqlite3": func(*types.RunOpts) {
			source := input.Files[0].Text
			input.Files[0].Text = source
			if strings.HasPrefix(source, ".") {
				opts.Command = "sqlite3 chinook.db " + source
			} else {
				opts.Command = "sqlite3 -header chinook.db < runbox.sql"
			}
			opts.FileExt = "sql"
		},
		"tex": func(*types.RunOpts) {
			opts.Command = `touch oblivoir.sty && pdflatex -halt-on-error runbox.tex && convert -strip runbox.pdf p%d.png && echo ==` + input.RunID + `== && ls *.png 2>/dev/null | head -10 | xargs -i sh -c "echo; base64 -w0 {}"`
			opts.FileExt = "tex"
			opts.Postflight = func(o *types.Output) {
				if len(o.Images) > 0 {
					o.Logs = []string{}
				}
			}
		},
	}
	f, ok := langFunc[input.Lang]
	if !ok {
		return nil, ErrInvalidLanguage
	}
	f(opts)
	return opts, nil
}

func Run(input Input, extraOpts ...map[string]int) (*types.Output, error) {
	if len(input.Files) == 0 {
		return nil, ErrNoFiles
	}
	if input.RunID == "" {
		input.RunID = runid.New("multi", input.Lang)
	}
	opts, err := getRunOpts(input)
	if err != nil {
		if err != ErrInvalidLanguage {
			klog.Warningf("unknown err: %s", err.Error())
		}
		return nil, err
	}
	if opts.Image == "" {
		opts.Image = fmt.Sprintf("jmnote/runbox:%s", input.Lang)
	}

	// override extraOpts
	for _, o := range extraOpts {
		for k, v := range o {
			if k == "timeoutSeconds" {
				opts.TimeoutSeconds = v + 1
			}
		}
	}

	// write files
	binds, err := writeFiles(input, opts)
	if err != nil {
		return nil, fmt.Errorf("writeFiles err: %w", err)
	}

	// timedCommand
	var command = opts.Command
	if !strings.Contains(command, "TIME") {
		command = "TIME " + command
	}
	command = strings.Replace(command, "TIME", "/usr/bin/time -f @@"+input.RunID+"@@%E,%U,%S,%M@@", 1)
	command = fmt.Sprintf("%s %d %s -c '%s'", opts.TimeoutCommand, opts.TimeoutSeconds, opts.Shell, command)
	dockerOptions := docker.Options{
		RunID:          input.RunID,
		PidsLimit:      300,
		TimeoutSeconds: 60,
		OutputLimit:    8000,
		Image:          opts.Image,
		Command:        command,
		Binds:          binds,
		Env:            opts.Env,
		WorkingDir:     opts.WorkingDir,
	}
	cli, err := docker.New()
	if err != nil {
		return nil, fmt.Errorf("docker new err: %w", err)
	}
	runResult, err := cli.Run(dockerOptions)
	if err != nil {
		return nil, fmt.Errorf("Run err: %w", err)
	}
	output := toOutput(runResult, input.RunID)
	if opts.Postflight != nil {
		opts.Postflight(output)
	}
	return output, nil
}

func writeFiles(input Input, opts *types.RunOpts) ([]string, error) {
	bindSrcRoot := "/data/files/" + input.RunID
	bindSrcDir := bindSrcRoot + opts.VolSubPath
	bindDstDir := opts.WorkingDir + opts.VolSubPath
	if err := os.MkdirAll(bindSrcDir, 0777); err != nil {
		return nil, fmt.Errorf("MkdirAll err: %w, name: %s", err, bindSrcDir)
	}
	var binds []string
	for _, f := range input.Files {
		fileName := f.Name
		content := f.Text
		if fileName == "" {
			fileName = opts.FileName + "." + opts.FileExt
		}
		if opts.ModifySourceFunc != nil && filepath.Ext(fileName) == "."+opts.FileExt {
			content = opts.ModifySourceFunc(content)
		}
		src := bindSrcDir + "/" + fileName
		dst := bindDstDir + "/" + fileName
		binds = append(binds, src+":"+dst)
		if err := os.WriteFile(src, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("WriteFile err: %w, src: %s", err, src)
		}
	}
	return binds, nil
}
