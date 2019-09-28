package pbplugintest

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	ProtocBinary string
)

func init() {
	ProtocBinary, _ = exec.LookPath("protoc")
}

// Load the specified .proto files into a FileDescriptorSet. Any errors in loading/parsing will
// immediately fail the test. If no files are specified, all files will be read in.
func ReadFileDescriptorSet(t *testing.T, includePath string, filenames ...string) *descriptor.FileDescriptorSet {
	if ProtocBinary == "" {
		t.Fatalf("can't find 'protoc' binary (is it in your $PATH?)")
	}

	// If no files were specified, include all '.proto' files in the include path.
	if len(filenames) == 0 {
		matches, err := filepath.Glob(includePath + "/**/*.proto")
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range matches {
			filenames = append(filenames, match)
		}
	}

	// Use protoc to output descriptor info for the specified .proto files.
	var args []string
	args = append(args, "--descriptor_set_out=/dev/stdout")
	args = append(args, "--include_source_info")
	args = append(args, "--include_imports")
	args = append(args, "--proto_path="+includePath)
	args = append(args, filenames...)
	cmd := exec.Command(ProtocBinary, args...)
	stdoutBuf := bytes.Buffer{}
	stderrBuf := bytes.Buffer{}
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Run()
	if err != nil {
		t.Fatalf("failed to load descriptor set (%s): %s: %s",
			strings.Join(cmd.Args, " "), err.Error(), stderrBuf.String())
	}
	fds := &descriptor.FileDescriptorSet{}
	err = proto.Unmarshal(stdoutBuf.Bytes(), fds)
	if err != nil {
		t.Fatalf("failed to parse protoc output as FileDescriptorSet: %s", err.Error())
	}
	return fds
}
