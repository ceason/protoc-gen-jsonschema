package pbplugin

import (
	"flag"
	"io"
	"os"
	"strings"
	"text/template"
)

var OptionDelimiter = ","

func init() {
	flag.Usage = func() {
		printUsage(os.Stderr)
	}
}

func pluginName() string {
	parts := strings.Split(os.Args[0], "protoc-gen-")
	if len(parts) > 1 {
		return parts[1]
	}
	return "UNKNOWN_PLUGIN_NAME"
}

// Parse 'flag' package variables from plugin options.
func ParseFlagsFromOptions(s string) {
	// Panicking (rather than the default behavior of exiting) will allow
	// us to catch & forward the parsing error so it may be returned in the protoc response.
	flag.CommandLine.Init("protoc-gen-"+pluginName(), flag.ExitOnError)

	// The 'flag' package expects flags to begin with dashes.
	var args []string
	for _, o := range strings.Split(s, OptionDelimiter) {
		args = append(args, "--"+o)
	}

	// Reset all flag values to their defaults before parsing.
	flag.VisitAll(func(f *flag.Flag) {
		err := f.Value.Set(f.DefValue)
		if err != nil {
			panic(err)
		}
	})
	_ = flag.CommandLine.Parse(args)
	return
}

func printUsage(w io.Writer) {
	type opt struct {
		Name  string
		Usage string
	}
	var opts []opt
	flag.VisitAll(func(f *flag.Flag) {
		_, usageTxt := flag.UnquoteUsage(f)
		opts = append(opts, opt{f.Name, usageTxt})
	})
	usage := struct {
		PluginName      string
		Options         []opt
		OptionDelimiter string
	}{}
	usage.PluginName = pluginName()
	usage.Options = opts
	usage.OptionDelimiter = OptionDelimiter

	tmpl, err := template.New("usage").Parse(`This is a protoc plugin and should be invoked like:

   protoc --{{.PluginName}}_out=[OPTION[=VALUE]{{.OptionDelimiter}}OPTION...:]OUT_DIR [PROTOC_ARGS...]

This plugin has the following options (for PROTOC_ARGS see 'protoc --help'):
   {{range .Options}}
   {{.Name}}
      {{.Usage}}
   {{end}}
`)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(w, usage)
	if err != nil {
		panic(err)
	}
}
