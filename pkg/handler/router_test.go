package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRouter(t *testing.T) {
	r := NewRouter()
	require.NotNil(t, r)
}
