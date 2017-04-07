package serializedcall

import (
	"errors"
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

func TestMarshalXML(t *testing.T) {
	_, err := MarshalXMLArgs(1, 2)
	assert.NoError(t, err, "Support ints")

	_, err = MarshalXMLArgs(-1, 0)
	assert.NoError(t, err, "Support negative ints")

	_, err = MarshalXMLArgs(1.0)
	assert.NoError(t, err, "Support float64")

	_, err = MarshalXMLArgs("1")
	assert.NoError(t, err, "Support strings")

	_, err = MarshalXMLArgs(1, "")
	assert.NoError(t, err, "Support mixed")

	_, err = MarshalXMLArgs(nil)
	assert.NoError(t, err, "Support nil")

	type TestType struct {
		a int
		b string
	}

	_, err = MarshalXMLArgs(TestType{1, ""})
	assert.NoError(t, err, "Support struct")

	_, err = MarshalXMLArgs(&TestType{1, ""})
	assert.NoError(t, err, "Support struct pointer")
}

func TestUnmarshalXML(t *testing.T) {
	j, err := MarshalXMLArgs(1, 0, -1)
	assert.NoError(t, err, "Marshal ints")
	args, err := UnmarshalXMLArgs(j)
	assert.NoError(t, err, "Unmarshal ints")
	assert.Equal(t, 1, args[0].(int), "Support ints")
	assert.Equal(t, 0, args[1].(int), "Support ints")
	assert.Equal(t, -1, args[2].(int), "Support ints")

	j, err = MarshalXMLArgs(1.0)
	pp.Println(string(j))
	assert.NoError(t, err, "Support f64")
	args, err = UnmarshalXMLArgs(j)
	assert.NoError(t, err, "Support f64")
	assert.IsType(t, 1.0, args[0], "Should be float64")
	pp.Println(args)
	assert.Equal(t, 1.0, args[0].(float64), "Should be 1.0")

}

func TestCallWithXML(t *testing.T) {

	testfunc1 := func(x int) (int, error) {
		if x == 0 {
			return 0, errors.New("")
		} else {
			return x, nil
		}
	}

	args := []byte(`{"0":0}`)
	_, err := CallWithXML(testfunc1, args)
	assert.NoError(t, err, "")

	args = []byte(`{"0":null}`)
	_, err = CallWithXML(testfunc1, args)
	assert.Error(t, err, "")

}

func TestNilArgCall(t *testing.T) {

	testfunc1 := func(x *int) int {
		if x == nil {
			return 0
		}
		return *x
	}

	args := []byte(`{"0":null}`)
	_, err := CallWithXML(testfunc1, args)
	assert.NoError(t, err, "")

}
