package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/pseudomuto/protokit"
	"google.golang.org/protobuf/proto"

	pluginpb "github.com/golang/protobuf/protoc-gen-go/plugin"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/strings"
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
var SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

// Generator describes a protoc code generate plugin. It's an implementation of Generator from github.com/pseudomuto/protokit
type Generator struct {
	Suffix   string
	Template string
}

// Generate compiles the documentation and generates the CodeGeneratorResponse to send back to protoc. It does this
// by rendering a template based on the options parsed from the CodeGeneratorRequest.
func (p *Generator) Generate(r *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	resp := new(pluginpb.CodeGeneratorResponse)

	for _, fds := range protokit.ParseCodeGenRequest(r) {
		data := ProtoFile{
			Package:  fds.GetPackage(),
			Name:     fds.GetName(),
			Services: buildServices(fds.GetServices()),
		}

		f, errExecute := executeTemplate(p.Template, data)
		if errExecute != nil {
			return nil, errExecute
		}

		resp.File = append(
			resp.File,
			&pluginpb.CodeGeneratorResponse_File{
				Name:    proto.String(strings.Replace("-", "_", strings.TrimSuffix(".", data.Name)+p.Suffix)),
				Content: proto.String(f),
			},
		)
	}

	resp.SupportedFeatures = proto.Uint64(SupportedFeatures)

	return resp, nil
}

func executeTemplate(tmplSrc string, data interface{}) (string, error) {
	var tmplFuncs = template.FuncMap{
		"trimSuffix": strings.TrimSuffix,
		"baseName":   strings.BaseName,
		"replace":    strings.Replace,
		"split":      strings.Split,
		"join":       strings.Join,
	}

	tmpl, err := template.New("").Funcs(tmplFuncs).Parse(tmplSrc)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	buf := new(bytes.Buffer)

	if err = tmpl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
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

	for _, m := range in {
		out = append(
			out, Method{
				Name:            m.GetName(),
				Comment:         m.GetComments().String(),
				Request:         m.GetInputType(),
				Response:        m.GetOutputType(),
				ClientStreaming: m.GetClientStreaming(),
				ServerStreaming: m.GetServerStreaming(),
			},
		)
	}

	return out
}
