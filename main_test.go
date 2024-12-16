package main

import (
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	go main()
	time.Sleep(time.Duration(100) * time.Millisecond)
}
