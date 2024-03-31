package main

import "github.com/zetaoss/zetarun/pkg/handler"

func main() {
	r := handler.NewRouter()
	r.Run(":8080")
}
