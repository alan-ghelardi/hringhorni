package config

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/ghodss/yaml"
	"github.com/gobuffalo/flect"
	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/configmap"
)

// Parser is a function that takes a string read from a value in a config map
// and attempts to turn it into a higher object when applicable. The function
// returns an error with a meaningful message when the value is considered
// invalid.
type Parser func(value string) (interface{}, error)

// Helper functions to parse several values.

func ParseBool(value string) (interface{}, error) {
	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	}
	return false, fmt.Errorf("expected true or false, but got %q", value)
}

func ParseEnum(values ...string) Parser {
	return func(value string) (interface{}, error) {
		if sets.NewString(values...).Has(value) {
			return value, nil
		}
		return "", fmt.Errorf("expected one of %s, but got %q", strings.Join(values, ", "), value)
	}
}

func ParseInt(value string) (interface{}, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", value)
	}
	return i, nil
}

func ParseString(value string) (interface{}, error) {
	return value, nil
}

func ParseImageRef(value string) (interface{}, error) {
	if strings.HasPrefix(value, "ko://") {
		return nil, fmt.Errorf("invalid image ref %s", value)
	}
	return value, nil
}

func ParseURL(value string) (interface{}, error) {
	if _, err := url.ParseRequestURI(value); err != nil {
		return "", fmt.Errorf("invalid URL \"%s\"", value)
	}
	return value, nil
}

func ParseYAMLMap(value string) (interface{}, error) {
	var m map[string]string
	if err := yaml.Unmarshal([]byte(value), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func ParseYAMLSlice(value string) (interface{}, error) {
	var s []string
	if err := yaml.Unmarshal([]byte(value), &s); err != nil {
		return nil, err
	}
	return s, nil
}

func ParseSetOfStrings(value string) (interface{}, error) {
	elements, err := ParseYAMLSlice(value)
	if err != nil {
		return nil, err
	}
	return sets.NewString(elements.([]string)...), nil
}

// Unmarshal is a helper function to fill out a given object with values parsed
// from the supplied config map.
//
// Provide a map of string to parser func where the key is the respective field
// name in kebab case. Thus, if the key webhook-url associated to the func
// parseURL is given, the Unmarshal func will attempt to parse whatever value is
// assigned to this key in the config map and assign the resulting value to the
// field WebhookUrl in the supplied object.
//
// This function returns the aggregated errors encountered during the parsing process.
func Unmarshal(object interface{}, configMap *corev1.ConfigMap, parsers map[string]Parser) (errors *apis.FieldError) {
	for key, value := range configMap.Data {
		if parser, exists := parsers[key]; exists {
			if x, err := parser(value); err != nil {
				errors = errors.Also(apis.ErrGeneric(err.Error(), "data."+key))
			} else {
				fieldName := flect.Pascalize(key)
				field := reflect.ValueOf(object).Elem().FieldByName(fieldName)
				if !field.IsValid() {
					panic(fmt.Errorf("field %q not found in the provided struct", fieldName))
				}
				field.Set(reflect.ValueOf(x))
			}
		} else if key != configmap.ExampleKey {
			errors = errors.Also(apis.ErrDisallowedFields("data." + key))
		}
	}
	return errors
}
