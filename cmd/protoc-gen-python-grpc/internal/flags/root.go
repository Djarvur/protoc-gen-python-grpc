package flags

import (
	"github.com/spf13/cobra"
	"os"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/generator"
	"github.com/Djarvur/protoc-gen-python-grpc/internal/plugin"
)

const (
	templateFlag          = "template"
	fileNameSuffixFlag    = "suffix"
	fileNameSuffixDefault = "_pb2_grpc.py"
)

func Root() *cobra.Command {
	templateSource := newTemplateValue()

	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "protoc-gen-python-grpc",
		Short: "protoc plugin to generate Python gRPC code",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runRoot(
				must(cmd.Flags().GetString(fileNameSuffixFlag)),
				templateSource.source,
			)
		},
	}

	cmd.Flags().String(fileNameSuffixFlag, fileNameSuffixDefault, "generated file(s) name suffix")
	cmd.Flags().Var(templateSource, templateFlag, "template to be used")

	return cmd
}

func runRoot(suffix, templateSource string) {
	p, err := plugin.New(must(generator.New(suffix, templateSource)), os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}

	err = p.Run()
	if err != nil {
		panic(err)
	}
}
