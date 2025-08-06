package generator

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/Djarvur/protokit"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/flags"
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
const SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL |
	pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)

var _ protokit.Plugin = (*generator)(nil)

// generator describes a protoc code generate plugin.
// It's an implementation of generator from github.com/Djarvur/protokit.
type generator struct{}

func New() *generator {
	return &generator{}
}

// Generate compiles the code and generates the CodeGeneratorResponse to send back to protoc. It does this
// by rendering a template based on the options parsed from the CodeGeneratorRequest.
func (p *generator) Generate(r *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	params := flags.Parse(r.Parameter)

	resp := new(pluginpb.CodeGeneratorResponse)

	for _, fds := range protokit.ParseCodeGenRequest(r) {
		data := ProtoFile{
			Package:  fds.GetPackage(),
			Name:     fds.GetName(),
			Services: buildServices(fds.GetServices()),
		}

		content, errExecute := executeTemplate(params.Template.Template, data)
		if errExecute != nil {
			return nil, errExecute
		}

		resp.File = append(
			resp.File,
			&pluginpb.CodeGeneratorResponse_File{ //nolint:exhaustruct
				Name: proto.String(
					strings.ReplaceAll(
						strings.TrimSuffix(data.Name, ".proto")+params.Suffix,
						"-",
						"_",
					),
				),
				Content: proto.String(content),
			},
		)
	}

	resp.SupportedFeatures = proto.Uint64(SupportedFeatures)
	resp.MinimumEdition = proto.Int32(int32(descriptorpb.Edition_EDITION_PROTO2))
	resp.MaximumEdition = proto.Int32(int32(descriptorpb.Edition_EDITION_2024))

	return resp, nil
}

func executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
	buf := new(bytes.Buffer)

	if err := tmpl.Execute(buf, data); err != nil {
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
