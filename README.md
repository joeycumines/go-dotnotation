# go-dotnotation

A simple dot notation accessor, to get and set values from data structures.
The intended use case is to simplify retrieving information from unstructured or
data in unknown formats, decoded from JSON, but it can be adapted for other
purposes by providing custom handlers.

The package can be imported as `github.com/joeycumines/go-dotnotation/dotnotation`.

```go

// Accessor provides two methods, Get and Set, that can be configured to handle custom data structures via the
// exported properties, Parser, Getter, and Setter.
type Accessor struct {
	// Getter returns the property value of a given target, or an error.
	Getter func(target interface{}, property string) (interface{}, error)
	// Setter sets the property value of a given target, to a given value, or returns an error.
	Setter func(target interface{}, property string, value interface{}) error
	// Parser converts a given key into a list of properties to access in order to get or set.
	Parser func(key string) []string
}

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
```

## Notes

- `DefaultParser` just splits on `.`
- `DefaultGetter` and `DefaultSetter` support `[]interface{}`,
    `map[string]interface{}`, as well as those two types with one level of
    pointer indirection (`*[]interface{}` and `*map[string]interface{}`)
- Setting the next index (like `len(slice)`) of a `*[]interface{}` type
    will append to the slice.
