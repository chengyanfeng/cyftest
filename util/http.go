package utils

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"common.dh.cn/def"
)

func HttpGet(url string, header *P, param *P) (body string, e error) {
	r, err := HttpGetBytes(url, header, param)
	if err != nil {
		Error("HttpGet异常:", err.Error())
	}
	e = err
	body = string(r)
	return
}

func HttpGetBytes(url string, header *P, param *P) (body []byte, e error) {
	return HttpDo("GET", url, header, param)
}

func HttpPost(url string, header *P, param *P) (body string, err error) {
	r, e := HttpDo("POST", url, header, param)
	if e != nil {
		Error("HttpPost异常:", e.Error())
		body = e.Error()
		err = e
	} else {
		body = string(r)
	}
	return
}

func HttpDelete(url string, header *P, param *P) (body []byte, e error) {
	return HttpDo("DELETE", url, header, param)
}

func HttpDo(method string, httpurl string, header *P, param *P) (body []byte, err error) {
	client := &http.Client{Timeout: time.Duration(def.DEFAULT_HTTP_TIMEOUT)}
	var req *http.Request
	vs := url.Values{}
	if param != nil {
		for k, v := range *param {
			key := ToString(k)
			if IsMapArray(v) {
				vs.Set(key, JsonEncode(v))
			} else if IsArray(v) {
				a, _ := v.([]interface{})
				for i, iv := range a {
					if i == 0 {
						vs.Set(key, ToString(iv))
					} else {
						vs.Add(key, ToString(iv))
					}
				}
			} else {
				vs.Set(key, ToString(v))
			}
		}
	}
	method = strings.ToUpper(method)
	req, err = http.NewRequest(method, httpurl, strings.NewReader(vs.Encode()))
	if header != nil {
		for k, v := range *header {
			req.Header.Set(ToString(k), ToString(v))
		}
	}
	if method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	resp, err := client.Do(req)
	if err != nil {
		Error("HttpDo异常:", err.Error())
		return []byte(ToString(resp)), err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func HttpRequest(method string, httpurl string, param *P) (body []byte, cookies []*http.Cookie, err error) {
	client := &http.Client{Timeout: time.Duration(def.DEFAULT_HTTP_TIMEOUT)}
	var req *http.Request
	vs := url.Values{}
	if param != nil {
		for k, v := range *param {
			key := ToString(k)
			if IsMapArray(v) {
				vs.Set(key, JsonEncode(v))
			} else if IsArray(v) {
				a, _ := v.([]interface{})
				for i, iv := range a {
					if i == 0 {
						vs.Set(key, ToString(iv))
					} else {
						vs.Add(key, ToString(iv))
					}
				}
			} else {
				vs.Set(key, ToString(v))
			}
		}
	}
	method = strings.ToUpper(method)
	req, err = http.NewRequest(method, httpurl, strings.NewReader(vs.Encode()))
	resp, err := client.Do(req)
	if err != nil {
		Error("HttpRequest异常:", err.Error())
		return []byte(ToString(resp)), nil, err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	cookies = resp.Cookies()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func HttpPostBody(url string, header *P, body []byte) (string, error) {
	client := &http.Client{Timeout: def.DEFAULT_HTTP_TIMEOUT}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	if header != nil {
		for k, v := range *header {
			req.Header.Set(ToString(k), ToString(v))
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		Error("HttpPostBody异常:", err.Error())
		return ToString(resp), err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	b, err := ioutil.ReadAll(resp.Body)
	return string(b), err
}

func UrlEncoded(str string) (string, error) {
	str = strings.Replace(str, "%", "%25", -1)
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func Upload(url, file string) (body []byte, err error) {
	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add your file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	fw, err := w.CreateFormFile("bin", file)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	// Add the other fields
	if fw, err = w.CreateFormField("key"); err != nil {
		return
	}
	if _, err = fw.Write([]byte("KEY")); err != nil {
		return
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	// Now that you have a form, you can submit it to your handler.
	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		Error("Upload异常:", err.Error())
		return []byte(ToString(res)), err
	}
	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()
	body, err = ioutil.ReadAll(res.Body)
	return
}
