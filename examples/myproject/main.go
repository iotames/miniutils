package main

import (
	"fmt"
	"log"

	"github.com/iotames/miniutils"
)

func main() {
	// 构建HTTP请求(默认GET方法)
	req := miniutils.NewHttpRequest("https://www.baidu.com")
	// 设置HTTP请求头
	req.SetRequestHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.81 Safari/533.33")
	// 执行HTTP请求
	err := req.Do(nil)
	if err != nil {
		log.Println(err)
		return
	}
	body := string(req.BodyBytes)
	// 提取网页中，以`<span class="title-content-title">`字符串开头，`</span>` 字符串结尾的所有字符串片段。
	strfind := miniutils.NewStrfind(body).SetRegexpBeginEnd(`<span class="title-content-title">`, `</span>`)
	hots := strfind.DoFind().GetAll(false)
	hotsMsg := "百度热搜:\n"
	// logger := miniutils.GetLogger("")
	for i, v := range hots {
		hotTitle := fmt.Sprintf("[%d]----->[%s]\n", i, v)
		// logger.Debug(hotTitle)
		hotsMsg += hotTitle
	}
	log.Println(hotsMsg)
}
