package linecfg

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	var (
		line string
		err  error
	)

	// Happy path
	line, err = Encode(&struct {
		Foo string
		Bar int `linecfg:"override"`
	}{"hello", 3})
	assert.NoError(t, err)
	assert.Equal(t, "foo=hello override=3", line)

	// Hide null values
	line, err = Encode(&struct {
		Foo string
		Bar int
	}{"", 0})
	assert.NoError(t, err)
	assert.Equal(t, "", line)

	// Bad data
	line, err = Encode(nil)
	assert.Equal(t, ErrInvalidType{reflect.TypeOf(nil)}, err)
	assert.Equal(t, "", line)
}
