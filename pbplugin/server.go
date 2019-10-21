package pbplugin

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type CodeGenerator func(req *plugin_go.CodeGeneratorRequest) (files []*plugin_go.CodeGeneratorResponse_File, err error)

type codeGeneratorError struct {
	error
}

// This should be used to indicate errors in .proto files which prevent the
// code generator from generating correct code.
func Errorf(format string, a ...interface{}) error {
	return codeGeneratorError{fmt.Errorf(format, a...)}
}

func Serve(fn CodeGenerator) {
	// Print usage if the plugin is not being invoked properly.
	stat, _ := os.Stdin.Stat()
	if !((stat.Mode() & os.ModeCharDevice) == 0) {
		printUsage(os.Stderr)
		os.Exit(2)
	}
	// Serve request & output any errors.
	err := fn.serve(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func (fn CodeGenerator) serve(r io.Reader, w io.Writer) error {
	// Read request from STDIN, pass it to the handler, then write response to STDOUT.
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("could not read codegen request from stdin: %s", err.Error())
	}
	req := &plugin_go.CodeGeneratorRequest{}
	err = proto.Unmarshal(input, req)
	if err != nil {
		return fmt.Errorf("could not read codegen request from stdin: %s", err.Error())
	}
	resError := ""
	files, err := fn(req)
	if err != nil {
		switch err.(type) {
		case codeGeneratorError:
			resError = err.Error()
			err = nil
		default:
			return err
		}
	}
	res := &plugin_go.CodeGeneratorResponse{
		Error: &resError,
		File:  files,
	}
	bytes, err := proto.Marshal(res)
	if err != nil {
		return fmt.Errorf("could not write serialized response: %s", err.Error())
	}
	_, err = w.Write(bytes)
	if err != nil {
		return fmt.Errorf("could not write serialized response: %s", err.Error())
	}
	return nil
}
