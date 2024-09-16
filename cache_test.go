package gonvoy

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache_StoreAndLoad(t *testing.T) {
	lc := newInternalCache()

	t.Run("pointer object", func(t *testing.T) {
		source := bytes.NewReader([]byte("testing"))
		lc.Store("foo", source)

		receiver := new(bytes.Reader)
		_, ok, err := LoadValue[*bytes.Reader](lc, "foo")
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, source, receiver)
	})

	t.Run("literal object", func(t *testing.T) {
		type mystruct struct{}
		src := mystruct{}
		lc.Store("bar", src)

		dest := mystruct{}
		_, ok, err := LoadValue[mystruct](lc, "bar")
		require.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, src, dest)
	})

	t.Run("a nil receiver, returns an error", func(t *testing.T) {
		type mystruct struct{}
		_, ok, err := LoadValue[mystruct](lc, "bar")
		assert.False(t, ok)
		assert.ErrorIs(t, err, ErrNilReceiver)
	})

	t.Run("receiver has incompatibility data type with the source, returns an error", func(t *testing.T) {
		type mystruct struct{}
		src := new(mystruct)
		lc.Store("foobar", src)

		_, ok, err := LoadValue[mystruct](lc, "foobar")
		assert.False(t, ok)
		assert.ErrorIs(t, err, ErrIncompatibleReceiver)
	})

	t.Run("if no data found during a Load, then returns false without an error", func(t *testing.T) {
		_, ok, err := LoadValue[struct{}](lc, "data-not-exists")
		assert.False(t, ok)
		assert.NoError(t, err)
	})
}
