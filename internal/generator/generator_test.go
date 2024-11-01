package generator_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/generator"
	"github.com/Djarvur/protoc-gen-python-grpc/internal/template"
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

func TestGenerator_New(t *testing.T) {
	t.Parallel()

	// Arrange
	suffix := "suffix"
	tmplSrc := "template"

	// Act
	gen := generator.New(suffix, tmplSrc)

	// Assert
	require.NotNil(t, gen, "got nil generator")
	assert.Equal(t, suffix, gen.Suffix)
	assert.Equal(t, tmplSrc, gen.Template)
}

func TestGenerator_Generate_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	req := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(inBytes, req); err != nil {
		t.Fatalf("error unmarshalling testdata/in.bin: %v", err)
	}

	templateSource := template.NewTemplateValue()
	gen := generator.New(fileNameSuffixDefault, templateSource.Source())

	// Act
	resp, err := gen.Generate(req)

	// Assert
	require.NoError(t, err, "error generating response")
	require.NotNil(t, resp, "got nil response")

	outData, err := proto.Marshal(resp)
	if err != nil {
		t.Fatalf("error marshalling generated response: %v", err)
	}

	require.Equal(t, outBytes, outData)
}

func TestGenerator_Generate_EmptyRequestSuccess(t *testing.T) {
	t.Parallel()

	// Arrange
	req := new(pluginpb.CodeGeneratorRequest)

	expectedResp := new(pluginpb.CodeGeneratorResponse)
	expectedResp.SupportedFeatures = proto.Uint64(generator.SupportedFeatures)

	expectedData, err := proto.Marshal(expectedResp)
	if err != nil {
		t.Fatalf("error marshaling expected response: %v", err)
	}

	gen := generator.New("", "")

	// Act
	resp, err := gen.Generate(req)

	// Assert
	require.NoError(t, err, "error generating response")
	require.NotNil(t, resp, "got nil response")

	outData, err := proto.Marshal(resp)
	if err != nil {
		t.Fatalf("error marshalling generated response: %v", err)
	}

	require.Equal(t, expectedData, outData)
}

func TestGenerator_Generate_BuildTemplateError(t *testing.T) {
	t.Parallel()

	// Arrange
	req := new(pluginpb.CodeGeneratorRequest)
	gen := generator.New("", "{{ if }}")

	// Act
	_, err := gen.Generate(req)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, generator.ErrTemplateBuild)
}

func TestGenerator_Generate_ExecuteTemplateError(t *testing.T) {
	t.Parallel()

	// Arrange
	req := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(inBytes, req); err != nil {
		t.Errorf("error unmarshalling testdata/in.bin: %v", err)
	}

	gen := generator.New(fileNameSuffixDefault, "{{ .Data }}")

	// Act
	_, err := gen.Generate(req)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, generator.ErrTemplateExec)
}
