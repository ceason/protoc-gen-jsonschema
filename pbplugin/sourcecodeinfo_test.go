package pbplugin

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"os/exec"
	"strings"
	"testing"
)

func TestSourceInfoLookup(t *testing.T) {
	// Read in the test file & get references to the things we've declared.
	// Note that the hardcoded indexes must reflect the declaration order in
	// the .proto file.
	fds := mustReadProtoFiles(t, "UNIMPLEMENTED", "MessageWithComments.proto")
	protoFile := fds.File[0]
	msgWithComments := protoFile.MessageType[0]
	msgWithComments_name1 := msgWithComments.Field[0]

	// Create an instance of our thing and test that it returns the expected
	// source data for each of our above declarations.
	src := NewSourceInfo(fds.File)
	assertCommentsMatch(t, src.GetMessage(msgWithComments), &descriptor.SourceCodeInfo_Location{
		LeadingComments: strPtr(" This is a message level comment and talks about what this message is and why you should care about it!\n"),
	})
	assertCommentsMatch(t, src.GetField(msgWithComments_name1), &descriptor.SourceCodeInfo_Location{
		LeadingComments: strPtr(" This field is supposed to represent blahblahblah\n"),
	})
}

func assertCommentsMatch(t *testing.T, actual, expected *descriptor.SourceCodeInfo_Location) {
	if len(actual.LeadingDetachedComments) != len(expected.LeadingDetachedComments) {
		t.Fatalf("Wrong value for LeadingDetachedComments.\n   got: %v\n   want: %v",
			actual.LeadingDetachedComments, expected.LeadingDetachedComments)
	}
	for i := 0; i < len(actual.LeadingDetachedComments); i++ {
		if actual.LeadingDetachedComments[i] != expected.LeadingDetachedComments[i] {
			t.Fatalf("Wrong value for LeadingDetachedComments.\n   got: %v\n   want: %v",
				actual.LeadingDetachedComments, expected.LeadingDetachedComments)
		}
	}
	if actual.GetTrailingComments() != expected.GetTrailingComments() {
		t.Fatalf("Wrong value for TrailingComments.\n   got: %q\n   want: %q",
			actual.GetTrailingComments(), expected.GetTrailingComments())
	}
	if actual.GetLeadingComments() != expected.GetLeadingComments() {
		t.Fatalf("Wrong value for LeadingComments.\n   got: %q\n   want: %q",
			actual.GetLeadingComments(), expected.GetLeadingComments())
	}
}

// Go doesn't have syntax for addressing a string literal, so this is the next best thing.
func strPtr(s string) *string {
	return &s
}

// Load the specified .proto files into a FileDescriptorSet. Any errors in loading/parsing will
// immediately fail the test.
func mustReadProtoFiles(t *testing.T, includePath string, filenames ...string) *descriptor.FileDescriptorSet {
	protocBinary, err := exec.LookPath("protoc")
	if err != nil {
		t.Fatalf("Can't find 'protoc' binary in $PATH: %s", err.Error())
	}

	// Use protoc to output descriptor info for the specified .proto files.
	var args []string
	args = append(args, "--descriptor_set_out=/dev/stdout")
	args = append(args, "--include_source_info")
	args = append(args, "--include_imports")
	args = append(args, "--proto_path="+includePath)
	args = append(args, filenames...)
	cmd := exec.Command(protocBinary, args...)
	stdoutBuf := bytes.Buffer{}
	stderrBuf := bytes.Buffer{}
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
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
