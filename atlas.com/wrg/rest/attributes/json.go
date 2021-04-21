package attributes

import (
	"encoding/json"
	"io"
)

// ToJSON serializes the given interface into a string based JSON format
func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

// FromJSON deserializes the object from JSON string
// in an io.Reader to the given interface
func FromJSON(i interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(i)
}

// objectMap is a simple representation of a json map. key is a string, and value is a nested object
type objectMap map[string]interface{}

// objectMapper maps a objectMap to a concrete data type
type objectMapper func(objectMap) interface{}

// conditionalMapperProvider returns a string representing the type of object the objectMapper handles, as well as the
// objectMapper itself.
type conditionalMapperProvider func() (string, objectMapper)

// concreteObjectProvider returns a new (empty) concrete object
type concreteObjectProvider func() interface{}

// dataSegment is a slice of concrete objects
type dataSegment []interface{}

// unmarshalRoot will take a raw byte array, and using mapper functions, produce a data dataSegment and includes
// dataSegment representing a jsonapi.org request response
func unmarshalRoot(data []byte, root objectMapper, options ...conditionalMapperProvider) (dataSegment, dataSegment, error) {
	var dataResult dataSegment
	var includeResult dataSegment

	var single = struct {
		Data     objectMap
		Included []objectMap
	}{}
	err := json.Unmarshal(data, &single)
	if err == nil {
		dataResult = produceSingle(single.Data, root)
		includeResult = produceList(single.Included, includeMapper(options...))
		return dataResult, includeResult, nil
	} else {
		var list = struct {
			Data     []objectMap
			Included []objectMap
		}{}
		err = json.Unmarshal(data, &list)
		if err == nil {
			dataResult = produceList(list.Data, root)
			includeResult = produceList(list.Included, includeMapper(options...))
			return dataResult, includeResult, nil
		}
		return nil, nil, err
	}
}

// includeMapper represents a objectMapper for handling the jsonapi.org includes data section
func includeMapper(options ...conditionalMapperProvider) objectMapper {
	return func(o objectMap) interface{} {
		return addInclude(o, options...)
	}
}

// produceSingle will produce a dataSegment given a single objectMap from a objectMapper
func produceSingle(o objectMap, m objectMapper) dataSegment {
	return append(make([]interface{}, 0), m(o))
}

// produceList will produce a dataSegment given a objectMap slice from a objectMapper
func produceList(o []objectMap, m objectMapper) dataSegment {
	if len(o) > 0 {
		var result = make([]interface{}, 0)
		for _, x := range o {
			result = append(result, m(x))
		}
		return result
	}
	return nil
}

// transformMap will populate the concrete object provided by the concreteObjectProvider from the objectMap and return it.
func transformMap(dpf concreteObjectProvider, x objectMap) interface{} {
	b, err := json.Marshal(x)
	if err == nil {
		r := dpf()
		err = json.Unmarshal(b, r)
		if err == nil {
			return r
		}
	}
	return nil
}

func mapperFunc(dpf concreteObjectProvider) objectMapper {
	return func(x objectMap) interface{} {
		return transformMap(dpf, x)
	}
}

func unmarshalData(tf string, dpf concreteObjectProvider) (string, objectMapper) {
	return tf, mapperFunc(dpf)
}

// addInclude processes a map structure representing the jsonapi.org data object to produce a concrete struct.
// given a set of functions which could produce a concrete struct
func addInclude(x objectMap, options ...conditionalMapperProvider) interface{} {
	t := x["type"].(string)
	for _, o := range options {
		tf, mf := o()
		if t == tf {
			return mf(x)
		}
	}
	return nil
}
