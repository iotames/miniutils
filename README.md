<div align="center">
   <span style="font-size:100px">🧰</span>
  <br><a href="https://github.com/iotames/miniutils">Github</a> | <a href="https://gitee.com/catmes/miniutils">Gitee</a>
  <br>入门级，简单，易用的Go小工具<br>助你成功转职Golang工程师！
</div>


## 简介

[![GoDoc](https://badgen.net/badge/Go/referenct)](https://pkg.go.dev/github.com/iotames/miniutils)
[![License](https://badgen.net/badge/License/MIT/green)](https://github.com/iotames/miniutils/blob/main/LICENSE)

Go实用小工具: 封装了高效开发的常用工具、函数集，简单易用。


## 快速开始

```
# 创建本地项目 myproject
go mod init myproject
# 新建入口文件 main.go
vim main.go
```

入口文件 `main.go`

```
package main

import (
	"fmt"
	"github.com/iotames/miniutils"
)

func main() {
    strfind := miniutils.NewStrfind("https://d.168.com/offer/356789.html")
	dofind := strfind.SetRegexp(`offer/(\d+)\.html`).DoFind()
	fmt.Println(dofind.GetOne(false)) // "356789"
	fmt.Println(dofind.GetOne(true)) // "offer/356789.html"
}
```

```
# 更新依赖
go mod tidy
# 运行
go run .
```

## 示例

[测试样例](/utils_test.go)

### JWT 工具

`JWT`: 全称JSON Web Token，互联网API通讯接口身份验证的行业标准。通过JWT字符串的解密和验签，进行用户身份认证。参见: https://jwt.io/

```
package main

import (
	"fmt"
	"log"
	"time"
	"github.com/iotames/miniutils"
)

func main() {
	secret := miniutils.GetRandString(32) // 设置JWT签名密钥
	jwt := miniutils.NewJwt(secret) // 初始化JWT小工具
	jwtInfo := map[string]interface{}{"id": 1519512704946016256, "name": "Harvey", "age": 16, "mobile": "15988888888"}
	// 设置原始数据jwtInfo，有效期3600秒，创建JWT字符串tokenStr
	tokenStr, err := jwt.Create(jwtInfo, time.Second*3600)
	if err != nil {
		fmt.Printf("jwt.Create error: %v", err)
        return
	}
	log.Println("create JWT:", tokenStr)
	// 解码 JWT 字符串. 返回 map[string]interface{} 格式的数据。
	info, err := miniutils.NewJwt("").Decode(tokenStr)
	if err != nil {
		fmt.Printf("jwt.Decode error: %v", err)
        return
	}
	log.Println("jwt Decode:", info)
	// 解码 JWT 字符串并验签，验证有效期。 返回 map[string]interface{} 格式的数据。
	claims, err := jwt.Parse(tokenStr)
	if err != nil {
		fmt.Printf("jwt.Parse error: %v", err)
        return
	}
	log.Println("jwt Parse:", claims)
    
    _, err = jwt.Parse(tokenStr + "sign error")
	if err == miniutils.ErrTokenSign {
		fmt.Printf("JWT 签名错误")
	}
}
```

### 日志记录

```
	logger := miniutils.GetLogger("")
	logger.Debug("first log 11111")
	logger.Info("second log 22222")
	logger = miniutils.NewLogger("runtime/mylogs")
	logger.Debug("my logs 2333")
	logger.CloseLogFile()
```

### 字符串提取工具

```
    strfind := miniutils.NewStrfind("https://d.168.com/offer/356789.html")
	dofind := strfind.SetRegexp(`offer/(\d+)\.html`).DoFind()
	fmt.Println(dofind.GetOne(false)) // "356789"
	fmt.Println(dofind.GetOne(true)) // "offer/356789.html"
```

### HTTP请求工具

```
	// 构建HTTP请求(默认GET方法)
    req := miniutils.NewHttpRequest("https://httpbin.org/get")
	// 设置HTTP请求头
	req.SetRequestHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.81 Safari/533.33")
	// 执行HTTP请求
	err := req.Do(nil)
	if err != nil {
		log.Println(err)
	}
	// 打印响应内容消息体
	log.Println(string(req.BodyBytes))

	req = miniutils.NewHttpRequest("https://httpbin.org/post")
	// 构建POST请求
	req.SetRequestPostByString("hello=word&some=2333")
	// 执行HTTP请求
	req.SetRequestHeader("xkey", "secretttkeyyy")
	err = req.Do(nil)
	if err != nil {
		log.Printf("request post do err (%v) \n", err)
	}
	// 打印HTTP响应对象
	log.Println(*req.Response)

	// 下载图片到本地
	miniutils.NewHttpRequest("https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png").Download("runtime/baidu.png")
```