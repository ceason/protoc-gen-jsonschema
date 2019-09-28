package main

import (
	"github.com/chrusty/protoc-gen-jsonschema/internal/converter"
	"github.com/chrusty/protoc-gen-jsonschema/pbplugin"
	"log"
)

func main() {
	err := pbplugin.Serve(&converter.Converter{})
	if err != nil {
		log.Fatal(err)
	}
}
