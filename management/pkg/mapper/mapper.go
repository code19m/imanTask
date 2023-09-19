package mapper

import (
	jsoniter "github.com/json-iterator/go"
)

// Map function facilitates the mapping of data from the source object to the destination object.
//
// Note: The output parameter should be a pointer to the destination object to allow the Map function to modify
// the destination object directly. If the output parameter is not a pointer, the Map function will not be able
// to modify the destination object's fields.
//
// This function relies on JSON marshaling and unmarshaling, so it might not be the most performant option
// for mapping large amounts of data or for situations requiring highly efficient data copying.
func Map(input interface{}, output interface{}) error {
	rawBytes, err := jsoniter.Marshal(input)
	if err != nil {
		return err
	}
	return jsoniter.Unmarshal(rawBytes, output)
}
