package template

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

//go:embed pb2_grpc.py.tmpl
var defaultTemplateSrc string

var _ pflag.Value = (*TemplateValue)(nil)

type TemplateValue struct {
	name   string
	source string
}

func (v *TemplateValue) String() string {
	return v.name
}

func NewTemplateValue() *TemplateValue {
	return &TemplateValue{
		name:   "EMBEDDED",
		source: defaultTemplateSrc,
	}
}

// Set is a method to set the template value.
func (v *TemplateValue) Set(s string) error {
	b, err := os.ReadFile(s)
	if err != nil {
		return fmt.Errorf("reading template %q: %w", s, err)
	}

	v.source = string(b)

	return nil
}

// Type required to implement pflag.Value.
func (v *TemplateValue) Type() string {
	return "text/template"
}

func (v *TemplateValue) Source() string {
	return v.source
}
