package miniutils

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
	// "encoding/hex"
)

var once sync.Once
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyz")

func GetRandString(l int) string {
	once.Do(func() {
		rand.Seed(time.Now().UnixNano())
	})
	b := make([]rune, l)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Md5 get the MD5 hash.
func Md5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
	// m := md5.New()
	// m.Write([]byte(s))
	// return hex.EncodeToString(m.Sum(nil))
}

// GetSha256 get SHA256 hash. The len of SHA256 value is 64.
func GetSha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashBytes := h.Sum(nil)
	return fmt.Sprintf("%x", hashBytes)
}

// GetSha256BySecret get SHA256 hash by Key
func GetSha256BySecret(secret string, keyBytes []byte) []byte {
	hash := crypto.SHA256
	if !hash.Available() {
		return make([]byte, 0)
	}
	hasher := hmac.New(hash.New, keyBytes)
	hasher.Write([]byte(secret))
	return hasher.Sum(nil)
}

// GetFileSha256 get SHA256 hash of file.
func GetFileSha256(path string) (hash string, err error) {
	var f *os.File
	f, err = os.Open(path)
	if err != nil {
		err = fmt.Errorf("file open error:%w", err)
		return
	}
	defer f.Close()
	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		err = fmt.Errorf("io.Copy error:%w", err)
		return
	}
	hashBytes := h.Sum(nil)
	hash = fmt.Sprintf("%x", hashBytes)
	return
}

// Base64UrlEncode Encode JWT specific base64url encoding with padding stripped
func Base64UrlEncode(seg []byte) string {
	return base64.RawURLEncoding.EncodeToString(seg)
}

// Base64UrlDecode Decode JWT specific base64url encoding with padding stripped
func Base64UrlDecode(seg string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(seg)
}
