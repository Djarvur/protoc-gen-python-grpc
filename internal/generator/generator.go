package generator

import (
	"bytes"
	_ "embed"
	"errors"
	"text/template"

	pluginpb "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"
	"google.golang.org/protobuf/proto"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/strings"
)

var (
	ErrTemplateBuild = errors.New("template building error")
	ErrTemplateExec  = errors.New("template executing error")
)

type ProtoFile struct {
	Name     string
	Package  string
	Services []Service
}

type Service struct {
	Name    string
	Comment string
	Methods []Method
}

type Method struct {
	Name            string
	Comment         string
	Request         string
	Response        string
	ClientStreaming bool
	ServerStreaming bool
}

// SupportedFeatures describes a flag setting for supported features.
const SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

var _ protokit.Plugin = (*Generator)(nil)

// Generator describes a protoc code generate plugin.
// It's an implementation of Generator from github.com/pseudomuto/protokit.
type Generator struct {
	Suffix   string
	Template string
}

func New(suffix, tmplSrc string) *Generator {
	return &Generator{
		Suffix:   suffix,
		Template: tmplSrc,
	}
}

// Generate compiles the documentation and generates the CodeGeneratorResponse to send back to protoc. It does this
// by rendering a template based on the options parsed from the CodeGeneratorRequest.
func (p *Generator) Generate(request *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	tmpl, err := buildTemplate(p.Template)
	if err != nil {
		return nil, err
	}

	resp := new(pluginpb.CodeGeneratorResponse)

	for _, fds := range protokit.ParseCodeGenRequest(request) {
		data := ProtoFile{
			Package:  fds.GetPackage(),
			Name:     fds.GetName(),
			Services: buildServices(fds.GetServices()),
		}

		content, errExecute := executeTemplate(tmpl, data)
		if errExecute != nil {
			return nil, errExecute
		}

		resp.File = append(
			resp.File,
			&pluginpb.CodeGeneratorResponse_File{ //nolint:exhaustruct
				Name:    proto.String(strings.Replace("-", "_", strings.TrimSuffix(".", data.Name)+p.Suffix)),
				Content: proto.String(content),
			},
		)
	}

	resp.SupportedFeatures = proto.Uint64(SupportedFeatures)

	return resp, nil
}

func buildTemplate(tmplSrc string) (*template.Template, error) {
	tmplFuncs := template.FuncMap{
		"trimSuffix": strings.TrimSuffix,
		"baseName":   strings.BaseName,
		"replace":    strings.Replace,
		"split":      strings.Split,
		"join":       strings.Join,
	}

	tmpl, err := template.New("").Funcs(tmplFuncs).Parse(tmplSrc)
	if err != nil {
		return nil, errors.Join(ErrTemplateBuild, err)
	}

	return tmpl, nil
}

func executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, data); err != nil {
		return "", errors.Join(ErrTemplateExec, err)
	}

	return buf.String(), nil
}

func buildServices(in []*protokit.ServiceDescriptor) []Service {
	out := make([]Service, 0, len(in))

	for _, svc := range in {
		out = append(
			out,
			Service{
				Name:    svc.GetName(),
				Comment: svc.GetComments().String(),
				Methods: buildMethods(svc.GetMethods()),
			},
		)
	}

	return out
}

func buildMethods(in []*protokit.MethodDescriptor) []Method {
	out := make([]Method, 0, 1)

	for _, method := range in {
		out = append(
			out, Method{
				Name:            method.GetName(),
				Comment:         method.GetComments().String(),
				Request:         method.GetInputType(),
				Response:        method.GetOutputType(),
				ClientStreaming: method.GetClientStreaming(),
				ServerStreaming: method.GetServerStreaming(),
			},
		)
	}

	return out
}
