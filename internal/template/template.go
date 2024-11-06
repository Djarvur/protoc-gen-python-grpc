package template

import (
	_ "embed"
	"errors"
	"io"
	"os"

	"github.com/spf13/pflag"
)

//go:embed pb2_grpc.py.tmpl
var defaultTemplateSrc string

var ErrTemplateRead = errors.New("reading template")

var _ pflag.Value = (*Value)(nil)

type Value struct {
	name   string
	source string
}

func DefaultValue() *Value {
	return &Value{
		name:   "EMBEDDED",
		source: defaultTemplateSrc,
	}
}

// NewValue is a method to create a new Value from a reader.
func NewValue(name string, reader io.Reader) (*Value, error) {
	val := &Value{
		name:   name,
		source: "",
	}

	source, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(ErrTemplateRead, err)
	}

	val.source = string(source)

	return val, nil
}

func (v *Value) Name() string {
	return v.name
}

func (v *Value) Source() string {
	return v.source
}

// String required to implement pflag.Value.
func (v *Value) String() string {
	return v.Name()
}

// Set required to implement pflag.Value.
// It sets the template value from a file path.
func (v *Value) Set(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return errors.Join(ErrTemplateRead, err)
	}

	newValue, err := NewValue(filepath, file)
	if err != nil {
		return err
	}

	v.name = newValue.Name()
	v.source = newValue.Source()

	return nil
}

// Type required to implement pflag.Value.
func (v *Value) Type() string {
	return "text/template"
}
