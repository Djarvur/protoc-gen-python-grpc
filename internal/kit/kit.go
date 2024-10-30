package kit

import (
	"bytes"
	"io"
	"os"

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
	in, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile("testdata/in.bin", in, 0o644); err != nil {
		panic(err)
	}

	var out bytes.Buffer

	errPlugin := protokit.RunPluginWithIO(p, bytes.NewBuffer(in), &out)

	if err = os.WriteFile("testdata/out.bin", out.Bytes(), 0o644); err != nil {
		panic(err)
	}

	if _, err = io.Copy(w, &out); err != nil {
		panic(err)
	}

	return errPlugin
}
