package gonvoy

import (
	"strconv"
	"testing"

	mock_envoy "github.com/commoddity/gonvoy/test/mock/envoy"
	"github.com/stretchr/testify/assert"
)

func TestBodyRead(t *testing.T) {
	b := &bodyWriter{}

	assert.Zero(t, b.Bytes())
	assert.Zero(t, b.String())

	bufferMock := mock_envoy.NewBufferInstance(t)

	bw := &bodyWriter{
		buffer: bufferMock,
		bytes:  []byte("lorem_ipsum"),
	}

	assert.Equal(t, "lorem_ipsum", bw.String())
	assert.Equal(t, []byte("lorem_ipsum"), bw.Bytes())
}

func TestBody_Write(t *testing.T) {
	input := []byte(`{"name":"John Doe"}`)

	t.Run("body writer is writable, returns no error", func(t *testing.T) {
		headerMock := mock_envoy.NewRequestHeaderMap(t)
		headerMock.EXPECT().Get(HeaderContentLength).Return("", true)
		headerMock.EXPECT().Set(HeaderContentLength, strconv.Itoa(len(input)))

		bufferMock := mock_envoy.NewBufferInstance(t)
		bufferMock.EXPECT().Set(input).Return(nil)
		bufferMock.EXPECT().Len().Return(len(input))

		writer := &bodyWriter{
			writable:              true,
			preserveContentLength: true,
			header:                headerMock,
			buffer:                bufferMock,
		}

		n, err := writer.Write(input)
		assert.NoError(t, err)
		assert.Equal(t, len(input), n)
	})

	t.Run("body writer is not writable, returns an error", func(t *testing.T) {
		headerMock := mock_envoy.NewRequestHeaderMap(t)
		bufferMock := mock_envoy.NewBufferInstance(t)
		writer := &bodyWriter{
			header: headerMock,
			buffer: bufferMock,
		}

		n, err := writer.Write(input)
		assert.ErrorIs(t, err, ErrOperationNotPermitted)
		assert.Zero(t, n)
	})

	t.Run("body writer is writable, but preserveContentLength is enabled", func(t *testing.T) {
		headerMock := mock_envoy.NewRequestHeaderMap(t)

		bufferMock := mock_envoy.NewBufferInstance(t)
		bufferMock.EXPECT().Set(input).Return(nil)
		bufferMock.EXPECT().Len().Return(len(input))

		writer := &bodyWriter{
			writable:              true,
			preserveContentLength: false,
			header:                headerMock,
			buffer:                bufferMock,
		}

		n, err := writer.Write(input)
		assert.NoError(t, err)
		assert.Equal(t, len(input), n)
	})
}

func TestBody_WriteString(t *testing.T) {
	input := "new data"

	t.Run("body writer is writable, returns no error", func(t *testing.T) {
		headerMock := mock_envoy.NewRequestHeaderMap(t)
		headerMock.EXPECT().Get(HeaderContentLength).Return("", true)
		headerMock.EXPECT().Set(HeaderContentLength, strconv.Itoa(len(input)))

		bufferMock := mock_envoy.NewBufferInstance(t)
		bufferMock.EXPECT().SetString(input).Return(nil)
		bufferMock.EXPECT().Len().Return(len(input))

		writer := &bodyWriter{
			writable:              true,
			preserveContentLength: true,
			header:                headerMock,
			buffer:                bufferMock,
		}

		n, err := writer.WriteString(input)
		assert.NoError(t, err)
		assert.Equal(t, len(input), n)
	})

	t.Run("body writer is not writable, returns an error", func(t *testing.T) {
		headerMock := mock_envoy.NewRequestHeaderMap(t)
		bufferMock := mock_envoy.NewBufferInstance(t)
		writer := &bodyWriter{
			header: headerMock,
			buffer: bufferMock,
		}

		n, err := writer.WriteString(input)
		assert.ErrorIs(t, err, ErrOperationNotPermitted)
		assert.Zero(t, n)
	})
}
