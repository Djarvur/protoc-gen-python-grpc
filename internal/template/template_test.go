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

	// Arrange

	// Act
	val := template.DefaultValue()

	// Assert
	require.NotNil(t, val)
	require.Equal(t, "EMBEDDED", val.Name())
	require.NotEmpty(t, val.Source())
}

func TestNewValue(t *testing.T) {
	t.Parallel()

	// Arrange
	name := "test"
	source := "test template content"
	reader := strings.NewReader(source)

	// Act
	val, err := template.NewValue(name, reader)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, val)
	require.Equal(t, name, val.Name())
	require.Equal(t, source, val.Source())
}

func TestNewValue_Error(t *testing.T) {
	t.Parallel()

	// Arrange
	reader := newMockReader()
	reader.On("Read", mock.Anything).Return(0, errReader)

	// Act
	_, err := template.NewValue("test", reader)

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, template.ErrTemplateRead)
	require.ErrorIs(t, err, errReader)
}

func TestValue_Set(t *testing.T) {
	t.Parallel()

	// Arrange
	file, err := os.CreateTemp("", "*_template_test")
	if err != nil {
		t.Fatalf("error creating temp file: %v", err)
	}
	defer os.Remove(file.Name())

	source := "file template content"

	_, err = file.WriteString(source)
	if err != nil {
		t.Fatalf("error writing to temp file: %v", err)
	}

	err = file.Close()
	if err != nil {
		t.Fatalf("error closing temp file: %v", err)
	}

	val := &template.Value{}

	// Act
	err = val.Set(file.Name())

	// Assert
	require.NoError(t, err)
	require.Equal(t, file.Name(), val.Name())
	require.Equal(t, source, val.Source())
}

func TestValue_Set_Error(t *testing.T) {
	t.Parallel()

	// Arrange
	val := &template.Value{}

	// Act
	err := val.Set("non_existent_file")

	// Assert
	require.Error(t, err)
	require.ErrorIs(t, err, template.ErrTemplateRead)
}

func TestValue_Type(t *testing.T) {
	t.Parallel()

	// Arrange
	val := &template.Value{}

	// Act
	valType := val.Type()

	// Assert
	require.Equal(t, "text/template", valType)
}

type mockReader struct{ mock.Mock }

func newMockReader() *mockReader {
	return &mockReader{}
}

func (mr *mockReader) Read(p []byte) (int, error) {
	args := mr.Called(p)

	return args.Int(0), args.Error(1)
}
