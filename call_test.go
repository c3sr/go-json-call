package jsoncall

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T) {

	testfunc1 := func(x int) (int, error) {
		if x == 0 {
			return 0, errors.New("")
		} else {
			return x, nil
		}
	}

	args := []byte(`{"0":0}`)
	_, err := CallWithJSON(testfunc1, args)
	assert.NoError(t, err, "")

	args = []byte(`{"0":null}`)
	_, err = CallWithJSON(testfunc1, args)
	assert.Error(t, err, "")

}

func TestNilArg(t *testing.T) {

	testfunc1 := func(x *int) int {
		if x == nil {
			return 0
		}
		return *x
	}

	args := []byte(`{"0":null}`)
	_, err := CallWithJSON(testfunc1, args)
	assert.NoError(t, err, "")

}
