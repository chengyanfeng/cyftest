package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/astaxie/beego/cache"
	"github.com/henrylee2cn/mahonia"
	"gopkg.in/mgo.v2/bson"

	. "common.dh.cn/def"
)

var localCache cache.Cache
var CronAuth = &P{"Authorization": JoinStr("Basic ", Base64Encode([]byte("mrocker:mrocker")))}

func init() {
	c, err := cache.NewCache("memory", `{"interval":60}`)
	//c, err := cache.NewCache("file", `{"CachePath":"./dhcache","FileSuffix":".cache","DirectoryLevel":2,"EmbedExpiry":120}`)
	if err != nil {
		Error(err)
	} else {
		localCache = c
	}
}

func IsInt(s interface{}) bool {
	_, e := strconv.ParseInt(ToString(s), 10, 64)
	return e == nil
}

func IsFloat(s interface{}) bool {
	_, e := strconv.ParseFloat(ToString(s), 64)
	return e == nil
}

func Md5(s ...interface{}) (r string) {
	return Hash("md5", s...)
}

func Hash(algorithm string, s ...interface{}) (r string) {
	r = hex.EncodeToString(HashBytes(algorithm, s...))
	return
}

func HashBytes(algorithm string, s ...interface{}) (r []byte) {
	var h hash.Hash
	switch algorithm {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha2", "sha256":
		h = sha256.New()
	}
	for _, value := range s {
		switch value.(type) {
		case []byte:
			h.Write(value.([]byte))
		default:
			h.Write([]byte(ToString(value)))
		}
	}
	r = h.Sum(nil)
	return
}

func HashMac(algorithm string, key []byte, s ...interface{}) (r []byte) {
	var mac hash.Hash
	switch algorithm {
	case "md5":
		mac = hmac.New(md5.New, key)
	case "sha1":
		mac = hmac.New(sha1.New, key)
	case "sha2", "sha256":
		mac = hmac.New(sha256.New, key)
	}
	for _, value := range s {
		switch value.(type) {
		case []byte:
			mac.Write(value.([]byte))
		default:
			mac.Write([]byte(ToString(value)))
		}
	}
	r = mac.Sum(nil)
	return
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(s string) []byte {
	r, e := base64.StdEncoding.DecodeString(s)
	if e != nil {
		Error(e)
	}
	return r
}

func InArray(s string, a []string) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}

func StartsWith(s string, a ...string) bool {
	for _, x := range a {
		if strings.HasPrefix(s, x) {
			return true
		}
	}
	return false
}

func EndsWith(s string, a ...string) bool {
	for _, x := range a {
		if strings.HasSuffix(s, x) {
			return true
		}
	}
	return false
}

func Unset(p P, keys ...string) {
	for _, x := range keys {
		delete(p, x)
	}
}

func Rand(start int, end int) int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(end)
	if r < start {
		r = start + rand.Intn(end-start)
	}
	//time.Sleep(1 * time.Nanosecond)
	return r
}

func Replace(src string, find []string, r string) string {
	for _, v := range find {
		src = strings.Replace(src, v, r, -1)
	}
	return src
}

func Count(src string, find []string) (c int) {
	for _, v := range find {
		c += strings.Count(src, v)
	}
	return
}

func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}
	switch v.(type) {
	case P:
		return len(v.(P)) == 0
	case []interface{}:
		return len(v.([]interface{})) == 0
	case []P:
		return len(v.([]P)) == 0
	case *[]P:
		return len(*v.(*[]P)) == 0
	}
	return ToString(v) == ""
}

func Trim(str string) string {
	return strings.TrimSpace(str)
}

func Ip2Int(ip string) int64 {
	sec := strings.Split(ip, ".")
	if len(sec) == 4 {
		return int64(ToInt(sec[0]))<<24 + int64(ToInt(sec[1]))<<16 + int64(ToInt(sec[2]))<<8 + int64(ToInt(sec[3]))
	}
	return 0
}

func GetCronStr(sec int) (str string) {
	ss := sec % 60
	ii := sec / 60
	hh := sec / 3600
	if ii == 0 && hh == 0 {
		str = fmt.Sprintf("0/%v * * * * *", sec)
	} else if ii > 0 && hh == 0 {
		str = fmt.Sprintf("%v */%v * * * *", ss, ii)
	} else if hh > 0 {
		str = fmt.Sprintf("%v %v */%v * * *", ss, ii%60, hh)
	} else {
		str = "0/60 * * * * *"
	}
	return
}

func Gbk2Utf(str string) string {
	enc := mahonia.NewDecoder("gbk")
	return enc.ConvertString(str)
}

func RenderTpl(tpl string, data interface{}) string {
	var bb bytes.Buffer
	//t, err := template.ParseFiles(tpl)
	t, err := template.New(Md5(tpl)).Parse(tpl)
	if err != nil {
		Error(err)
	}
	t.Execute(&bb, data)
	return bb.String()
}

func AddInOid(oids *[]bson.ObjectId, nid bson.ObjectId) {
	for _, oid := range *oids {
		if oid.Hex() == nid.Hex() {
			return
		}
	}
	*oids = append(*oids, nid)
	return
}

