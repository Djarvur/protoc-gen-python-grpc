package kit

import (
	"io"

	"github.com/pseudomuto/protokit"
	"google.golang.org/protobuf/types/pluginpb"
)

type Plugin interface {
	Generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error)
}

type Kit struct {
}

func New() Kit {
	return Kit{}
}

func (k Kit) RunPluginWithIO(p Plugin, r io.Reader, w io.Writer) error {
	return protokit.RunPluginWithIO(p, r, w)
}
