package main

import (
	"os"

	"github.com/pseudomuto/protokit"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/generator"
)

func main() {
	if err := protokit.RunPluginWithIO(generator.New(), os.Stdin, os.Stdout); err != nil {
		panic(err)
	}
}
