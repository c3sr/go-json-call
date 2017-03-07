package jsoncall

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/labstack/gommon/log"
)

// MarshalArgs converts its arguments in order into the json argument
// format compatible with go-json-call. Float values that are precise ints
// are converted to integers, unfortunately.
func MarshalArgs(args ...interface{}) ([]byte, error) {

	argMap := map[string]interface{}{}

	for argIdx, arg := range args {
		argIdxStr := strconv.Itoa(argIdx)
		argMap[argIdxStr] = arg
	}
	jsonBytes, err := json.Marshal(argMap)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

// UnmarshalArgs converts a json argument format into a slice of
// How to handle ambiguous json encodings (1.0 (go) -> 1 (json) -> 1 or 1.0 (go))
func UnmarshalArgs(jsonBytes []byte) ([]interface{}, error) {

	var args map[string]interface{}

	dec := json.NewDecoder(bytes.NewBuffer([]byte(jsonBytes)))
	dec.UseNumber()
	err := dec.Decode(&args)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%#v\n", args)

	// Convert to slice
	argSlice := make([]interface{}, 0)
	for i := 0; i < len(args); i++ {
		idxStr := strconv.Itoa(i)
		if arg, ok := args[idxStr]; ok {
			switch arg := arg.(type) {
			case json.Number: // try parsing as int, if fail, assume float
				i32, err := strconv.ParseInt(arg.String(), 10, 32)
				if err == nil {
					argSlice = append(argSlice, int(i32))
					continue
				}
				f64, err := strconv.ParseFloat(arg.String(), 64)
				if err != nil {
					return nil, err
				}
				argSlice = append(argSlice, f64)
			default:
				argSlice = append(argSlice, arg)
			}
		} else {
			log.Info("Didn't find argument", i)
			break
		}
	}
	return argSlice, nil
}

func canBeNil(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		return true
	}
	return false
}

func callFunction(args []interface{}, function interface{}) ([]reflect.Value, error) {
	functionTy := reflect.TypeOf(function)
	functionVal := reflect.ValueOf(function)

	// fmt.Println("Function has", functionTy.NumIn(), "parameters")
	// fmt.Println("Function has", functionTy.NumOut(), "results")
	// fmt.Println("Function sig:", functionTy)

	if functionTy.NumIn() != len(args) {
		return nil, errors.New("Function has " + strconv.Itoa(functionTy.NumIn()) + " params, got " + strconv.Itoa(len(args)) + " arguments.")
	}

	// Set up argument values
	argVals := make([]reflect.Value, len(args))
	for i, arg := range args {
		paramTy := functionTy.In(i)
		argTy := reflect.TypeOf(arg)

		if argTy == nil { // argument is nil interface
			if canBeNil(paramTy) {
				// argVals[i] = reflect.ValueOf(arg)
				argVals[i] = reflect.Zero(paramTy) // call with zero value
				continue
			} else {
				return nil, errors.New("Argument " + strconv.Itoa(i) + " type mismatch: got nil interface, but " + paramTy.String() + " cannot be nil")
			}
		}

		if paramTy != argTy {
			return nil, errors.New("Argument " + strconv.Itoa(i) + " type mismatch: got " + argTy.String() + " but needed " + paramTy.String())
		}

		argVals[i] = reflect.ValueOf(arg)
	}

	// Call the function
	resultVals := functionVal.Call(argVals)

	return resultVals, nil
}

func buildResults(vals []reflect.Value) (map[string]interface{}, error) {
	// Convert the results from values to go types
	results := map[string]interface{}{}
	for i, resultVal := range vals {
		idxStr := strconv.Itoa(i)
		resultKind := resultVal.Kind()
		// fmt.Println(resultTy, resultKind)

		switch resultKind {
		case reflect.Int:
			results[idxStr] = resultVal.Int()
		case reflect.Bool:
			results[idxStr] = resultVal.Bool()
		case reflect.String:
			results[idxStr] = resultVal.String()
		case reflect.Interface:
			results[idxStr] = resultVal.Interface()
		default:
			return nil, errors.New("Unhandled kind for result " + strconv.Itoa(i) + ": " + resultKind.String())
		}
	}
	return results, nil
}

// CallWithJSON Calls a function f with arguments from a JSON byte stream
func CallWithJSON(f interface{}, jsonBytes []byte) ([]byte, error) {
	args, err := UnmarshalArgs(jsonBytes)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v\n", args)

	vals, err := callFunction(args, f)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%#v\n", vals)

	results, err := buildResults(vals)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%#v\n", results)

	return json.Marshal(results)
}
