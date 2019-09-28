package pbplugin

import (
	"reflect"
	"testing"
)

func TestUnmarshalOptions(t *testing.T) {
	type asdf struct {
		SomeBooleanOption bool   `pbplugin:"An example boolean option"`
		AStringOption     string `pbplugin:"An example string option"`
		AnIntOption       int    `pbplugin:"An example integer option"`
	}
	optionStr := "some_boolean_option,a_string_option=asdf,an_int_option=1234"
	expected := asdf{
		SomeBooleanOption: true,
		AStringOption:     "asdf",
		AnIntOption:       1234,
	}
	ts := asdf{}
	err := unmarshalOptions(optionStr, &ts)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expected, ts) {
		t.Errorf(`Did not unmarshal expected option values
  want: %v
   got: %v
`, expected, ts)
	}
}

func TestGetOptionsFromTags(t *testing.T) {
	type blah struct {
		APluginOption bool `pbplugin:"Description of this option"`
		AnotherOption bool `pbplugin:"Description for another option"`
	}
	expected := []option{
		{"a_plugin_option", "Description of this option", "bool", "APluginOption"},
		{"another_option", "Description for another option", "bool", "AnotherOption"},
	}
	ts := &blah{}
	optionMap := map[string]option{}
	options, err := getOptionsFromTags(ts)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range options {
		optionMap[o.Name] = o
	}
	for _, want := range expected {
		got, ok := optionMap[want.Name]
		if !ok {
			t.Errorf("Missing expected option '%s'", want.Name)
			continue
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf(`Did not get expected option info from tag
  want: %v
   got: %v
`, want, got)
		}
	}

}
