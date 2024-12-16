package testutil

import (
	"testing"

	"github.com/docker/docker/client"
	"github.com/maxatome/go-testdeep/td"
)

func TestNew(tt *testing.T) {
	t := td.NewT(tt, td.ContextConfig{IgnoreUnexported: true})
	want := &client.Client{}
	got := NewDocker()
	t.Cmp(got, want)
}
