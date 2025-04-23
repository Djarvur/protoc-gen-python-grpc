package main

import (
	"os"

	"github.com/Djarvur/protoc-gen-python-grpc/internal/generator"
	"github.com/Djarvur/protoc-gen-python-grpc/internal/kit"
)

func main() {
	if err := kit.New().RunPluginWithIO(generator.New(), os.Stdin, os.Stdout); err != nil {
		panic(err)
	}
}
