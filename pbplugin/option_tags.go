package pbplugin

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	OptionDelimiter = ","
	OptionTag       = "pbplugin"
)

var optionTypes = []string{
	"bool",
	"int",
	"string",
}

type option struct {
	Name        string
	Description string
	OptType     string
	fieldName   string
}

func unmarshalField(value *string, field reflect.Value) error {
	switch field.Kind() {
	case reflect.Int:
		if value == nil {
			return fmt.Errorf("no value provided")
		}
		parsed, err := strconv.Atoi(*value)
		if err != nil {
			return err
		}
		field.SetInt(int64(parsed))

	case reflect.String:
		if value == nil {
			return fmt.Errorf("no value provided")
		}
		field.SetString(*value)

	case reflect.Bool:
		// An empty value defaults to true.
		if value == nil {
			field.SetBool(true)
			return nil
		}
		// Explicitly provided values use strconv.
		parsed, err := strconv.ParseBool(*value)
		if err != nil {
			return err
		}
		field.SetBool(parsed)

	default:
		panic(fmt.Sprintf("Unimplemented kind '%s'", field.Kind()))
	}
	return nil
}

func unmarshalOptions(s string, v interface{}) error {
	// Split option name & value, defaulting to nil if no value is specified.
	optionMap := map[string]*string{}
	for _, opt := range strings.Split(s, OptionDelimiter) {
		if opt == "" {
			continue // Skip empty options
		}
		parts := strings.SplitN(opt, "=", 2)
		if len(parts) == 1 {
			optionMap[parts[0]] = nil
		} else {
			optionMap[parts[0]] = &parts[1]
		}
	}

	// Set struct field values based on parsed options.
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("unmarshaling plugin options requires a **pointer to a struct**, but received a %s instead", val.Kind())
	}
	val = val.Elem()
	options, err := getOptionsFromTags(v)
	if err != nil {
		return err
	}
UnmarshalOptionField:
	for _, o := range options {
		f := val.FieldByName(o.fieldName)
		optValue, isSet := optionMap[o.Name]

		// Options not explicitly provided are set to the zero value.
		if !isSet {
			f.Set(reflect.Zero(f.Type()))
			continue UnmarshalOptionField
		}
		err := unmarshalField(optValue, f)
		if err != nil {
			return fmt.Errorf("failed to parse option '%s': %s", o.Name, err.Error())
		}
	}

	// Check if any nonexistent options were provided.
	for _, o := range options {
		delete(optionMap, o.Name)
	}
	for o, _ := range optionMap {
		return fmt.Errorf("no such option '%s' exists", o)
	}
	return nil
}

func getOptionsFromTags(val interface{}) ([]option, error) {
	var opts []option
	t := reflect.TypeOf(val)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		description, isOption := field.Tag.Lookup(OptionTag)

		// Skip this field if it's not tagged as an option.
		if !isOption {
			continue
		}

		// Check that this field is a supported type.
		if !isSupportedType(field.Type) {
			return nil, fmt.Errorf(
				"Option field '%s' has unsupported type '%s' (supported types: %s)",
				field.Name, field.Type.Name(), strings.Join(optionTypes, ", "))
		}

		// Add the field to our output.
		opts = append(opts, option{
			Name:        toSnakeCase(field.Name),
			Description: description,
			OptType:     field.Type.Name(),
			fieldName:   field.Name,
		})
	}
	return opts, nil
}

func isSupportedType(t reflect.Type) bool {
	for _, optionType := range optionTypes {
		if t.Name() == optionType {
			return true
		}
	}
	return false
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
