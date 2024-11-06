package template_test

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/template"
)

var errReader = errors.New("read error")

func TestDefaultValue(t *testing.T) {
	t.Parallel()

	val := template.DefaultValue()

	require.NotNil(t, val)
	require.Equal(t, "EMBEDDED", val.Name())
	require.NotEmpty(t, val.Source())
}

func TestNewValue(t *testing.T) {
	t.Parallel()

	name := "test"
	source := "test template content"
	reader := strings.NewReader(source)

	val, err := template.NewValue(name, reader)

	require.NoError(t, err)
	require.NotNil(t, val)
	require.Equal(t, name, val.Name())
	require.Equal(t, source, val.Source())
}

func TestNewValue_Error(t *testing.T) {
	t.Parallel()

	reader := newMockReader()
	reader.On("Read", mock.Anything).Return(0, errReader)

	_, err := template.NewValue("test", reader)

	require.Error(t, err)
	require.ErrorIs(t, err, template.ErrTemplateRead)
	require.ErrorIs(t, err, errReader)
}

func TestValue_Set(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp("", "*_template_test")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	source := "file template content"
	_, err = file.WriteString(source)
	require.NoError(t, err)
	require.NoError(t, file.Close())

	val := &template.Value{}

	err = val.Set(file.Name())

	require.NoError(t, err)
	require.Equal(t, file.Name(), val.Name())
	require.Equal(t, source, val.Source())
}

func TestValue_Set_Error(t *testing.T) {
	t.Parallel()

	val := &template.Value{}

	err := val.Set("non_existent_file")

	require.Error(t, err)
	require.ErrorIs(t, err, template.ErrTemplateRead)
}

func TestValue_Type(t *testing.T) {
	t.Parallel()

	val := &template.Value{}

	require.Equal(t, "text/template", val.Type())
}

type mockReader struct{ mock.Mock }

func newMockReader() *mockReader {
	return &mockReader{}
}

func (mr *mockReader) Read(p []byte) (int, error) {
	args := mr.Called(p)

	return args.Int(0), args.Error(1)
}
