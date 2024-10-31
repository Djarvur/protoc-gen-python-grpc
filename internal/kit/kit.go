package kit

import (
	"io"

	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"
)

type Plugin interface {
	Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error)
}

type Kit struct {
}

func New() Kit {
	return Kit{}
}

func (k Kit) RunPluginWithIO(p Plugin, r io.Reader, w io.Writer) error {
	return protokit.RunPluginWithIO(p, r, w)
}
