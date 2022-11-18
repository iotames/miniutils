package miniutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	neturl "net/url"
	"os"
	"strings"
	"time"
)

type HttpRequest struct {
	Client      *http.Client
	Request     *http.Request
	OnRequest   func(r *http.Request)
	Response    *http.Response
	BodyBytes   []byte
	Url         string
	RetryTimes  uint8
	postPayload string
}

// NewHttpRequest 构建http请求，默认为GET方法的HTTP请求。
func NewHttpRequest(url string) *HttpRequest {
	h := &HttpRequest{}
	h.Url = url
	h.Client = &http.Client{}
	h.Client.Transport = http.DefaultTransport
	h.setRequestGet()
	return h
}

// Do 执行之构建好的HTTP请求
func (h *HttpRequest) Do(callback func(h *HttpRequest)) error {
	err := h.checkRequest()
	if err != nil {
		return err
	}
	if h.OnRequest != nil {
		h.OnRequest(h.Request)
	}

	fmt.Printf("\n\n----RequestBegin---%s: %s -------------------------------------------\n\n", h.Request.Method, h.Request.URL)
	for key, value := range h.Request.Header {
		hvalue := ""
		for _, vv := range value {
			hvalue += "[" + vv + "]"
		}
		fmt.Println(key+": ", hvalue)
	}
	if h.Request.Method == "POST" {
		fmt.Printf("\n----POST---Payload: %s-----------\n", h.postPayload)
	}

	resp, err := h.Client.Do(h.Request)
	h.Response = resp
	if err != nil {
		fmt.Println("--Error---------RequestError--------", err)

		if h.RetryTimes > 0 {
			for i := 0; i < int(h.RetryTimes); i++ {
				fmt.Println("-------RetryRequest------", i+1)
				resp, err = h.Client.Do(h.Request)
				if err == nil {
					fmt.Printf("-------RetryRequestOK--times: %d----url: %s---\n\n", i+1, h.Request.URL)
					h.Response = resp
					break
				}
			}
		}
		if err != nil {
			fmt.Printf("------RetryRequestFail---Error: %v-----------\n", err)
			// log.Fatal(err)
			return err
		}
	}
	fmt.Printf("\n\n----RequestEnd---(%s)----------------------\n\n", resp.Status)

	defer resp.Body.Close()

	if callback == nil {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			// log.Fatal(err)
			log.Printf("---Error Happened. ioutil.ReadAll:%v-------\n", err)
			return err
		}
		h.BodyBytes = bodyBytes
	} else {
		callback(h)
	}
	return nil
}

// ParseJsonResponse 解码HTTP响应的json内容。 when response body is a json string
// param model: a struct point
func (h *HttpRequest) ParseJsonResponse(model interface{}) error {
	err := json.Unmarshal(h.BodyBytes, model)
	if err != nil {
		log.Printf("---Error Happened. ParseJsonResponse:%v-------\n", err)
		return err
	}
	return nil
}

// SetRequestPost 构建POST方法请求，设置POST数据内容。
func (h *HttpRequest) SetRequestPost(data interface{}) error {
	pdata, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return h.SetRequestPostByString(string(pdata))
}

func (h *HttpRequest) setRequestGet() error {
	req, err := http.NewRequest("GET", h.Url, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	h.Request = req
	return nil
}

// SetRequestPostByString 构建POST方法请求，设置POST数据字符串。
func (h *HttpRequest) SetRequestPostByString(data string) error {
	// return http.NewRequestWithContext(context.Background(), method, url, body)
	req, err := http.NewRequest("POST", h.Url, strings.NewReader(data))
	if err != nil {
		// panic(req)
		log.Fatal(err)
		return err
	}
	h.Request = req
	h.postPayload = data
	return nil
}

func (h *HttpRequest) checkRequest() error {
	if h.Request == nil {
		return errors.New("request can not be nil. Please Set Request before")
	}
	return nil
}

// SetRequestHeader 设置HTTP请求头
func (h *HttpRequest) SetRequestHeader(key string, value string) error {
	err := h.checkRequest()
	if err != nil {
		log.Fatal(err)
		return err
	}
	h.Request.Header.Set(key, value)
	return nil
}

// 设置HTTP代理。 例: SetProxy("http://127.0.0.1:1080")
func (h *HttpRequest) SetProxy(proxyUrl string) {
	fmt.Println("------SetProxy------", proxyUrl)
	h.Client.Transport.(*http.Transport).Proxy = func(_ *http.Request) (*neturl.URL, error) {
		return neturl.Parse(proxyUrl)
	}
}

// SetTimeout 设置请求超时时间
func (h *HttpRequest) SetTimeout(seconds uint8) {
	h.Client.Timeout = time.Duration(seconds) * time.Second
}

// Download 下载文件. example:
// NewHttpRequest("https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png").Download("runtime/baidu.png")
func (h *HttpRequest) Download(file string) error {
	err := h.checkRequest()
	if err != nil {
		log.Fatal(err)
		return err
	}
	f, err := os.Create(file)
	if err != nil {
		log.Printf("---Error Happened. os.Create:%v-------\n", err)
		// panic(err)
		return err
	}
	fmt.Println("------Download To: " + file)
	err = h.Do(func(h *HttpRequest) {
		_, err := io.Copy(f, h.Response.Body)
		if err != nil {
			// panic(err)
			log.Printf("---Error Happened. io.Copy:%v-------\n", err)
		}
	})
	return err
}
