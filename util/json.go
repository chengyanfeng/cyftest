package utils

import (
	"encoding/json"

	"github.com/clbanning/mxj"
)
//1.json字符串转换P类型
func JsonDecode(b []byte) (p *P) {
	p = &P{}
	err := json.Unmarshal(b, p)
	if err != nil {
		Error("JsonDecode", string(b), err)
	}
	return
}
//2.interface类型转换成json格式的字符串
func JsonEncode(v interface{}) (r string) {
	b, err := json.Marshal(v)
	if err != nil {
		Error(err)
	}
	r = string(b)
	return
}
//3.判断是不是json
func IsJson(b []byte) bool {
	var j json.RawMessage
	return json.Unmarshal(b, &j) == nil
}
//4.解析json数组成为P类型
func JsonDecodeArray(b []byte) (p []P, e error) {
	p = []P{}
	e = json.Unmarshal(b, &p)
	return
}
//5.解析json为数组
func JsonDecodeArray_str(b []byte) (p []string, e error) {
	p = []string{}
	e = json.Unmarshal(b, &p)
	return
}
//6.解析json数组为P
func JsonDecodeArrays(b []byte) (p *[]P) {
	p = &[]P{}
	e := json.Unmarshal(b, p)
	if e != nil {
		Error(e)
	}
	return
}
//7.解析json为数组
func JsonDecodeStrings(s string) (r []string) {
	r = []string{}
	e := json.Unmarshal([]byte(s), &r)
	if e != nil {
		Error(e, s)
	}
	return
}
//8.解析多个参数为字符串
func JoinStr(val ...interface{}) (r string) {
	for _, v := range val {
		r += ToString(v)
	}
	return
}
//解析xml为字符串
func Xml2Json(src string) (s string, err error) {
	m, err := mxj.NewMapXml([]byte(src))
	return JsonEncode(m), err
}
