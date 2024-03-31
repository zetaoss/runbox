package main

import "github.com/zetaoss/runbox/pkg/handler"

func main() {
	r := handler.NewRouter()
	err := r.Run(":8080")
	panic(err)
}
