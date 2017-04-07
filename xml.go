package serializedcall

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/k0kubun/pp"
)

// MarshalXMLArgs converts a list of arguments into XML
func MarshalXMLArgs(args ...interface{}) ([]byte, error) {

	argSlice := []interface{}{}

	for _, arg := range args {
		argSlice = append(argSlice, arg)
	}
	xmlBytes, err := json.Marshal(argSlice)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(xmlBytes))
	return xmlBytes, nil
}

// UnmarshalXMLArgs converts XML into a slice of go data
func UnmarshalXMLArgs(xmlBytes []byte) ([]interface{}, error) {

	var args []interface{}

	dec := xml.NewDecoder(bytes.NewBuffer([]byte(xmlBytes)))
	err := dec.Decode(&args)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%#v\n", args)

	return args, nil
}

// CallWithXML Calls a function f with arguments from a XML byte stream
func CallWithXML(f interface{}, xmlBytes []byte) ([]byte, error) {
	args, err := UnmarshalXMLArgs(xmlBytes)
	if err != nil {
		return nil, err
	}

	pp.Println(args)

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
