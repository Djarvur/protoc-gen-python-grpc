package kit_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/kit"
)

var (
	errGenerator = errors.New("generate error")
	errReader    = errors.New("read error")
	errWriter    = errors.New("write error")
)

func TestRunPluginWithIO_ReadError(t *testing.T) {
	t.Parallel()

	// Arrange
	gen := newMockGenerator()
	gen.On("Generate", mock.Anything).Return(nil)

	reader := newMockReader()
	reader.On("Read", mock.Anything).Return(0, errReader)

	out := &bytes.Buffer{}

	// Act
	err := kit.New().RunPluginWithIO(gen, reader, out)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, kit.ErrRun)
	require.ErrorIs(t, err, errReader)
}

func TestRunPluginWithIO_BadInputDataError(t *testing.T) {
	t.Parallel()

	// Arrange
	gen := newMockGenerator()
	gen.On("Generate", mock.Anything).Return(nil)

	bad := []byte{0x00}
	out := &bytes.Buffer{}

	// Act
	err := kit.New().RunPluginWithIO(gen, bytes.NewBuffer(bad), out)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, kit.ErrRun)
	require.ErrorIs(t, err, proto.Error)
}

func TestRunPluginWithIO_NoFileToGenerateError(t *testing.T) {
	t.Parallel()

	// Arrange
	gen := newMockGenerator()
	gen.On("Generate", mock.Anything).Return(new(pluginpb.CodeGeneratorResponse), nil)

	req := new(pluginpb.CodeGeneratorRequest)

	inBytes, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("error marshaling generator request")
	}

	out := &bytes.Buffer{}

	// Act
	err = kit.New().RunPluginWithIO(gen, bytes.NewBuffer(inBytes), out)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, kit.ErrRun)
	require.ErrorContains(t, err, "no files were supplied to the generator")
}

func TestRunPluginWithIO_WriteError(t *testing.T) {
	t.Parallel()

	// Arrange
	gen := newMockGenerator()
	gen.On("Generate", mock.Anything).Return(new(pluginpb.CodeGeneratorResponse), nil)

	req := new(pluginpb.CodeGeneratorRequest)
	req.FileToGenerate = []string{""}

	inBytes, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("error marshaling generator request")
	}

	writer := newMockWriter()
	writer.On("Write", mock.Anything).Return(0, errWriter)

	// Act
	err = kit.New().RunPluginWithIO(gen, bytes.NewBuffer(inBytes), writer)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, kit.ErrRun)
	require.ErrorIs(t, err, errWriter)
}

func TestRunPluginWithIO_GenerateError(t *testing.T) {
	t.Parallel()

	// Arrange
	gen := newMockGenerator()
	gen.On("Generate", mock.Anything).Return((*pluginpb.CodeGeneratorResponse)(nil), errGenerator)

	req := new(pluginpb.CodeGeneratorRequest)
	req.FileToGenerate = []string{""}

	inBytes, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("error marshaling generator request")
	}

	out := &bytes.Buffer{}

	// Act
	err = kit.New().RunPluginWithIO(gen, bytes.NewBuffer(inBytes), out)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, kit.ErrRun)
	require.ErrorIs(t, err, errGenerator)
}

func TestRunPluginWithIO_Success(t *testing.T) {
	t.Parallel()

	// Arrange
	req := new(pluginpb.CodeGeneratorRequest)
	req.FileToGenerate = []string{""}

	inBytes, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("error marshaling generator request")
	}

	expectedResp := new(pluginpb.CodeGeneratorResponse)
	expectedResp.SupportedFeatures = proto.Uint64(uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL))

	outBytes, err := proto.Marshal(expectedResp)
	if err != nil {
		t.Fatalf("error marshaling expected response: %v", err)
	}

	gen := newMockGenerator()
	gen.On("Generate", mock.Anything).Return(expectedResp, nil)

	out := &bytes.Buffer{}

	// Act
	err = kit.New().RunPluginWithIO(gen, bytes.NewBuffer(inBytes), out)

	// Assert
	require.NoError(t, err)
	require.Equal(t, out.Bytes(), outBytes)
}

type mockReader struct{ mock.Mock }

func newMockReader() *mockReader {
	return &mockReader{}
}

func (mr *mockReader) Read(p []byte) (int, error) {
	args := mr.Called(p)

	return args.Int(0), args.Error(1)
}

type mockWriter struct{ mock.Mock }

func newMockWriter() *mockWriter {
	return &mockWriter{}
}

func (mw *mockWriter) Write(p []byte) (int, error) {
	args := mw.Called(p)

	return args.Int(0), args.Error(1)
}

type mockGenerator struct{ mock.Mock }

func newMockGenerator() *mockGenerator {
	return &mockGenerator{}
}

func (mg *mockGenerator) Generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	args := mg.Called(req)
	//nolint:forcetypeassert
	return args.Get(0).(*pluginpb.CodeGeneratorResponse), args.Error(1)
}
