package dotnotation

import "errors"

// Accessor provides two methods, Get and Set, that can be configured to handle custom data structures via the
// exported properties, Getter and Setter.
type Accessor struct {
	// Getter returns the property value of a given target, or an error.
	Getter func(target interface{}, property string) (interface{}, error)
	// Setter sets the property value of a given target, to a given value, or returns an error.
	Setter func(target interface{}, property string, value interface{}) error
	// Parser converts a given key into a list of properties to access in order to get or set.
	Parser func(key string) []string
}

func (p Accessor) Set(target interface{}, key string, value interface{}) error {
	properties := p.parser(key)

	for i, property := range properties {
		if i == (len(properties) - 1) {
			// we reached the last property
			return p.setter(target, property, value)
		}

		// attempt to get the next level, so we can set the last property
		var err error
		target, err = p.getter(target, property)
		if err != nil {
			return err
		}
	}

	return errors.New("no properties parsed from key: " + key)
}

func (p Accessor) Get(target interface{}, key string) (interface{}, error) {
	properties := p.parser(key)

	for i, property := range properties {
		if i == (len(properties) - 1) {
			// we reached the last property
			return p.getter(target, property)
		}

		// attempt to get the next level
		var err error
		target, err = p.getter(target, property)
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("no properties parsed from key: " + key)
}

func (p Accessor) getter(target interface{}, property string) (interface{}, error) {
	if p.Getter == nil {
		return DefaultGetter(target, property)
	}

	return p.Getter(target, property)
}

func (p Accessor) setter(target interface{}, property string, value interface{}) error {
	if p.Setter == nil {
		return DefaultSetter(target, property, value)
	}

	return p.Setter(target, property, value)
}

func (p Accessor) parser(key string) []string {
	if p.Parser == nil {
		return DefaultParser(key)
	}

	return p.Parser(key)
}
