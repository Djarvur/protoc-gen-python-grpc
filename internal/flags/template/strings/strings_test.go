package strings_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/flags/template/strings"
)

func TestTrimExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in  string
		sep string
		out string
	}{
		{in: "aaa", sep: ".", out: "aaa"},
		{in: "aaa.bbb", sep: ".", out: "aaa"},
		{in: ".bbb", sep: ".", out: ""},
		{in: ".bbb.ccc", sep: ".", out: ".bbb"},
		{in: "qqq.aaa.bbb", sep: ".", out: "qqq.aaa"},
		{in: "qqq.aaa.bbb", sep: "/", out: "qqq.aaa.bbb"},
		{in: "qqq/aaa/bbb", sep: "/", out: "qqq/aaa"},
	}

	for i, tt := range tests {
		t.Run(
			strconv.Itoa(i),
			func(t *testing.T) {
				t.Parallel()

				out := strings.TrimSuffix(tt.sep, tt.in)
				require.Equal(t, tt.out, out)
			},
		)
	}
}

func TestBaseName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in  string
		sep string
		out string
	}{
		{in: "aaa", sep: ".", out: "aaa"},
		{in: "aaa.bbb", sep: ".", out: "bbb"},
		{in: ".bbb", sep: ".", out: "bbb"},
		{in: ".bbb.ccc", sep: ".", out: "ccc"},
		{in: "qqq.aaa.bbb", sep: ".", out: "bbb"},
		{in: "qqq.aaa.bbb", sep: "/", out: "qqq.aaa.bbb"},
		{in: "qqq/aaa/bbb", sep: "/", out: "bbb"},
	}

	for i, tt := range tests {
		t.Run(
			strconv.Itoa(i),
			func(t *testing.T) {
				t.Parallel()

				out := strings.BaseName(tt.sep, tt.in)
				require.Equal(t, tt.out, out)
			},
		)
	}
}
