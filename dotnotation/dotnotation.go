package dotnotation

// Parser is the entry point to the implementation that this package provides, and encapsulates configuration
// as well as the actual logic to get and set via dot or bracket notation.
type Parser struct {
	// Getter returns the property value of a given target, or an error.
	Getter func(target interface{}, property string) (interface{}, error)

	// Setter sets the property value of a given target, to a given value, or returns an error.
	Setter func(target interface{}, property string, value interface{}) error
}
