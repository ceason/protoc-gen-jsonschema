package main

import (
	"github.com/chrusty/protoc-gen-jsonschema/internal/converter"
	"github.com/chrusty/protoc-gen-jsonschema/pbplugin"
)

func main() {
	pbplugin.Serve(converter.CodeGenerator)
}
