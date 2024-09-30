package flags

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

//go:embed pb2_grpc.py.tmpl
var defaultTemplateSrc string

var _ pflag.Value = (*sourceValue)(nil)

type sourceValue struct {
	name   string
	source string
}

func (r *sourceValue) String() string {
	return r.name
}

func newTemplateValue() *sourceValue {
	return &sourceValue{
		name:   "EMBEDDED",
		source: defaultTemplateSrc,
	}
}

// Set is a method to set the template value.
func (r *sourceValue) Set(s string) error {
	b, err := os.ReadFile(s)
	if err != nil {
		return fmt.Errorf("reading template %q: %w", s, err)
	}

	r.source = string(b)

	return nil
}

// Type required to implement pflag.Value.
func (*sourceValue) Type() string {
	return "text/template"
}
