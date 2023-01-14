package optional

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyOptional(t *testing.T) {
	opt := NewEmpty[int]()
	assert.False(t, opt.HasValue())
}

func TestOptionalWithValue(t *testing.T) {
	opt := New(42)
	assert.True(t, opt.HasValue())
}

func TestOptionalValue(t *testing.T) {
	opt := New(42)
	assert.Equal(t, opt.Value(), 42)
}

func TestEmptyOptionalValue(t *testing.T) {
	opt := NewEmpty[int]()
	assert.Panics(t, func() { opt.Value() })
}

func TestOptionalValueOr(t *testing.T) {
	opt := New(42)
	assert.Equal(t, opt.ValueOr(1337), 42)
}

func TestEmptyOptionalValueOr(t *testing.T) {
	opt := NewEmpty[int]()
	assert.Equal(t, opt.ValueOr(1337), 1337)
}

func TestEmptyOptionalClear(t *testing.T) {
	opt := New(42)
	opt.Clear()
	assert.False(t, opt.HasValue())
}

func TestInlineOptionalMethodCall(t *testing.T) {
	// If Optional methods have a pointer receiver, they can't be called on a
	// non-addressable value (such as the direct return of New or NewEmpty)
	assert.True(t, New(42).HasValue())
	assert.Equal(t, New(42).Value(), 42)
	assert.Equal(t, NewEmpty[int]().ValueOr(1337), 1337)
}
