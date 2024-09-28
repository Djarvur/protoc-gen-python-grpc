package main

import (
	"github.com/Djarvur/protoc-gen-python-grpc/cmd/protoc-gen-python-grpc/internal/flags"
)

func main() {

	if err := flags.Root().Execute(); err != nil {
		panic(err)
	}

}
