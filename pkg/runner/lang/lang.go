package lang

import (
	"fmt"
	"strings"

	"github.com/zetaoss/runbox/pkg/errors"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"k8s.io/utils/ptr"
)

type Lang struct {
	box *box.Box
}

func New(box *box.Box) *Lang {
	return &Lang{box}
}

type Input struct {
	Lang  string     `json:"lang"`
	Files []box.File `json:"files"`
	Main  int        `json:"main,omitempty"`
}

type LangOpts struct {
	Input              Input
	Command            string
	CollectImagesCount int
	Env                []string
	FileDir            string
	FileName           string
	FileExt            string
	FileMain           int
	ModifyMainFunc     func(string) string
	Shell              string
	TimeoutSeconds     int
	User               string
	WorkingDir         string
}

func (l *Lang) Run(input Input, extraOpts ...map[string]int) (*box.Result, error) {
	langOpts, err := toLangOpts(input)
	if err != nil {
		if err == errors.ErrInvalidLanguage || err == errors.ErrNoFiles {
			return nil, err
		}
		return nil, fmt.Errorf("toLangOpts err: %w", err)
	}
	for _, o := range extraOpts {
		for k, v := range o {
			if k == "timeoutSeconds" {
				langOpts.TimeoutSeconds = v + 1
			}
		}
	}
	boxOpts := toBoxOpts(*langOpts)
	return l.box.Run(&boxOpts)
}

func toBoxOpts(langOpts LangOpts) box.Opts {
	var files []box.File

	for i, f := range langOpts.Input.Files {
		fileName := f.Name
		if fileName == "" {
			fileName = langOpts.FileName + "." + langOpts.FileExt
		}
		if i == langOpts.Input.Main && langOpts.ModifyMainFunc != nil {
			f.Body = langOpts.ModifyMainFunc(f.Body)
		}
		file := box.File{
			Name: langOpts.WorkingDir + langOpts.FileDir + "/" + fileName,
			Body: f.Body,
		}
		files = append(files, file)
	}
	return box.Opts{
		CollectStats:       ptr.To(true),
		CollectImages:      true,
		CollectImagesCount: langOpts.CollectImagesCount,
		Command:            langOpts.Command,
		Env:                langOpts.Env,
		Files:              files,
		Image:              fmt.Sprintf("ghcr.io/zetaoss/runcontainers/%s", langOpts.Input.Lang),
		Shell:              langOpts.Shell,
		Timeout:            langOpts.TimeoutSeconds * 1000,
		User:               langOpts.User,
		WorkingDir:         langOpts.WorkingDir,
	}
}

func toLangOpts(input Input) (*LangOpts, error) {
	if len(input.Files) == 0 {
		return nil, errors.ErrNoFiles
	}
	var opts = &LangOpts{
		Input:          input,
		Command:        "",
		FileName:       "runbox",
		FileExt:        input.Lang,
		FileDir:        "",
		ModifyMainFunc: nil,
		Shell:          "sh",
		TimeoutSeconds: 10,
		WorkingDir:     "/home/user01",
	}
	switch input.Lang {
	case "bash":
		opts.Command = "/bin/bash runbox.sh"
		opts.FileExt = "sh"
		opts.Shell = "bash"
	case "c":
		opts.Command = "gcc runbox.c; ./a.out"
	case "cpp":
		opts.Command = "g++ runbox.cpp; ./a.out"
	case "csharp":
		opts.Command = "mcs runbox.cs; mono runbox.exe"
		opts.FileExt = "cs"
	case "java":
		opts.Command = `javac -d bin -cp "lib/*" src/*; java -cp "bin:lib/*" App`
		opts.FileDir = "/src"
		opts.FileName = "App"
		opts.WorkingDir = "/demo"
	case "kotlin":
		opts.Command = "kotlinc runbox.kt -include-runtime -d runbox.jar && java -jar runbox.jar"
		opts.FileExt = "kt"
		opts.TimeoutSeconds = 40
	case "go":
		opts.Command = "go run runbox.go"
		opts.Env = []string{"TINI_SUBREAPER=1"}
		opts.TimeoutSeconds = 30
	case "lua":
		opts.Command = "lua runbox.lua"
	case "mysql":
		opts.Command = "bash /tmp/entrypoint.sh"
		opts.FileExt = "sql"
		opts.TimeoutSeconds = 30
	case "perl":
		opts.Command = "perl runbox.pl"
		opts.FileExt = "pl"
	case "php":
		opts.Command = "php runbox.php"
		opts.ModifyMainFunc = func(source string) string {
			source = strings.TrimLeft(source, " \t\n")
			if source[:5] != "<?php" {
				source = "<?php\nrequire_once('vendor/autoload.php');\n" + source
			}
			return source
		}
	case "powershell":
		opts.Command = "pwsh runbox.ps"
		opts.FileExt = "ps"
	case "python":
		opts.Command = "python runbox.py"
		opts.FileExt = "py"
	case "r":
		opts.Command = "Rscript runbox.r"
		opts.ModifyMainFunc = func(source string) string {
			return "png(width=500,height=400);\n" + source + "\n" + `options(echo=F); invisible(dev.off());system('find . -name "*.pdf" -exec mogrify -density 80 -format png {} \\;',ignore.stdout=T,ignore.stderr=F);`
		}
	case "ruby":
		opts.Command = "ruby runbox.rb"
		opts.FileExt = "rb"
	case "sqlite3":
		opts.FileExt = "sql"
		source := input.Files[0].Body
		if strings.HasPrefix(source, ".") {
			opts.Command = "sqlite3 chinook.db " + source
		} else {
			opts.Command = "sqlite3 -header chinook.db < runbox.sql"
		}
		input.Files[0].Body = source
	case "tex":
		opts.FileExt = "tex"
		opts.CollectImagesCount = 10
		opts.Command = "touch oblivoir.sty && pdflatex -halt-on-error runbox.tex && convert runbox.pdf -strip p%d.png"
		opts.TimeoutSeconds = 30
		opts.User = "root"
	default:
		return nil, errors.ErrInvalidLanguage
	}
	return opts, nil
}
