package dotnotation

import (
	"fmt"
	"strconv"
	"strings"
)

// DefaultGetter returns the property value of a given target, or an error, supporting types like encoding/json.
// Supports one level of pointer indirection.
func DefaultGetter(target interface{}, property string) (interface{}, error) {
	// handle each type that is supported by simple unmarshalling of a json value
	// https://golang.org/pkg/encoding/json/#Unmarshal
	if v, ok := target.(*[]interface{}); ok {
		target = *v
	} else if v, ok := target.(*map[string]interface{}); ok {
		target = *v
	}
	switch v := target.(type) {
	case []interface{}:
		i, err := strconv.Atoi(property)

		if err != nil {
			return nil, fmt.Errorf("cannot get non-integer property '%s' on a slice", property)
		}

		if i < 0 || i >= len(v) {
			return nil, fmt.Errorf("cannot get out of range property '%s' on a slice", property)
		}

		return v[i], nil

	case map[string]interface{}:
		value, ok := v[property]

		if !ok {
			return nil, fmt.Errorf("cannot get non-existent property '%s' on a map", property)
		}

		return value, nil

	default:
		return nil, fmt.Errorf("cannot get property '%s' on type %T", property, target)
	}
}

// DefaultSetter sets the property value of a given target, to a given value, or returns an error, supporting types
// like encoding/json.
// Supports one level of pointer indirection, and appending to slices if a pointer is used.
func DefaultSetter(target interface{}, property string, value interface{}) error {
	// handle each type that is supported by simple unmarshalling of a json value
	// https://golang.org/pkg/encoding/json/#Unmarshal
	switch v := target.(type) {
	case []interface{}:
		i, err := strconv.Atoi(property)

		if err != nil {
			return fmt.Errorf("cannot set non-integer property '%s' on a slice", property)
		}

		if i < 0 || i >= len(v) {
			return fmt.Errorf("cannot set out of range property '%s' on a slice", property)
		}

		v[i] = value
		return nil

	case map[string]interface{}:
		v[property] = value
		return nil

	case *[]interface{}:
		i, err := strconv.Atoi(property)

		if err != nil {
			return fmt.Errorf("cannot set non-integer property '%s' on a slice", property)
		}

		if i < 0 || i > len(*v) {
			return fmt.Errorf("cannot set out of range property '%s' on a slice", property)
		}

		if i == len(*v) {
			*v = append(*v, value)
			return nil
		}

		(*v)[i] = value
		return nil

	case *map[string]interface{}:
		(*v)[property] = value
		return nil

	default:
		return fmt.Errorf("cannot set property '%s' on type %T", property, target)
	}
}

// DefaultParser simply converts a string key into a list of properties that must be accessed in order, to achieve
// the dot notation get or set.
func DefaultParser(key string) []string {
	return strings.Split(key, ".")
}
