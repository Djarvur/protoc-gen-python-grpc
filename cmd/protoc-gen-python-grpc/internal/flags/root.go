package flags

import (
	"github.com/pseudomuto/protokit"
	"github.com/spf13/cobra"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/generator"
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
	if err := protokit.RunPlugin(must(generator.New(suffix, templateSource))); err != nil {
		panic(err)
	}
}
