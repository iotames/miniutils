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
	log.Printf("jwt info:%+v(%T)", jwtInfo, jwtInfo["id"]) // int

	// 解码JWT
	info, err := NewJwt("").Decode(tokenStr)
	if err != nil {
		t.Errorf("jwt.Decode error: %v", err)
	}
	log.Printf("jwt Decode:%+v(%T)", info, info["id"]) // json.Number
	// 核对数据正确性
	// if !reflect.DeepEqual(jwtInfo, info) {
	if fmt.Sprintf("%v", jwtInfo) != fmt.Sprintf("%v", info) {
		t.Errorf("jwt.Decode error")
	}

	// 解码JWT并验签，验证有效期
	claims, err := jwt.Parse(tokenStr)
	if err != nil {
		t.Errorf("jwt.Parse error: %v", err)
	}
	log.Printf("jwt Parse:%+v(%T)", claims, claims["id"]) // json.Number
	// 核对数据正确性
	// if !reflect.DeepEqual(jwtInfo, claims) {
	if fmt.Sprintf("%v", jwtInfo) != fmt.Sprintf("%v", claims) {
		t.Errorf("jwt.Parse error")
	}

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

func TestStrfind(t *testing.T) {
	strfind := NewStrfind("https://d.168.com/offer/356789.html")
	strf := strfind.SetRegexp(`offer/(\d+)\.html`).DoFind()
	offerCode := strf.GetOne(false)
	if offerCode != "356789" {
		t.Errorf("strfind error")
	}
	allstr := strf.GetOne(true)
	if allstr != "offer/356789.html" {
		t.Errorf("strfind error")
	}
	allF := strf.GetAll(false)
	log.Println(allF)
	allT := strf.GetAll(true)
	log.Println(allT)
}

func TestLogger(t *testing.T) {
	logger := GetLogger("")
	logger.Debug("first log 11111")
	logger.Info("second log 22222")
	logger = NewLogger("runtime/loggg1")
	logger.Debug("my logs 123")
	logger.CloseLogFile()
}

func TestOsfile(t *testing.T) {
	err := CopyFile("README.md", "README.md.copy")
	if err != nil {
		t.Errorf("CopyFile err:%v", err)
	}
	err = CopyDir("hello123", "hello321")
	if err != nil {
		t.Errorf("CopyDir err:%v", err)
	}
}

func TestHttpRequest(t *testing.T) {
	req := NewHttpRequest("https://httpbin.org/get")
	req.SetRequestHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.81 Safari/533.33")
	err := req.Do(nil)
	if err != nil {
		t.Errorf("request get do err %v", err)
	}
	log.Println(string(req.BodyBytes))

	req = NewHttpRequest("https://httpbin.org/post")
	req.SetRequestPostByString("hello=word&some=2333")
	req.SetRequestHeader("xkey", "secretttkeyyy")
	err = req.Do(nil)
	if err != nil {
		t.Errorf("request post do err %v", err)
	}
	log.Println(*req.Response)

	req = NewHttpRequest("https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png")
	req.SetRequestHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.81 Safari/533.33")
	req.Download("runtime/baidu.png")
}

func TestArray(t *testing.T) {
	a1 := []int64{12, 16, 33, 55}
	b1 := []string{"aa", "bb", "cc"}
	if GetIndexOf(33, a1) != 2 {
		t.Errorf("GetIndexOf Error")
	}
	if GetIndexOf[int64](999, a1) != -1 {
		t.Errorf("GetIndexOf Error")
	}
	if GetIndexOf("bb", b1) != 1 {
		t.Errorf("GetIndexOf Error")
	}
	if GetIndexOf[string]("hello", b1) != -1 {
		t.Errorf("GetIndexOf Error")
	}
}

func TestNumber(t *testing.T) {
	if StrToInt("15,633,510", nil) != 15633510 {
		t.Errorf("StrToInt error")
	}
	if StrToInt("$ 156,335,10", []string{"$", ","}) != 15633510 {
		t.Errorf("StrToInt error")
	}
	price := GetPriceByText("$ 156,335,10.37")
	fmt.Printf("%.2f\n", price)
	if price != 15633510.37 {
		t.Errorf("StrToInt error")
	}
	price = GetPriceByText("£156,335,10.37")
	fmt.Printf("%.2f\n", price)
	if price != 15633510.37 {
		t.Errorf("StrToInt error")
	}
}
