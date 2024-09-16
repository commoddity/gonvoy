package util_test

import (
	"os"
	"testing"

	"github.com/commoddity/gonvoy/pkg/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplaceAllEmptySpace(t *testing.T) {
	message := "via Go extension"

	replaced := util.ReplaceAllEmptySpace(message)
	assert.Equal(t, "via_Go_extension", replaced)
}

func TestNewFrom(t *testing.T) {
	type dummyStruct struct {
		A string
		B string
		c string
	}

	src := dummyStruct{
		A: "foo",
		B: "bar",
		c: "foobar_private",
	}

	dest, err := util.NewFrom(src)
	require.NoError(t, err)
	assert.IsType(t, &dummyStruct{}, dest)
	assert.Zero(t, dest.(*dummyStruct).A)
	assert.Zero(t, dest.(*dummyStruct).B)
	assert.Zero(t, dest.(*dummyStruct).c)

	srcPtr := &dummyStruct{
		A: "foo",
		B: "bar",
		c: "foobar_private",
	}
	destPtr, err := util.NewFrom(srcPtr)
	require.NoError(t, err)
	assert.IsType(t, srcPtr, destPtr)
	assert.Zero(t, destPtr.(*dummyStruct).A)
	assert.Zero(t, destPtr.(*dummyStruct).B)
	assert.Zero(t, destPtr.(*dummyStruct).c)
	assert.NotSame(t, src, dest)
	assert.NotSame(t, srcPtr, destPtr)
}

func TestIsNil(t *testing.T) {
	var (
		a []int
		b *string
		c interface{}
		d int
		e string
		f struct{}
	)

	assert.True(t, util.IsNil(a))
	assert.True(t, util.IsNil(b))

	c = a
	assert.True(t, util.IsNil(c))

	assert.False(t, util.IsNil(d))
	assert.False(t, util.IsNil(e))
	assert.False(t, util.IsNil(f))
}

func TestIn(t *testing.T) {

	t.Run("true", func(t *testing.T) {
		val := util.In("woman", "man", "woman")
		assert.True(t, val)
	})

	t.Run("false", func(t *testing.T) {
		val := util.In(5, 1, 2, 3, 4)
		assert.False(t, val)
	})
}

func TestGetAbsPathFromCaller(t *testing.T) {
	path, err := util.GetAbsPathFromCaller(0)
	assert.NoError(t, err)

	wd, err := os.Getwd()
	assert.NoError(t, err)
	assert.Equal(t, wd, path)
}

func TestStringStartsWith(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		result := util.StringStartsWith("", "prefix")
		assert.False(t, result)
	})

	t.Run("string starts with prefix", func(t *testing.T) {
		result := util.StringStartsWith("prefix123", "prefix")
		assert.True(t, result)
	})

	t.Run("string does not start with prefix", func(t *testing.T) {
		result := util.StringStartsWith("suffix123", "prefix")
		assert.False(t, result)
	})

	t.Run("string starts with any of the prefixes", func(t *testing.T) {
		result := util.StringStartsWith("prefix123", "suffix", "prefix")
		assert.True(t, result)
	})

	t.Run("string does not start with any of the prefixes", func(t *testing.T) {
		result := util.StringStartsWith("suffix123", "prefix1", "prefix2")
		assert.False(t, result)
	})
}
