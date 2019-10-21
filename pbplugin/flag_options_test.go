package pbplugin

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseOptions(t *testing.T) {
	someBooleanOption := flag.Bool("some_boolean_option", false, "An example boolean option")
	anotherBooleanOption := flag.Bool("another_boolean_option", false, "Another example boolean option")
	aStringOption := flag.String("a_string_option", "", "An example string option")
	anIntOption := flag.Int("an_int_option", 0, "An example integer option")

	ParseFlagsFromOptions("some_boolean_option,a_string_option=asdf,an_int_option=1234")
	assert.Equal(t, true, *someBooleanOption)
	assert.Equal(t, false, *anotherBooleanOption)
	assert.Equal(t, "asdf", *aStringOption)
	assert.Equal(t, 1234, *anIntOption)

	// Make sure options are reset to defaults before reparsing.
	ParseFlagsFromOptions("another_boolean_option=true,a_string_option=fdsa")
	assert.Equal(t, false, *someBooleanOption)
	assert.Equal(t, true, *anotherBooleanOption)
	assert.Equal(t, "fdsa", *aStringOption)
	assert.Equal(t, 0, *anIntOption)
}
