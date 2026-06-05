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
			}
			// NotExist 分支：jsonparser.ArrayEach 的 callback 签名不返回 error，
			// 外层 _, err 检查的是 ArrayEach 自身的错误，与 callback 内同名的 err
			// 形参无关。这里无法像 ObjectEach 路径那样在 callback 里中止迭代+
			// 透出错误，故 NotExist 在数组元素中被静默跳过。如需行为改变（让
			// NotExist 触发 panic），需要在闭包外捕获一个 *error 再 xpanic。
		})
		xpanic.WhenError(err)
		gen.Out()
		gen.P("}")
	}
}
