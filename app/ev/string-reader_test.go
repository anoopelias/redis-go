package ev

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringRead(t *testing.T) {
	arr := []byte{67, 104, 101, 99, 111, 10}
	asr := NewArrayStringReader(arr)

	str, err := asr.ReadString('\n')
	assert.Nil(t, err)
	assert.Equal(t, "Checo\n", str)
}

func TestStringReadMulti(t *testing.T) {
	arr := []byte{67, 104, 101, 99, 111, 10, 80, 101, 114, 101, 122, 10}
	asr := NewArrayStringReader(arr)

	str, err := asr.ReadString('\n')
	assert.Nil(t, err)
	assert.Equal(t, "Checo\n", str)

	str, err = asr.ReadString('\n')
	assert.Nil(t, err)
	assert.Equal(t, "Perez\n", str)
}

func TestStringReadError(t *testing.T) {
	arr := []byte{67, 104, 101, 99, 111, 10}
	asr := NewArrayStringReader(arr)

	asr.ReadString('\n')
	_, err := asr.ReadString('\n')
	assert.NotNil(t, err)
}
