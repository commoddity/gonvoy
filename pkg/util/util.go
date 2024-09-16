package util

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// ReplaceAllEmptySpace replaces all empty space as well as reserved escape characters such as
// tab, newline, carriage return, and so forth.
func ReplaceAllEmptySpace(s string) string {
	replacementMaps := []string{
		" ", "_",
		"\t", "_",
		"\n", "_",
		"\v", "_",
		"\r", "_",
		"\f", "_",
	}

	replacer := strings.NewReplacer(replacementMaps...)

	return replacer.Replace(s)
}

// NewFrom dynamically creates a new variable from the specified data type.
// However, the returned Value's Type is always a PointerTo{dataType}.
func NewFrom(in interface{}) (out interface{}, err error) {
	f := reflect.TypeOf(in)
	var v reflect.Value
	switch f.Kind() {
	case reflect.Ptr:
		v = reflect.New(f.Elem())
	case reflect.Struct:
		v = reflect.New(f)
	default:
		return out, fmt.Errorf("data type must be either pointer to struct or literal struct. Got %s instead", f.Name())
	}

	return v.Interface(), nil
}

func IsNil(i interface{}) bool {
	v := reflect.ValueOf(i)
	if !v.IsValid() {
		return true
	}

	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// In returns true if a given value exists in the list.
func In[T comparable](value T, list ...T) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// GetAbsPathFromCaller returns the absolute path of the file that calls this function.
func GetAbsPathFromCaller(skip int) (string, error) {
	if skip < 0 {
		skip = 0
	}

	_, file, _, ok := runtime.Caller(skip)
	if !ok {
		return "", os.ErrNotExist
	}

	absPath, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	return filepath.Dir(absPath), nil
}

// StringStartsWith checks if a string starts with any of the prefixes.
func StringStartsWith(s string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}

	return false
}
