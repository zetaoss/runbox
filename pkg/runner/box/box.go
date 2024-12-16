package box

import (
	"github.com/docker/docker/client"
)

type Box struct {
	cli *client.Client
}

func New(cli *client.Client) *Box {
	return &Box{cli}
}

func (b *Box) Run(opts *Opts) (*Result, error) {
	s := NewSession(b.cli, opts)
	if err := s.run(); err != nil {
		return nil, err
	}
	return &s.result, nil
}
