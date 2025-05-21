package template

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"text/template"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/flags/template/strings"
)

//go:embed pb2_grpc.py.tmpl
var defaultTemplateSrc string

var _ flag.Value = (*TemplateValue)(nil)

type TemplateValue struct {
	name     string
	Template *template.Template
}

func (r *TemplateValue) String() string {
	return r.name
}

func NewTemplateValue() *TemplateValue {
	return &TemplateValue{
		name:     "EMBEDDED",
		Template: template.Must(buildTemplate("EMBEDDED", defaultTemplateSrc)),
	}
}

// Set is a method to set the template value.
func (r *TemplateValue) Set(s string) error {
	b, err := os.ReadFile(s)
	if err != nil {
		return fmt.Errorf("template %q: reading: %w", s, err)
	}

	r.name = s

	r.Template, err = buildTemplate(s, string(b))
	if err != nil {
		return fmt.Errorf("template %q: %w", s, err)
	}

	return nil
}

// Type required to implement pflag.Value.
func (*TemplateValue) Type() string {
	return "text/template"
}

func buildTemplate(name, src string) (*template.Template, error) {
	tmplFuncs := template.FuncMap{
		"trimSuffix": strings.TrimSuffix,
		"baseName":   strings.BaseName,
		"replace":    strings.Replace,
		"split":      strings.Split,
		"join":       strings.Join,
	}

	tmpl, err := template.New(name).Funcs(tmplFuncs).Parse(src)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	return tmpl, nil
}
