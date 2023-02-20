package xgen

import (
	"fmt"
	"testing"
)

var testCode = []byte(`
/*
 * I am C style comments
 */
 // Code generated by ProtoKit. DO NOT EDIT. (client)
 // from sheet@ErrorCode|Design
 syntax = "proto3";
 package rawdata;
 option go_package = "example/gen/golang/rawdata";
 option csharp_namespace = "gen.rawdata";
 // annotation@annotation_type(type="rawdata")
 // annotation@rawdata(row="ErrorCode",table="ErrorCodeMap",table_wrapper="ErrorCodeConf",map="true",map_key_type="int32",data="ErrorCodeConf")
 // annotation@ab(table_ab="ErrorCodeMapAB",table_ab_patch="ErrorCodeMapABPatch",table_ab_value="ErrorCodeMapABValue",ab_patch="false")
 // annotation@filter(just_server="false",just_client="false")
 // annotation@ErrorCode(id="id")


 message ErrorCode{
   // 错误信息id
   // annotation@field_id(id="true")
   int32 id = 1;
   // 错误信息key
   string Key = 2;
 }
 message ErrorCodeMapABValue{
   map<int32, ErrorCode> ErrorCodeMap = 1;
 }
 message ErrorCodeConf{
   map<int32, ErrorCode> ErrorCodeMap = 1;
   map<string, ErrorCodeMapABValue> ErrorCodeMapAB = 2;
   map<string, ErrorCodeMapABValue> ErrorCodeMapABPatch = 3;
 }
`)

func TestRmoveCAndCppCommentAndBlanklines(t *testing.T) {
	fmt.Println(string(RmoveCAndCppCommentAndBlanklines(testCode)))
}
