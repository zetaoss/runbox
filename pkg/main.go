package main

import "github.com/zetaoss/runbox/pkg/handler"

func main() {
	r := handler.NewRouter()
	r.Run(":8080")
}
