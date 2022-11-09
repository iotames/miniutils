## Go实用小工具

为提高开发效率封装了一些实用工具、函数，尽量不依赖第三方库。
简单易用，拿来当Golang入门练习也是不错的选择。


## 示例

### JWT 工具

`JWT`: https://jwt.io/

```
import (
	"fmt"
	"log"
	"time"
)
func main() {
	secret := GetRandString(32) // 设置JWT签名密钥
	jwt := NewJwt(secret) // 初始化JWT小工具
	jwtInfo := map[string]interface{}{"id": 1519512704946016256, "name": "Harvey", "age": 16, "mobile": "15988888888"}
	tokenStr, err := jwt.Create(jwtInfo, time.Second*3600) // 设置原始数据jwtInfo，有效期3600秒，创建JWT字符串tokenStr
	if err != nil {
		fmt.Println("jwt.Create error: %v", err)
        return
	}
	log.Println("create JWT:", tokenStr)
	info, err := NewJwt("").Decode(tokenStr) // 解码 JWT 字符串. 返回 map[string]interface{} 格式的数据。
	if err != nil {
		fmt.Println("jwt.Decode error: %v", err)
        return
	}
	log.Println("jwt Decode:", info)

	claims, err := jwt.Parse(tokenStr) // 解码 JWT 字符串并验签，验证有效期。 返回 map[string]interface{} 格式的数据。
	if err != nil {
		fmt.Println("jwt.Parse error: %v", err)
        return
	}
	log.Println("jwt Parse:", claims)
    
    _, err = jwt.Parse(tokenStr + "sign error")
	if err == ErrTokenSign {
		fmt.Println("JWT 签名错误")
	}
}
```