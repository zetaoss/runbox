package docker

import (
	"os"
	"testing"

	"github.com/docker/docker/client"
	"github.com/maxatome/go-testdeep/td"
)

func TestNew(tt *testing.T) {
	t := td.NewT(tt, td.ContextConfig{IgnoreUnexported: true})

	t.Cmp(os.Getenv("DOCKER_TLS_VERIFY"), "1")
	t.Cmp(os.Getenv("DOCKER_HOST")[:6], "tcp://")
	t.Cmp(os.Getenv("DOCKER_CERT_PATH")[:1], "/")

	want := &client.Client{}
	cli, err := New()
	t.Nil(err)
	t.Cmp(cli, want)
}
