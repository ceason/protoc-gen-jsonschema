package pbplugin

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Handler interface {
	Handle(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error)
}

func Serve(handler Handler) error {
	// Print usage if the plugin is not being invoked properly.
	stat, _ := os.Stdin.Stat()
	if !((stat.Mode() & os.ModeCharDevice) == 0) {
		return printUsage(os.Stderr, handler)
	}

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
	err = unmarshalOptions(req.GetParameter(), handler)
	if err != nil {
		// todo: figure out which errors should go to codegen response, vs which should be returned to caller??
		return err
	}
	res, err := handler.Handle(req)
	if err != nil {
		errMsg := err.Error()
		res = &plugin_go.CodeGeneratorResponse{
			Error: &errMsg,
		}
	}
	bytes, err := proto.Marshal(res)
	if err != nil {
		return fmt.Errorf("could not write serialized response: %s", err.Error())
	}
	_, err = os.Stdout.Write(bytes)
	if err != nil {
		return fmt.Errorf("could not write serialized response: %s", err.Error())
	}
	return nil
}

func printUsage(w io.Writer, h Handler) error {
	options, err := getOptionsFromTags(h)
	if err != nil {
		return err
	}
	usage := struct {
		PluginName      string
		Options         []option
		OptionDelimiter string
	}{}
	usage.PluginName = strings.Split(os.Args[0], "protoc-gen-")[1]
	usage.Options = options
	usage.OptionDelimiter = OptionDelimiter

	tmpl, err := template.New("usage").Parse(`This is a protoc plugin and should be invoked like:

   protoc --{{.PluginName}}_out=[OPTION{{.OptionDelimiter}}OPTION...:]OUT_DIR [PROTOC_ARGS...]

This plugin has the following options (for PROTOC_ARGS see 'protoc --help'):
   {{range .Options}}
     {{.Name}}        ({{.OptType}}) {{.Description}}
   {{end}}
`)
	if err != nil {
		return err
	}
	err = tmpl.Execute(w, usage)
	if err != nil {
		return err
	}
	return nil
}
