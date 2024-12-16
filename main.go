package main

import (
	"github.com/zetaoss/runbox/pkg/handler"
	"github.com/zetaoss/runbox/pkg/runner/box"
	"github.com/zetaoss/runbox/pkg/runner/lang"
	"github.com/zetaoss/runbox/pkg/testutil"
)

func main() {
	d := testutil.NewDocker()
	b := box.New(d)
	langRunner := lang.New(b)
	r := handler.New(langRunner)
	_ = r.Run(":8080")
}