package testutil

import (
	"github.com/docker/docker/client"
	"github.com/zetaoss/runbox/pkg/docker"
)

func NewDocker() *client.Client {
	cli, err := docker.New()
	if err != nil {
		panic(err)
	}
	return cli
}
