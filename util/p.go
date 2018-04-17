package utils

import (
	"reflect"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type P map[string]interface{}
//复制
func (p *P) Copy() P {
	pn := make(P)
	for k, v := range *p {
		pn[k] = v
	}
	return pn
}

func (p P) CopyFrom(from P) {
	for k, v := range from {
		if IsEmpty(p[k]) {
			p[k] = v
		}
	}
}

func (p *P) ToInt(s ...string) {
	for _, k := range s {
		v := ToString((*p)[k])
		if !IsEmpty(v) {
			(*p)[k] = ToInt(v)
		}
	}
}

func (p *P) ToOid(s ...string) {
	for _, k := range s {
		v := ToString((*p)[k])
		if !IsEmpty(v) {
			if !IsOid(v) {
				Unset(*p, k)
				continue
			}
			(*p)[k] = ToOid(v)
		}
	}
}

func (p *P) ToOids(s ...string) {
	for _, k := range s {
		v := ToStrings((*p)[k])
		if !IsEmpty(v) && len(v) > 0 {
			(*p)[k] = ToOids(v)
		} else {
			Unset(*p, k)
		}
	}
}

func (p *P) Like(s ...string) {
	for _, k := range s {
		v := ToString((*p)[k])
		if !IsEmpty(v) {
			(*p)[k] = &bson.RegEx{Pattern: v, Options: "i"}
		}
	}
}

func (p *P) ToP(s ...string) (r P) {
	for _, k := range s {
		v := ToString((*p)[k])
		r = *JsonDecode([]byte(v))
		(*p)[k] = r
	}
	return
}

func (p *P) Get(k string, def interface{}) interface{} {
	r := (*p)[k]
	if r == nil {
		r = def
	}
	return r
}

func SetKv(p P, k string, v []string) {
	if len(v) == 1 {
		if len(v[0]) > 0 {
			p[k] = v[0]
		}
	} else {
		p[k] = v
	}
}

func ModelToP(o interface{}) P {
	info := P{}
	if o == nil {
		return info
	} else {
		s := reflect.ValueOf(o).Elem()
		for i := 0; i < s.NumField(); i++ {
			f := s.Type().Field(i)
			key := f.Tag.Get("json")
			if key == "" || key == "-" {
				continue
			}
			value := s.Field(i).Interface()
			if key == "update_time" {
				//value = ToBeiJingTime(value.(time.Time))
				value = value.(time.Time).Format("2006-01-02 15:04:05")
			}
			info[key] = value
		}
		return info
	}
}

func ModelToArrayP(o []*interface{}) []P {
	var array = []P{}
	for _, v := range o {
		array = append(array, ModelToP(v))
	}
	return array
}
