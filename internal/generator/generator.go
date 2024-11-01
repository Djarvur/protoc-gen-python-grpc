package generator

import (
	"bytes"
	"errors"
	"fmt"
	gostrings "strings"
	"text/template"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

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
const SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

const (
	// serviceMethodFieldNumber contains method field number in ServiceDescriptorProto message.
	serviceMethodFieldNumber = 2
	// serviceFieldNumber contains service field number in FileDescriptorProto message.
	serviceFieldNumber = 6
)

var (
	ErrTemplateParse = errors.New("template parsing error")
	ErrTemplateBuild = errors.New("template building error")
	ErrTemplateExec  = errors.New("template executing error")
)

// generator describes a protoc code generate plugin.
// It's an implementation of generator from github.com/pseudomuto/protokit.
type generator struct {
	Suffix   string
	Template string
}

func New(suffix, tmplSrc string) *generator {
	return &generator{
		Suffix:   suffix,
		Template: tmplSrc,
	}
}

// Generate compiles the documentation and generates the CodeGeneratorResponse to send back to protoc. It does this
// by rendering a template based on the options parsed from the CodeGeneratorRequest.
func (p *generator) Generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	tmpl, err := buildTemplate(p.Template)
	if err != nil {
		return nil, errors.Join(ErrTemplateBuild, err)
	}

	resp := new(pluginpb.CodeGeneratorResponse)

	for _, fds := range req.GetProtoFile() {
		data := ProtoFile{
			Package:  fds.GetPackage(),
			Name:     fds.GetName(),
			Services: buildServices(fds),
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
		return nil, errors.Join(ErrTemplateParse, err)
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

		buf := new(bytes.Buffer)

		if leading := scrub(loc.GetLeadingComments()); leading != "" {
			buf.WriteString(leading)
			buf.WriteString("\n\n")
		}

		buf.WriteString(scrub(loc.GetTrailingComments()))

		commentPath := gostrings.ReplaceAll(gostrings.Trim(fmt.Sprintf("%v", path), "[]"), " ", ".")
		comments[commentPath] = gostrings.TrimSpace(buf.String())
	}

	return comments
}

func scrub(str string) string {
	return gostrings.TrimSpace(gostrings.ReplaceAll(str, "\n ", "\n"))
}
