package main

import (
	"fmt"
	"log"
	"time"

	"github.com/iotames/miniutils"
)

func main() {
	secret := miniutils.GetRandString(32) // 设置JWT签名密钥
	jwt := miniutils.NewJwt(secret)       // 初始化JWT小工具
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
