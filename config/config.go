// Package config holds global project configuration
package config

type config map[string]any

var cfg config = map[string]any{}

// Get retrieves a value in the project config
func Get(key string) (any, bool) {
	v, ok := cfg[key]
	return v, ok
}

// GetOr retrieves a value in the project config if it is defined, or returns
// the fallback value
func GetOr(key string, fallback any) any {
	if v, ok := cfg[key]; ok {
		return v
	}
	return fallback
}

// GetAll returns the entire project config
func GetAll() map[string]any {
	return cfg
}

// Set sets a value in the project config
func Set(key string, value any) {
	cfg[key] = value
}
