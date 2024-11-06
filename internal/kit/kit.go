package kit

import (
	"errors"
	"io"

	plugingo "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"
)

var ErrRun = errors.New("generator run error")

type Plugin interface {
	Generate(req *plugingo.CodeGeneratorRequest) (*plugingo.CodeGeneratorResponse, error)
}

type Kit struct{}

func New() Kit {
	return Kit{}
}

func (k Kit) RunPluginWithIO(p Plugin, r io.Reader, w io.Writer) error {
	err := protokit.RunPluginWithIO(p, r, w)
	if err != nil {
		return errors.Join(ErrRun, err)
	}

	return nil
}
