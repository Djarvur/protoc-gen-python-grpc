package plugin

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
	"io"
)

// CodeGenerator describes an interface for generating code based on an incoming request
type CodeGenerator interface {
	Generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error)
}

type Plugin struct {
	r         io.Reader
	w         io.Writer
	g         CodeGenerator
	initialed bool
}

func New(g CodeGenerator, r io.Reader, w io.Writer) (*Plugin, error) {
	if g == nil {
		return nil, fmt.Errorf("generator must not be nil")
	}
	if r == nil {
		return nil, fmt.Errorf("reader must not be nil")
	}
	if w == nil {
		return nil, fmt.Errorf("writer must not be nil")
	}

	return &Plugin{
		r:         r,
		w:         w,
		g:         g,
		initialed: true,
	}, nil
}

func (p *Plugin) Run() error {
	if !p.initialed {
		return fmt.Errorf("use New to setup plugin")
	}

	req, err := readRequest(p.r)
	if err != nil {
		return err
	}

	resp, err := p.g.Generate(req)
	if err != nil {
		return nil
	}

	return writeResponse(p.w, resp)
}

func readRequest(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	req := new(pluginpb.CodeGeneratorRequest)
	if err = proto.Unmarshal(data, req); err != nil {
		return nil, err
	}

	if len(req.GetFileToGenerate()) == 0 {
		return nil, fmt.Errorf("no files were supplied to the plugin")
	}

	return req, nil
}

func writeResponse(w io.Writer, resp *pluginpb.CodeGeneratorResponse) error {
	data, err := proto.Marshal(resp)
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}
