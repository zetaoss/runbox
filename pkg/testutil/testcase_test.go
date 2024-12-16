package testutil

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestName(t *testing.T) {
	tests := []struct {
		input    []any
		expected string
	}{
		{[]any{1, "hello", "world"}, "01 hello world"},
		{[]any{12, "test", "example"}, "12 test example"},
		{[]any{123, "A/B", "C_D"}, "123 A%B C D"},
		{[]any{99, "this", "is", "a", "very", "long", "string", "that", "should", "be", "truncated", "to", "64", "characters", "in", "total"}, "99 this is a very long string that should be truncated to 64 ..."},
		{[]any{0, "test", nil, "example"}, "00 test example"},
		{[]any{"string only", "no number"}, "string only no number"},
	}

	for i, tc := range tests {
		t.Run(Name(i, tc.input), func(t *testing.T) {
			result := Name(tc.input...)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		input any
		want  string
	}{
		{
			map[string]any{"key1": "value1", "key2": "value2"},
			"value1 value2",
		},
		{
			[]any{1, "test", true},
			"1 test true",
		},
		{
			struct {
				A int
				B string
			}{
				A: 1,
				B: "test",
			},
			"1 test",
		},
		{
			struct {
				C any
			}{
				C: []any{"nested", "slice"},
			},
			"nested slice",
		},
		{
			nil,
			"",
		},
		{
			123,
			"123",
		},
	}

	for i, tc := range tests {
		t.Run(Name(i, tc.input), func(t *testing.T) {
			got := toString(tc.input)
			fmt.Printf("%#v\n", got)
			assert.Equal(t, tc.want, got)
		})
	}
}
