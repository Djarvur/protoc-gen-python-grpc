package flags

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/generator"
	"github.com/Djarvur/protoc-gen-python-grpc/internal/kit"
	"github.com/Djarvur/protoc-gen-python-grpc/internal/template"
)

const (
	templateFlag          = "template"
	fileNameSuffixFlag    = "suffix"
	fileNameSuffixDefault = "_pb2_grpc.py"
)

func Root() *cobra.Command {
	templateSource := template.DefaultValue()

	cmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "protoc-gen-python-grpc",
		Short: "protoc plugin to generate Python gRPC code",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runRoot(
				must(cmd.Flags().GetString(fileNameSuffixFlag)),
				templateSource.Source(),
			)
		},
	}

	cmd.Flags().String(fileNameSuffixFlag, fileNameSuffixDefault, "generated file(s) name suffix")
	cmd.Flags().Var(templateSource, templateFlag, "template to be used")

	return cmd
}

func runRoot(suffix, templateSource string) {
	if err := kit.New().RunPluginWithIO(generator.New(suffix, templateSource), os.Stdin, os.Stdout); err != nil {
		panic(err)
	}
}
