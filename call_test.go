package jsoncall

import (
	"errors"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	_, err := MarshalArgs(1, 2)
	assert.NoError(t, err, "Support ints")

	_, err = MarshalArgs(-1, 0)
	assert.NoError(t, err, "Support negative ints")

	_, err = MarshalArgs(1.0)
	assert.NoError(t, err, "Support float64")

	_, err = MarshalArgs("1")
	assert.NoError(t, err, "Support strings")

	_, err = MarshalArgs(1, "")
	assert.NoError(t, err, "Support mixed")

	_, err = MarshalArgs(nil)
	assert.NoError(t, err, "Support nil")

	type TestType struct {
		a int
		b string
	}

	_, err = MarshalArgs(TestType{1, ""})
	assert.NoError(t, err, "Support struct")

	_, err = MarshalArgs(&TestType{1, ""})
	assert.NoError(t, err, "Support struct pointer")
}

func TestUnmarshal(t *testing.T) {
	j, err := MarshalArgs(1, 0, -1)
	assert.NoError(t, err, "Support ints")
	args, err := UnmarshalArgs(j)
	assert.NoError(t, err, "Support ints")
	assert.Equal(t, 1, args[0].(int), "Support ints")
	assert.Equal(t, 0, args[1].(int), "Support ints")
	assert.Equal(t, -1, args[2].(int), "Support ints")

	j, err = MarshalArgs(1.0)
	pp.Println(string(j))
	assert.NoError(t, err, "Support f64")
	args, err = UnmarshalArgs(j)
	assert.NoError(t, err, "Support f64")
	assert.IsType(t, 1.0, args[0], "Should be float64")
	pp.Println(args)
	assert.Equal(t, 1.0, args[0].(float64), "Should be 1.0")

}

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
