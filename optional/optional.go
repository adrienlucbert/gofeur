// Package optional provides an implementation for an optional type wrapper.
package optional

// Optional manages an optional contained value.
type Optional[T any] struct {
	ptr *T
}

// New initializes an optional containing the given value.
func New[T any](v T) Optional[T] {
	return Optional[T]{&v}
}

// NewEmpty initializes an optional containing no value.
func NewEmpty[T any]() Optional[T] {
	return Optional[T]{}
}

// HasValue returns true if the optional holds a value.
func (o Optional[T]) HasValue() bool {
	return o.ptr != nil
}

// Value returns the value held in the optional.
func (o Optional[T]) Value() T {
	if !o.HasValue() {
		panic("Trying to retrieve value from an empty optional")
	}
	return *o.ptr
}

// Clear sets the value held in the optional to nil.
func (o *Optional[T]) Clear() {
	o.ptr = nil
}

// Set sets the value held in the optional.
func (o *Optional[T]) Set(v T) {
	o.ptr = &v
}

// ValueOr returns the value held in the optional if available, defaultValue
// otherwise.
func (o Optional[T]) ValueOr(defaultValue T) T {
	if !o.HasValue() {
		return defaultValue
	}
	return *o.ptr
}
