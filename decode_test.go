package linecfg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	var (
		c   *ExampleConfig
		err error
	)

	// Happy case
	c = new(ExampleConfig)
	err = Decode("host=localhost port=8080 connect_timeout=5s A=3 B=4.5 C=1", c)
	assert.NoError(t, err)
	assert.Equal(t, &ExampleConfig{
		Host:           "localhost",
		Port:           8080,
		ConnectTimeout: 5 * time.Second,
		A:              3,
		B:              4.5,
		C:              true,
	}, c)

	// Unknown key
	c = new(ExampleConfig)
	err = Decode("host=localhost port=8080 connect_timeout=5s random_key=444", c)
	assert.Error(t, err)
	assert.Equal(t, "linecfg: unknown field 'random_key'", err.Error())

	// Bad format
	c = new(ExampleConfig)
	err = Decode("host =localhost port=8080 connect_timeout=5s", c)
	assert.Error(t, err)
	assert.Equal(t, "linecfg: bad field 'host' in 'host =localhost port=8080 connect_timeout=5s'", err.Error())

	c = new(ExampleConfig)
	err = Decode("host= localhost port=8080 connect_timeout=5s", c)
	assert.Error(t, err)
	assert.Equal(t, "linecfg: bad field 'localhost' in 'host= localhost port=8080 connect_timeout=5s'", err.Error())

	// Bad keys
	c = new(ExampleConfig)
	err = Decode("unsupported=5", c)
	assert.Error(t, err)
	assert.Equal(t, "linecfg: cannot unmarshal 5 into Go value of type time.Time", err.Error())

	// Good format
	c = new(ExampleConfig)
	err = Decode(" host=localhost  port=8080   connect_timeout=5s  ", c)
	assert.NoError(t, err)

	// Bad value
	err = Decode("", nil)
	assert.Error(t, err)
	assert.Equal(t, "linecfg: invalid type: nil", err.Error())

	err = Decode("", struct{}{})
	assert.Error(t, err)
	assert.Equal(t, "linecfg: invalid type: non-pointer struct {}", err.Error())
}

type ExampleConfig struct {
	Host           string
	Port           int
	ConnectTimeout time.Duration `linecfg:"connect_timeout"`
	A              uint32
	B              float32
	C              bool
	Unsupported    time.Time
}
