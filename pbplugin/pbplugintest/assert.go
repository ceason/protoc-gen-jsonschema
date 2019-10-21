package pbplugintest

import (
	"github.com/chrusty/protoc-gen-jsonschema/pbplugin"
	"testing"
)

// Assert that `fn`, given all of the .proto files in `path`, will
// output all of the non-proto files in `path`.
func AssertCodegenOutput(t *testing.T, fn pbplugin.CodeGenerator, path string, pluginOpts ...string) {
	panic("Unimplemented")
}
