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
	return
}

func JsonDecodeUseNumber(infoBytes []byte, result interface{}) error {
	// 未设置UseNumber长整型会丢失精度
	decoder := json.NewDecoder(bytes.NewReader(infoBytes))
	decoder.UseNumber()
	return decoder.Decode(result)
}

// Decode reads the JsonWebToken string. Return the JWT decoded data.
func (j *JsonWebToken) Decode(jwtStr string) (segInfo map[string]interface{}, err error) {
	segTokens := strings.Split(jwtStr, ".")
	if len(segTokens) != 3 {
		err = ErrTokenFormat
		return
	}
	var infoBytes []byte
	infoBytes, err = Base64UrlDecode(segTokens[1])
	if err != nil {
		err = fmt.Errorf("Base64UrlDecode error: %w", err)
		return
	}
	err = JsonDecodeUseNumber(infoBytes, &segInfo)
	if err != nil {
		err = fmt.Errorf("JsonDecodeUseNumber error: %w", err)
	}
	return
}

// Decode reads the JsonWebToken string. Check the JWT decoded data and return.
func (j *JsonWebToken) Parse(tokenStr string) (result map[string]interface{}, err error) {
	tokenSplit := strings.Split(tokenStr, `.`)
	if len(tokenSplit) != 3 {
		err = ErrTokenFormat
		return
	}

	bodyInfoBytes, err := Base64UrlDecode(tokenSplit[1])
	if err != nil {
		err = fmt.Errorf("Base64UrlDecode error: %w", err)
		return
	}
	bodyInfo := map[string]interface{}{}
	// err = json.Unmarshal(bodyInfoBytes, &bodyInfo) 时间戳 int64 转json会变 float64
	err = JsonDecodeUseNumber(bodyInfoBytes, &bodyInfo)
	if err != nil {
		err = fmt.Errorf("JsonDecodeUseNumber error: %w", err)
		return
	}

	exp, ok := bodyInfo["exp"]
	if ok {
		expiredAt, _ := exp.(json.Number).Int64()
		if expiredAt < time.Now().Unix() {
			err = ErrTokenExpired
			return
		}
	}
	var okToken string
	okToken, err = GetJwtBySecret([]byte(j.secret), bodyInfo)
	if err != nil {
		err = fmt.Errorf("GetJwtBySecret error: %w", err)
		return
	}

	if okToken != tokenStr {
		err = ErrTokenSign
		return
	}
	result = bodyInfo
	err = nil
	return
}
