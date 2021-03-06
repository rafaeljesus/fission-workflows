package typedvalues

import (
	"encoding/json"

	"fmt"

	"github.com/fission/fission-workflows/pkg/types"
)

type Object map[string]interface{}

// The additional subtype (e.g. json/string) is needed to be able to evaluate inline functions on specific json types (e.g. strings vs. arrays)
const (
	FORMAT_JSON = "json"
	TYPE_STRING = "string"
	TYPE_OBJECT = "object"
	TYPE_INT    = "int"
	TYPE_ARRAY  = "array"
	TYPE_BOOL   = "bool"
)

var JSON_TYPES = []string{
	TYPE_STRING,
	TYPE_OBJECT,
	TYPE_ARRAY,
	TYPE_BOOL,
	TYPE_INT,
}

func isJsonValue(val *types.TypedValue) bool {
	for _, vtype := range JSON_TYPES {
		if FormatType(FORMAT_JSON, vtype) == val.Type {
			return true
		}
	}
	return false
}

type JsonParserFormatter struct{}

func (jp *JsonParserFormatter) Parse(i interface{}, allowedTypes ...string) (*types.TypedValue, error) {
	return jp.parsePrimitives(i, allowedTypes...)
}

func (jp *JsonParserFormatter) parsePrimitives(i interface{}, allowedTypes ...string) (*types.TypedValue, error) {
	var tp string
	switch i.(type) {
	case bool:
		tp = TYPE_BOOL
	case string:
		tp = TYPE_STRING
	case map[string]interface{}:
		tp = TYPE_OBJECT
	case []interface{}:
		tp = TYPE_ARRAY
	case float64:
		tp = TYPE_INT
	case int64:
		tp = TYPE_INT
	case int:
		tp = TYPE_INT
	default:
		return nil, fmt.Errorf("value '%v' cannot be parsed to json", i)
	}

	bs, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	return &types.TypedValue{
		Type:  FormatType(FORMAT_JSON, tp),
		Value: bs,
	}, nil
}

func (jp *JsonParserFormatter) Format(v *types.TypedValue) (interface{}, error) {
	if !isJsonValue(v) {
		return nil, fmt.Errorf("value '%v' is not a JSON type", v)
	}

	var i interface{}
	err := json.Unmarshal(v.Value, &i)
	return i, err
}
