package multi

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zetaoss/runbox/pkg/docker"
	"github.com/zetaoss/runbox/pkg/run/lang/types"
	"github.com/zetaoss/runbox/pkg/util/runid"
)

func getOutputAndLogs(runResult *docker.Result, runID string) (*types.Output, []docker.Log) {
	output := &types.Output{}
	logs := runResult.Logs
	if runResult.ExitCode == 124 {
		output.Timeout = true
	}
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
func getRunOpts(input types.MultiInput) (*types.RunOpts, error) {
	if len(input.Files) == 0 {
		return nil, fmt.Errorf("no files")
	}
	var opts = &types.RunOpts{
		Command:          "",
		Env:              []string{},
		FileName:         "runbox",
		FileExt:          input.Lang,
		PidsLimit:        15,
		ModifySourceFunc: nil,
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
			opts.WorkingDir = "/demo"
			opts.VolSubPath = "/src"
			opts.FileName = "App"
			opts.PidsLimit = 30
		},
		"kotlin": func(*types.RunOpts) {
			opts.Command = "kotlinc runbox.kt -include-runtime -d runbox.jar && java -jar runbox.jar"
			opts.TimeoutCommand = "timeout -s KILL"
			opts.FileExt = "kt"
			opts.PidsLimit = 50
			opts.TimeoutSeconds = 10
		},
		"go": func(*types.RunOpts) {
			opts.PidsLimit = 40
			opts.TimeoutSeconds = 30
			opts.Env = []string{"TINI_SUBREAPER=1"}
			opts.TimeoutCommand = "tini timeout"
			// opts.Command = "go mod tidy 2>/dev/null; go run runbox.go"
			opts.Command = "go run runbox.go"
		},
		"lua": func(*types.RunOpts) {
			opts.Command = "lua runbox.lua"
		},
		"mysql": func(*types.RunOpts) {
			opts.PidsLimit = 300
			opts.TimeoutSeconds = 30
			opts.FileExt = "sql"
			opts.Command = "bash /tmp/entrypoint.sh"
		},
		"perl": func(*types.RunOpts) {
			opts.FileExt = "pl"
			opts.Command = "perl runbox.pl"
		},
		"php": func(*types.RunOpts) {
			opts.ModifySourceFunc = func(source string) string {
				source = strings.TrimLeft(source, " \t\n")
				if source[:5] != "<?php" {
					source = "<?php\nrequire_once('vendor/autoload.php');\n" + source
				}
				return source
			}
			opts.Command = "php runbox.php"
		},
		"powershell": func(*types.RunOpts) {
			opts.Command = "pwsh runbox.ps"
			opts.FileExt = "ps"
			opts.PidsLimit = 50
		},
		"python": func(*types.RunOpts) {
			opts.Command = "python runbox.py"
			opts.FileExt = "py"
		},
		"r": func(*types.RunOpts) {
			opts.ModifySourceFunc = func(source string) string {
				return "png(width=500,height=400);\n" + source + "\n" +
					`cat('\n==` + input.RunID + `==\n'); options(echo=F); invisible(dev.off());` +
					`system('find . -name "*.pdf" -exec mogrify -density 80 -format png {} \\;',ignore.stdout=T,ignore.stderr=F);` +
					`system('ls Rplot00?.png 2>/dev/null | head -2 | xargs -i sh -c "echo; base64 -w0 {}"')`
			}
			opts.Command = "Rscript runbox.r"
			opts.PidsLimit = 20
		},
		"ruby": func(*types.RunOpts) {
			opts.FileExt = "rb"
			opts.Command = "ruby runbox.rb"
		},
		"sqlite3": func(*types.RunOpts) {
			opts.FileExt = "sql"
			source := input.Files[0].Content
			if strings.HasPrefix(source, ".") {
				opts.Command = "sqlite3 chinook.db " + source
			} else {
				opts.Command = "sqlite3 -header chinook.db < runbox.sql"
			}
			input.Files[0].Content = source
		},
	}
	f, ok := langFunc[input.Lang]
	if !ok {
		return nil, types.ErrInvalidLanguage
	}
	f(opts)
	return opts, nil
}

func Run(input types.MultiInput, extraOpts ...map[string]int) (*types.Output, error) {
	if input.RunID == "" {
		input.RunID = runid.New("multi", input.Lang)
	}
	opts, err := getRunOpts(input)
	if err != nil {
		if err == types.ErrInvalidLanguage {
			return nil, types.ErrInvalidLanguage
		}
		return nil, fmt.Errorf("getRunOpts err: %w", err)
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
	command := fmt.Sprintf("%s %d %s -c '/usr/bin/time -f @@%s@@%%E,%%U,%%S,%%M@@ %s'", opts.TimeoutCommand, opts.TimeoutSeconds, opts.Shell, input.RunID, opts.Command)
	dockerOptions := docker.Options{
		RunID:          input.RunID,
		PidsLimit:      opts.PidsLimit,
		TimeoutSeconds: 60,
		Image:          fmt.Sprintf("jmnote/runbox:%s", input.Lang),
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
	return toOutput(runResult, input.RunID), nil
}

func writeFiles(input types.MultiInput, opts *types.RunOpts) ([]string, error) {
	bindSrcRoot := "/data/files/" + input.RunID
	bindSrcDir := bindSrcRoot + opts.VolSubPath
	bindDstDir := opts.WorkingDir + opts.VolSubPath
	if err := os.MkdirAll(bindSrcDir, 0777); err != nil {
		return nil, fmt.Errorf("MkdirAll err: %w, name: %s", err, bindSrcDir)
	}
	var binds []string
	for _, f := range input.Files {
		fileName := f.Name
		content := f.Content
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
