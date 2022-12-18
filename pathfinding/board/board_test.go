package board

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBoardString(t *testing.T) {
	b := New(3, 2)
	assert.Equal(t, b.String(), "· · · \n· · · \n")
	b[0][0].Blocked = true
	assert.Equal(t, b.String(), "# · · \n· · · \n")
}

func TestBoardWidth(t *testing.T) {
	b := New(3, 2)
	assert.Equal(t, b.Width(), uint(3))
}

func TestBoardHeight(t *testing.T) {
	b := New(3, 2)
	assert.Equal(t, b.Height(), uint(2))
}

func TestBoardAt(t *testing.T) {
	b := New(3, 2)
	b.At(1, 1).Blocked = true
	b.At(2, 0).Blocked = true
	assert.Equal(t, b.String(), "· · # \n· # · \n")
}

func TestBoardContains(t *testing.T) {
	b := New(3, 2)
	assert.False(t, b.Contains(13, 37))
	assert.False(t, b.Contains(3, 0))
	assert.False(t, b.Contains(0, 2))
	assert.True(t, b.Contains(0, 0))
	assert.True(t, b.Contains(2, 1))
}
