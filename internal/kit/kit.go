package kit

import (
	"errors"
	"io"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	ErrGenerate         = errors.New("generate error")
	ErrNoFileToGenerate = errors.New("no files were supplied to the plugin")
	ErrReadRequest      = errors.New("read request error")
	ErrRun              = errors.New("generator run error")
	ErrWriteResponse    = errors.New("write response error")
)

type Generator interface {
	Generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error)
}

type Kit struct{}

func New() Kit {
	return Kit{}
}

func (k Kit) RunPluginWithIO(generator Generator, r io.Reader, writer io.Writer) error {
	req, err := k.readRequest(r)
	if err != nil {
		return errors.Join(ErrRun, err)
	}

	resp, err := generator.Generate(req)
	if err != nil {
		return errors.Join(ErrRun, ErrGenerate, err)
	}

	err = k.writeResponse(writer, resp)
	if err != nil {
		return errors.Join(ErrRun, err)
	}

	return nil
}

func (k Kit) readRequest(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Join(ErrReadRequest, err)
	}

	req := new(pluginpb.CodeGeneratorRequest)
	if err = proto.Unmarshal(data, req); err != nil {
		return nil, errors.Join(ErrReadRequest, err)
	}

	if len(req.GetFileToGenerate()) == 0 {
		return nil, errors.Join(ErrReadRequest, ErrNoFileToGenerate)
	}

	return req, nil
}

func (k Kit) writeResponse(writer io.Writer, resp *pluginpb.CodeGeneratorResponse) error {
	data, err := proto.Marshal(resp)
	if err != nil {
		return errors.Join(ErrWriteResponse, err)
	}

	if _, err := writer.Write(data); err != nil {
		return errors.Join(ErrWriteResponse, err)
	}

	return nil
}
