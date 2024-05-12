package gonvoy

import (
	"encoding/json"
	"time"

	"github.com/ardikabs/gonvoy/pkg/util"
	"github.com/envoyproxy/envoy/contrib/golang/common/go/api"
)

// MustGetProperty is an extended of GetProperty, only panic if value is not in acceptable format.
func MustGetProperty(c RuntimeContext, name, defaultVal string) string {
	value, err := c.GetProperty(name, defaultVal)
	if err != nil {
		panic(err)
	}

	return value
}

// NewMinimalJSONResponse creates a minimal JSON body as a form of bytes.
func NewMinimalJSONResponse(code, message string, errs ...error) []byte {
	bodyMap := make(map[string]interface{})
	bodyMap["code"] = code
	bodyMap["message"] = message
	bodyMap["errors"] = nil
	bodyMap["data"] = make(map[string]interface{}, 0)
	bodyMap["serverTime"] = time.Now().UnixMilli()

	listErrs := make([]string, len(errs))
	for i, err := range errs {
		listErrs[i] = err.Error()
	}
	bodyMap["errors"] = listErrs

	bodyByte, err := json.Marshal(bodyMap)
	if err != nil {
		bodyByte = []byte("{}")
	}

	return bodyByte
}

func checkBodyAccessibility(strict, allowRead, allowWrite bool, header api.HeaderMap) (read, write bool) {
	access := isBodyAccessible(header)

	if !strict {
		read = access && (allowRead || allowWrite)
		write = access && allowWrite
		return
	}

	operation, ok := header.Get(HeaderXContentOperation)
	if !ok {
		return
	}

	if util.In(operation, ContentOperationReadOnly, ContentOperationRO) {
		read = access && allowRead
		return
	}

	if util.In(operation, ContentOperationReadWrite, ContentOperationRW) {
		write = access && allowWrite
		read = write
		return
	}

	return
}

func isBodyAccessible(header api.HeaderMap) bool {
	contentLength, ok := header.Get(HeaderContentLength)
	if !ok {
		return false
	}

	isEmpty := contentLength == "" || contentLength == "0"
	return !isEmpty
}

func isRequestBodyAccessible(c Context) bool {
	return c.IsRequestBodyReadable() || c.IsRequestBodyWriteable()
}

func isResponseBodyAccessible(c Context) bool {
	return c.IsResponseBodyReadable() || c.IsResponseBodyWriteable()
}
