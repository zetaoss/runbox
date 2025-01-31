package apperror

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsAppError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "direct ErrNoFiles",
			err:  ErrNoFiles,
			want: true,
		},
		{
			name: "direct ErrNoSources",
			err:  ErrNoSources,
			want: true,
		},
		{
			name: "direct ErrInvalidLanguage",
			err:  ErrInvalidLanguage,
			want: true,
		},
		{
			name: "wrapped ErrNoSources",
			err:  fmt.Errorf("wrapped: %w", ErrNoSources),
			want: true,
		},
		{
			name: "wrapped ErrInvalidLanguage",
			err:  fmt.Errorf("another wrap: %w", fmt.Errorf("wrapped: %w", ErrInvalidLanguage)),
			want: true,
		},
		{
			name: "random error",
			err:  errors.New("some other error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsAppError(tt.err)
			if got != tt.want {
				t.Errorf("IsAppError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
