package miniutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrTokenFormat  = errors.New("token is not a JWT")
	ErrTokenExpired = errors.New("token is expired")
	ErrTokenSign    = errors.New("token sign error")
	ErrTokenExp     = errors.New("token lost field: exp")
)

func GetJwtBySecret(keyBytes []byte, bodyInfo map[string]interface{}) (string, error) {
	headerInfo := map[string]interface{}{"typ": "JWT", "alg": "HS256"}
	var sig, sstr string
	var err error
	if sstr, err = toJwtString(headerInfo, bodyInfo); err != nil {
		return "", err
	}
	sig = Base64UrlEncode(GetSha256BySecret(sstr, keyBytes))
	return strings.Join([]string{sstr, sig}, "."), nil
}

func toJwtString(headerInfo, bodyInfo map[string]interface{}) (string, error) {
	var err error
	parts := make([]string, 2)
	for i := range parts {
		var jsonValue []byte
		if i == 0 {
			jsonValue, err = json.Marshal(headerInfo)
		} else {
			jsonValue, err = json.Marshal(bodyInfo)
			// log.Println("----jwt--toJsonString---:", string(jsonValue))
		}
		if err != nil {
			return "", err
		}

		parts[i] = Base64UrlEncode(jsonValue)
	}

	return strings.Join(parts, "."), nil
}

type JsonWebToken struct {
	secret, TokenString string
	claims              map[string]interface{}
}

// NewJwt init JsonWebToken by secret string
func NewJwt(secret string) *JsonWebToken {
	return &JsonWebToken{secret: secret}
}

// Create JsonWebToken string
// Create(map[string]interface{}{"id": 123456789, "username": "Harvey"}, time.Second*time.Duration(3600))
func (j *JsonWebToken) Create(claims map[string]interface{}, expiresin time.Duration) (token string, err error) {
	_, ok := claims["exp"]
	if !ok {
		claims["exp"] = time.Now().Add(expiresin).Unix()
	}
	token, err = GetJwtBySecret([]byte(j.secret), claims)
	if err != nil {
		err = fmt.Errorf("GetJwtBySecret error: %w", err)
		return
	}
	j.TokenString = token
	j.claims = claims
	return
}

// JsonDecodeUseNumber 解析带数字的JSON
func JsonDecodeUseNumber(infoBytes []byte, result interface{}) error {
	// err = json.Unmarshal(infoBytes, result) 时间戳 int64 转json会变 float64
	// 未设置UseNumber长整型会丢失精度
	decoder := json.NewDecoder(bytes.NewReader(infoBytes))
	decoder.UseNumber()
	// fmt.Printf("----JsonDecodeUseNumber--(%p)-(%p)-----\n", result, &result)
	return decoder.Decode(result)
}

// Decode 解码JWT字符串。reads the JsonWebToken string. Return the JWT decoded data.
func (j *JsonWebToken) Decode(jwtStr string) (result map[string]interface{}, err error) {
	tokenSplit := strings.Split(jwtStr, ".")
	if len(tokenSplit) != 3 {
		err = ErrTokenFormat
		return
	}
	var infoBytes []byte
	infoBytes, err = Base64UrlDecode(tokenSplit[1])
	if err != nil {
		err = fmt.Errorf("Base64UrlDecode error: %w", err)
		return
	}
	// fmt.Printf("----Decode1--(%p)-(%p)--\n", result, &result)
	// result = make(map[string]interface{})
	// fmt.Printf("----Decode2---(%p)-(%p)--\n", result, &result)
	err = JsonDecodeUseNumber(infoBytes, &result) // &result 传递非空指针. 不加取址符&导致空指针错误: json: Unmarshal(non-pointer map[string]interface {})
	if err != nil {
		err = fmt.Errorf("JsonDecodeUseNumber error: %w", err)
	}
	j.claims = result
	return
}

// Parse 解码JWT字符串，并验证其有效性。reads the JsonWebToken string. Check the JWT decoded data and return.
func (j *JsonWebToken) Parse(jwtStr string) (result map[string]interface{}, err error) {
	if j.claims != nil {
		result = j.claims
	} else {
		result, err = j.Decode(jwtStr)
		if err != nil {
			return
		}
	}
	exp, ok := result["exp"]
	if ok {
		expiredAt, _ := exp.(json.Number).Int64()
		if expiredAt < time.Now().Unix() {
			err = ErrTokenExpired
			return
		}
	} else {
		err = ErrTokenExp
		return
	}
	var okToken string
	okToken, err = GetJwtBySecret([]byte(j.secret), result)
	if err != nil {
		err = fmt.Errorf("GetJwtBySecret error: %w", err)
		return
	}
	if okToken != jwtStr {
		err = ErrTokenSign
		return
	}
	return
}
