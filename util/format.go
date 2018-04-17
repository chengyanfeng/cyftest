package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"
	"reflect"
)

//1.转换成int类型
func ToInt(s interface{}, default_v ...int) int {
	i, e := strconv.Atoi(ToString(s))
	if e != nil && len(default_v) > 0 {
		return default_v[0]
	}
	return i
}
//2.转换成int64
func ToInt64(s interface{}, default_v ...int64) int64 {
	switch s.(type) {
	case int64:
		return s.(int64)
	case int:
		return int64(s.(int))
	case float64:
		return int64(s.(float64))
	}
	i64, e := strconv.ParseInt(ToString(s), 10, 64)
	if e != nil && len(default_v) > 0 {
		return default_v[0]
	}
	return i64
}
//3.转换成float类型
func ToFloat(s interface{}, default_v ...float64) float64 {
	f64, e := strconv.ParseFloat(ToString(s), 64)
	if e != nil && len(default_v) > 0 {
		return default_v[0]
	}
	return f64
}
//4.转换成string类型
func ToString(v interface{}, def ...string) string {
	if v != nil {
		switch v.(type) {
		case bson.ObjectId:
			return v.(bson.ObjectId).Hex()
		case []byte:
			return string(v.([]byte))
		case *P, P:
			var p P
			switch v.(type) {
			case *P:
				if v.(*P) != nil {
					p = *v.(*P)
				}
			case P:
				p = v.(P)
			}
			var keys []string
			for k := range p {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			r := "P{"
			for _, k := range keys {
				r = JoinStr(r, k, ":", p[k], " ")
			}
			r = JoinStr(r, "}")
			return r
		case map[string]interface{}, []P, []interface{}:
			return JsonEncode(v)
		case int64:
			return strconv.FormatInt(v.(int64), 10)
		case []string:
			s := ""
			for _, j := range v.([]string) {
				s = JoinStr(s, ",", j)
			}
			if len(s) > 0 {
				s = s[1:]
			}
			return s
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	if len(def) > 0 {
		return def[0]
	} else {
		return ""
	}
}
//5.转换成p类型
func ToP(v interface{}) P {
	if v != nil {
		switch v.(type) {
		case P:
			return v.(P)
		case *P:
			return *v.(*P)
		case string:
			return *JsonDecode([]byte(v.(string)))
		case map[string]interface{}:
			return v.(map[string]interface{})
		default:
			Error("ToP fail", ToString(v))
		}
	}
	return P{}
}
//6.转换成string数组
func ToStrings(v interface{}) []string {
	strs := []string{}
	if v != nil {
		switch v.(type) {
		case []interface{}:
			for _, i := range v.([]interface{}) {
				strs = append(strs, ToString(i))
			}
		case []string:
			for _, i := range v.([]string) {
				strs = append(strs, i)
			}
		case string, interface{}:
			strs = append(strs, ToString(v))
		}
	}
	return strs
}
//6.string 类型转换成string数组
func ToFields(s string, div string) (r []string) {
	s = Replace(s, []string{`""`}, "")
	tmp := strings.Split(s, div)
	r = []string{}
	state := ""
	seg := ""
	for _, v := range tmp {
		v = Trim(v)
		if len(v) > 1 && StartsWith(v, `"`) && !EndsWith(v, `"`) {
			state = `s`
		} else if !StartsWith(v, `"`) && EndsWith(v, `"`) {
			state = "e"
		} else if state == `s` && v == `"` {
			state = "e"
		}
		if state == "s" {
			seg += "," + v
			seg = TransFunc(seg)
		} else if state == "e" {
			seg += "," + v
			if len(seg) > 1 {
				seg = seg[1:]
			}
			seg = TransFunc(seg)
			r = append(r, seg)
			seg = ""
			state = ""
		} else {
			v = TransFunc(v)
			r = append(r, v)
		}
	}
	return
}
//8.转换成oid interface类型
func ToOid(id interface{}) (oid bson.ObjectId) {
	s := ToString(id)
	if bson.IsObjectIdHex(s) {
		oid = bson.ObjectIdHex(s)
	}
	return
}

func ToOids(ids interface{}) (oids []bson.ObjectId) {
	oids = []bson.ObjectId{}
	switch ids.(type) {
	case []string:
		for _, id := range ids.([]string) {
			if IsOid(id) {
				oids = append(oids, ToOid(id))
			}
		}
	case []interface{}:
		for _, id := range ids.([]interface{}) {
			if IsOid(ToString(id)) {
				oids = append(oids, ToOid(ToString(id)))
			}
		}
	case []bson.ObjectId:
		oids = ids.([]bson.ObjectId)
	}
	return
}
//9.结构体转换成P类型，也就是map
func StructToMap(obj interface{}) P {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[strings.ToLower(t.Field(i).Name)] = v.Field(i).Interface()
	}
	return data
}

//func StructToMapArray(obj interface{}) (datas []P) {
// 	if objs, ok := obj.([]*models.DhBase); ok {
// 		for _, value := range objs {
// 			datas = append(datas, StructToMap(*value))
// 		}
// 	}
// 	return
// }