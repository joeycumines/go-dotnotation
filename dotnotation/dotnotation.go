// Package dotnotation provides dot notation getters and setters for manipulating data structures.
package dotnotation

var DefaultAccessor = Accessor{}

func Set(target interface{}, key string, value interface{}) error {
	return DefaultAccessor.Set(target, key, value)
}

func Get(target interface{}, key string) (interface{}, error) {
	return DefaultAccessor.Get(target, key)
}