// 缓存接口，存 S("key", value)，取 S("key")
func S(key string, p ...interface{}) (v interface{}) {
	md5 := Md5(key)
	if len(p) == 0 {
		return localCache.Get(md5)
	} else {
		if len(p) == 2 {
			var ttl int64
			switch p[1].(type) {
			case int:
				ttl = int64(p[1].(int))
			case int64:
				ttl = p[1].(int64)
			}
			localCache.Put(md5, p[0], time.Duration(ttl)*time.Second)
		} else if len(p) == 1 {
			localCache.Put(md5, p[0], time.Duration(CACHE_TTL_DEFAULT)*time.Second)
		}
		return p[0]
	}
}

func SDel(key string) error {
	return localCache.Delete(Md5(key))
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}

func RegSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:]
	return result
}

func IsCsvEnd(s string, half bool) (b bool) {
	if half {
		b = Count(s, []string{`"`})%2 == 1
	} else {
		b = Count(s, []string{`"`})%2 == 0
	}
	return
}

func TransFunc(o string) (n string) {
	if StartsWith(o, "to_date(") {
		o = Trim(Replace(o, []string{"to_date(", ")"}, ""))
		tmp := strings.Split(o, " as ")
		field := ""
		as := ""
		field = tmp[0]
		if len(tmp) > 1 {
			as = tmp[1]
		}
		tmp = strings.Split(field, ",")
		if len(tmp) > 1 {
			if !IsEmpty(as) {
				n = JoinStr(n, " as ", as)
			}
		}
	} else if StartsWith(o, `"`) && EndsWith(o, `"`) && Count(o, []string{","}) == 0 && len(o) > 1 {
		n = o[1 : len(o)-1]
	} else {
		n = o
	}
	return
}

func Exec(cmd string, exp ...int) (str string, e error) {
	osname := runtime.GOOS
	var r *exec.Cmd
	Info("Exec:", cmd)
	if osname == "windows" {
		r = exec.Command("cmd", "/c", cmd)
	} else {
		r = exec.Command("/bin/bash", "-c", cmd)
	}
	var buf bytes.Buffer
	r.Stdout = &buf
	r.Stderr = &buf
	r.Start()

	if len(exp) < 1 {
		exp = []int{60}
	}
	done := make(chan error)
	go func() { done <- r.Wait() }()

	timeout := time.After(time.Duration(exp[0]) * time.Second)
	select {
	case <-timeout:
		r.Process.Kill()
		str = "Command timed out"
		e = errors.New(str)
	case e = <-done:
		str = buf.String()
	}
	if e != nil {
		Error("Exec", str, e)
	}
	return
}

func Cwd() string {
	cwd, _ := os.Getwd()
	return cwd
}

func NewId() bson.ObjectId {
	return bson.NewObjectId()
}

func IsOid(id string) bool {
	return bson.IsObjectIdHex(id)
}

func IsArray(v interface{}) bool {
	if IsEmpty(v) {
		return false
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

func IsMapArray(v interface{}) bool {
	a, b := v.([]interface{})
	if b {
		for _, m := range a {
			switch m.(type) {
			case map[string]interface{}:
				return true
			default:
				return false
			}
		}
	}
	return false
}

func Invoke(any interface{}, name string, args ...interface{}) {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(any).MethodByName(name).Call(inputs)
}

func Cooper(funcs ...func()) {
	wg := new(sync.WaitGroup)
	for _, f := range funcs {
		wg.Add(1)
		go func(f1 func()) {
			defer wg.Done()
			f1()
		}(f)
	}
	wg.Wait()
}

func RunJs(js string, data ...string) (r string, e error) {
	cmd := "js " + js
	for _, v := range data {
		datafile := Md5(v)
		datafile = JoinStr("/data/upload/", datafile, ".js")
		if !FileExists(datafile) {
			Mkdir("/data/upload")
			WriteFile(datafile, []byte(v))
		}
		cmd += " " + datafile
	}
	r, e = Exec(cmd, 3)
	return
}

func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func GetRandomString(number int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < number; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func ReplaceWhere(p P) {
	wherebytes := []byte(ToString(p["where"]))
	where := []P{}
	if len(wherebytes) > 0 {
		where, _ = JsonDecodeArray(wherebytes)
	}
	hasUrlParam := false
	for k, v := range p {
		if StartsWith(k, "__") {
			hasUrlParam = true
			for _, wp := range where {
				if wp["v"] == k {
					wp["v"] = v
				}
			}
		}
	}
	if hasUrlParam {
		r := []P{}
		for _, wp := range where {
			if !StartsWith(ToString(wp["v"]), "__") {
				r = append(r, wp)
			}
		}
		p["where"] = JsonEncode(r)
	} else {
		p["where"] = JsonEncode(where)
	}
	Info("ReplaceWhere", p["where"])
}

func ParseTableHead(th interface{}) []P {
	tmp, _ := JsonDecodeArray([]byte(ToString(th)))
	for _, v := range tmp {
		if !IsEmpty(v["o"]) {
			v["o"] = strings.Replace(ToString(v["o"]), " ", "_", -1)
		}

		tp := ToString(v["type"])
		switch tp {
		case "number", "long", "float", "int", "numeric":
			v["type"] = "number"
		case "date", "datetime", "timestamp":
			v["type"] = "date"
		case "location":
			v["type"] = "location"
		default:
			v["type"] = "series"
		}
	}
	return tmp
}

func CopyToP(from, to P) {
	for k, v := range from {
		to[k] = v
	}
}

func ReplaceRegx(src string, regex []string, r string) string {
	for _, v := range regex {
		src = strings.Replace(src, v, r, -1)
		re := regexp.MustCompile(v)
		src = re.ReplaceAllString(src, r)
	}
	return src
}
