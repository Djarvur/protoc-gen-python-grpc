package kit_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/generator"
	"github.com/Djarvur/protoc-gen-python-grpc/internal/kit"
	"github.com/Djarvur/protoc-gen-python-grpc/internal/template"
	"github.com/stretchr/testify/require"
)

const (
	fileNameSuffixDefault = "_pb2_grpc.py"
)

var (
	//go:embed testdata/in.bin
	inBytes []byte

	//go:embed testdata/out.bin
	outBytes []byte
)

func TestRunPluginWithIO(t *testing.T) {
	t.Parallel()

	out := &bytes.Buffer{}

	err := kit.New().RunPluginWithIO(
		generator.New(fileNameSuffixDefault, template.DefaultValue().Source()),
		bytes.NewBuffer(inBytes),
		out,
	)
	require.NoError(t, err)

	require.Equal(t, out.Bytes(), outBytes)
}
