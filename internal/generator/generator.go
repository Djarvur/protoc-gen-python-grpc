package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
	gostrings "strings"
	"text/template"

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

const (
	// serviceMethodFieldNumber contains method field number in ServiceDescriptorProto message
	serviceMethodFieldNumber = 2
	// serviceFieldNumber contains service field number in FileDescriptorProto message
	serviceFieldNumber = 6
)

// SupportedFeatures describes a flag setting for supported features.
const SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

// Generator describes a protoc code generate plugin.
type Generator struct {
	Suffix   string
	Template *template.Template
}

func New(suffix, tmplSrc string) (*Generator, error) {
	tmpl, err := buildTemplate(tmplSrc)
	if err != nil {
		return nil, err
	}

	return &Generator{
		Suffix:   suffix,
		Template: tmpl,
	}, nil
}

// Generate compiles the documentation and generates the CodeGeneratorResponse to send back to protoc. It does this
// by rendering a template based on the options parsed from the CodeGeneratorRequest.
func (p *Generator) Generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	resp := new(pluginpb.CodeGeneratorResponse)

	for _, fds := range req.GetProtoFile() {
		data := ProtoFile{
			Package:  fds.GetPackage(),
			Name:     fds.GetName(),
			Services: buildServices(fds),
		}

		content, errExecute := executeTemplate(p.Template, data)
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
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	return tmpl, nil
}

func executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return buf.String(), nil
}

func buildServices(fd *descriptorpb.FileDescriptorProto) []Service {
	out := make([]Service, 0, len(fd.GetService()))

	comments := parseServicesComments(fd)

	for serviceIdx, svc := range fd.GetService() {

		servicePath := fmt.Sprintf("%d.%d", serviceFieldNumber, serviceIdx)
		service := Service{
			Name:    svc.GetName(),
			Comment: comments[servicePath],
			Methods: []Method{},
		}

		for methodIdx, method := range svc.GetMethod() {
			methodPath := fmt.Sprintf("%s.%d.%d", servicePath, serviceMethodFieldNumber, methodIdx)
			service.Methods = append(
				service.Methods, Method{
					Name:            method.GetName(),
					Comment:         comments[methodPath],
					Request:         method.GetInputType(),
					Response:        method.GetOutputType(),
					ClientStreaming: method.GetClientStreaming(),
					ServerStreaming: method.GetServerStreaming(),
				},
			)
		}

		out = append(out, service)
	}

	return out
}

func parseServicesComments(fd *descriptorpb.FileDescriptorProto) map[string]string {
	comments := make(map[string]string)

	for _, loc := range fd.GetSourceCodeInfo().GetLocation() {
		if loc.GetLeadingComments() == "" && loc.GetTrailingComments() == "" {
			continue
		}

		path := loc.GetPath()

		if len(path) < 2 || path[0] != serviceFieldNumber {
			continue
		}

		b := new(bytes.Buffer)
		leading := scrub(loc.GetLeadingComments())
		if leading != "" {
			b.WriteString(leading)
			b.WriteString("\n\n")
		}

		b.WriteString(scrub(loc.GetTrailingComments()))

		commentPath := gostrings.ReplaceAll(gostrings.Trim(fmt.Sprintf("%v", path), "[]"), " ", ".")
		comments[commentPath] = gostrings.TrimSpace(b.String())
	}

	return comments
}

func scrub(str string) string {
	return gostrings.TrimSpace(gostrings.Replace(str, "\n ", "\n", -1))
}
