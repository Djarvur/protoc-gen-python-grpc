package flags

import (
	"flag"
	"strings"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/flags/template"
)

type Flags struct {
	Template *template.TemplateValue
	Suffix   string
}

func Parse(in *string) Flags {
	out := Flags{
		Template: template.NewTemplateValue(),
		Suffix:   "_pb2_grpc.py",
	}

	if in == nil || *in == "" {
		return out
	}

	parser := flag.NewFlagSet("", flag.ExitOnError)

	parser.Var(out.Template, "template", "The template to use for generation.")
	parser.StringVar(&out.Suffix, "suffix", out.Suffix, "generated file(s) name suffix")

	parser.Parse(strings.Split(*in, ","))

	return out
}
