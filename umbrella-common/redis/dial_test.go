package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDial(t *testing.T) {
	assert := assert.New(t)

	_, err := Dial("tcp", "127.0.0.1:6379")
	assert.Nil(err)
}

func TestDialTimeout(t *testing.T) {
	assert := assert.New(t)

	_, err := DialTimeout("tcp", "127.0.0.1:6379", time.Duration(time.Millisecond))
	assert.Nil(err)
}
