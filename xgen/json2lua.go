package xgen

import (
	"errors"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/sandwich-go/boost/xpanic"
)

func MustJSON2Lua(inBytes []byte) []byte {
	gen := NewGeneratorWithSpaceIndent(4)
	gen.P("return ")
	json2LuaWithType(inBytes, jsonparser.Object, gen, 1)
	return gen.Bytes()
}

func json2LuaWithType(inBytes []byte, dataType jsonparser.ValueType, gen *Gen, indent int) {
	if dataType == jsonparser.Object {
		gen.P("{")
		gen.In()
		err := jsonparser.ObjectEach(inBytes, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
			if dataType == jsonparser.Object {
				subGen := NewGeneratorWithSpaceIndent(4)
				subGen.SetPrefixIndent(indent - 1)
				subGen.In()
				json2LuaWithType(value, dataType, subGen, indent+1)
				gen.PFormat("[\"%s\"] = %s,", string(key), strings.TrimSpace(strings.TrimRight(subGen.String(), "\n")))
			} else if dataType == jsonparser.Array {
				subGen := NewGeneratorWithSpaceIndent(4)
				subGen.SetPrefixIndent(indent - 1)
				subGen.In()
				json2LuaWithType(value, dataType, subGen, indent+1)
				gen.PFormat("[\"%s\"] = %s,", string(key), strings.TrimSpace(strings.TrimRight(subGen.String(), "\n")))
			} else if dataType == jsonparser.String {
				gen.PFormat("[\"%s\"] = \"%s\",", string(key), string(value))
			} else if dataType == jsonparser.Boolean || dataType == jsonparser.Number || dataType == jsonparser.Null || dataType == jsonparser.Unknown {
				gen.PFormat("[\"%s\"] = %s,", string(key), string(value))
			} else if dataType == jsonparser.NotExist {
				return errors.New("unsupported type")
			}
			return nil
		})
		xpanic.WhenError(err)
		gen.Out()
		gen.P("}")
	}
	if dataType == jsonparser.Array {
		gen.P("{")
		gen.In()
		_, err := jsonparser.ArrayEach(inBytes, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			if dataType == jsonparser.Object || dataType == jsonparser.Array {
				subGen := NewGeneratorWithSpaceIndent(4)
				subGen.SetPrefixIndent(indent - 1)
				subGen.In()
				json2LuaWithType(value, dataType, subGen, indent+1)
				gen.PFormat("%s,", strings.TrimSpace(strings.TrimRight(subGen.String(), "\n")))
			} else if dataType == jsonparser.String {
				gen.PFormat("\"%s\",", string(value))
			} else if dataType == jsonparser.Boolean || dataType == jsonparser.Number || dataType == jsonparser.Null || dataType == jsonparser.Unknown {
				gen.PFormat("%s,", string(value))
			} else if dataType == jsonparser.NotExist {
				err = errors.New("unsupported type")
			}
		})
		xpanic.WhenError(err)
		gen.Out()
		gen.P("}")
	}
}
