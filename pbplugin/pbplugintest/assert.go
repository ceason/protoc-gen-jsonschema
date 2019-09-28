package pbplugintest

import (
	"github.com/chrusty/protoc-gen-jsonschema/pbplugin"
	"testing"
)

// Assert that `h`, given all of the .proto files in `path`, will
// output all of the non-proto files in `path`.
func AssertCodegenOutput(t *testing.T, h pbplugin.Handler, path string) bool {
	panic("Unimplemented")
}
