package generator_test

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/Djarvur/protokit"
	"github.com/stretchr/testify/require"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/generator"
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

	err := protokit.RunPluginWithIO(generator.New(), bytes.NewReader(inBytes), out)
	require.NoError(t, err)
	require.Equal(t, out.Bytes(), outBytes)
}
