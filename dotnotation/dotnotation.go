// Package dotnotation provides dot notation getters and setters for manipulating data structures.
package dotnotation

// DefaultAccessor, used for the exported Set and Get functions.
var DefaultAccessor Accessor

// Set sets a value using dot notation, by default it supports generic []interface{} and map[string]interface{} types.
// It's behaviour can be configured by modifying the DefaultAccessor variable.
func Set(target interface{}, key string, value interface{}) error {
	return DefaultAccessor.Set(target, key, value)
}

// Get sets a value using dot notation, by default it supports generic []interface{} and map[string]interface{} types.
// It's behaviour can be configured by modifying the DefaultAccessor variable.
func Get(target interface{}, key string) (interface{}, error) {
	return DefaultAccessor.Get(target, key)
}
