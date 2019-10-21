package converter

import (
	"flag"
	"github.com/chrusty/protoc-gen-jsonschema/pbplugin"
	plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	allowNullValues = flag.Bool("allow_null_values", false,
		`Allow NULL values (by default, JSONSchemas will reject NULL values unless we explicitly allow them)`)
	disallowAdditionalProperties = flag.Bool("disallow_additional_properties", false,
		`Disallow additional properties (JSONSchemas won't validate JSON containing extra parameters)`)
	disallowBigIntsAsStrings = flag.Bool("disallow_big_ints_as_strings", false,
		`Disallow permissive validation of big-integers as strings (eg scientific notation)`)
	useProtoAndJSONFieldnames = flag.Bool("proto_and_json_fieldnames", false,
		`???`)
	debug = flag.Bool("debug", false,
		`Enable debug logging`)
)

func CodeGenerator(req *plugin_go.CodeGeneratorRequest) ([]*plugin_go.CodeGeneratorResponse_File, error) {
	pbplugin.ParseFlagsFromOptions(req.GetParameter())
	logger := logrus.New()
	logger.SetOutput(os.Stderr)
	logger.SetLevel(logrus.ErrorLevel)
	if *debug {
		logger.SetLevel(logrus.DebugLevel)
	}
	c := Converter{
		AllowNullValues:              *allowNullValues,
		DisallowAdditionalProperties: *disallowAdditionalProperties,
		DisallowBigIntsAsStrings:     *disallowBigIntsAsStrings,
		UseProtoAndJSONFieldnames:    *useProtoAndJSONFieldnames,
		logger:                       logger,
	}
	res, err := c.convert(req)
	return res.File, err
}
