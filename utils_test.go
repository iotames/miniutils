package miniutils

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestJwt(t *testing.T) {
	secret := GetRandString(32) // 创建32位的签名密钥
	jwt := NewJwt(secret)
	jwtInfo := map[string]interface{}{"id": 1519512704946016256, "name": "Harvey", "age": 16, "mobile": "15900000001"}

	// 创建JWT字符串
	tokenStr, err := jwt.Create(jwtInfo, time.Second*3600)
	if err != nil {
		t.Errorf("jwt.Create error: %v", err)
	}
	log.Println("create JWT:", tokenStr)

	// 解码JWT
	info, err := NewJwt("").Decode(tokenStr)
	if err != nil {
		t.Errorf("jwt.Decode error: %v", err)
	}
	log.Println("jwt Decode:", info)

	// 解码并验签，验证有效期
	claims, err := jwt.Parse(tokenStr)
	if err != nil {
		t.Errorf("jwt.Parse error: %v", err)
	}
	log.Println("jwt Parse:", claims)

	// 核对数据正确性
	if fmt.Sprintf("%v", jwtInfo) != fmt.Sprintf("%v", info) {
		t.Errorf("jwt.Decode error")
	}
	if fmt.Sprintf("%v", jwtInfo) != fmt.Sprintf("%v", claims) {
		t.Errorf("jwt.Parse error")
	}
	log.Println("jwt info:", jwtInfo)

	// 构建验签失败的JWT
	_, err = jwt.Parse(tokenStr + "sign error")
	if err != ErrTokenSign {
		t.Errorf("jwt.Parse error")
	}

	// 构建超过有效期的JWT
	exp := jwtInfo["exp"].(int64)
	jwtInfo["exp"] = exp - 3601
	tokenStr, _ = jwt.Create(jwtInfo, time.Second*3600)
	_, err = jwt.Parse(tokenStr)
	if err != ErrTokenExpired {
		t.Errorf("jwt.Parse error")
	}
}