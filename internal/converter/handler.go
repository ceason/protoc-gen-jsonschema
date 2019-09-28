package converter

import plugin_go "github.com/golang/protobuf/protoc-gen-go/plugin"

func (c *Converter) Handle(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	return c.convert(req)
}
