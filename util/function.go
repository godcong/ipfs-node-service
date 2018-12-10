package util

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

/*CustomHeader xml header*/
const CustomHeader = `<?xml version="1.0" encoding="UTF-8" standalone="no"?>`

/*CDATA xml cdata defines */
type CDATA struct {
	XMLName xml.Name
	Value   string `xml:",cdata"`
}

/* error types */
var (
	ErrorSignType  = errors.New("sign type error")
	ErrorParameter = errors.New("JsonApiParameters() check error")
	ErrorToken     = errors.New("EditAddressParameters() token is nil")
)

/*RandomKind RandomKind */
type RandomKind int

/*random kinds */
const (
	RandomNum      RandomKind = iota // 纯数字
	RandomLower                      // 小写字母
	RandomUpper                      // 大写字母
	RandomLowerNum                   // 数字、小写字母
	RandomUpperNum                   // 数字、大写字母
	RandomAll                        // 数字、大小写字母
)

/*RandomString defines */
var (
	RandomString = map[RandomKind]string{
		RandomNum:      "0123456789",
		RandomLower:    "abcdefghijklmnopqrstuvwxyz",
		RandomUpper:    "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		RandomLowerNum: "0123456789abcdefghijklmnopqrstuvwxyz",
		RandomUpperNum: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		RandomAll:      "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	}
)

/*ParseNumber parse interface to number */
func ParseNumber(v interface{}) (float64, bool) {
	switch v0 := v.(type) {
	case float64:
		return v0, true
	case float32:
		return float64(v0), true
	}
	return 0, false
}

/*ParseInt parse interface to int64 */
func ParseInt(v interface{}) (int64, bool) {
	switch v0 := v.(type) {
	case int:
		return int64(v0), true
	case int32:
		return int64(v0), true
	case int64:
		return int64(v0), true
	case uint:
		return int64(v0), true
	case uint32:
		return int64(v0), true
	case uint64:
		return int64(v0), true
	case float64:
		return int64(v0), true
	case float32:
		return int64(v0), true
	default:
	}
	return 0, false
}

/*ParseString parse interface to string */
func ParseString(v interface{}) (string, bool) {
	switch v0 := v.(type) {
	case string:
		return v0, true
	case []byte:
		return string(v0), true
	case bytes.Buffer:
		return v0.String(), true
	default:
	}
	return "", false
}

/*Time get time string */
func Time(t ...time.Time) string {
	if t == nil {
		return strconv.Itoa(time.Now().Nanosecond())
	}
	return strconv.Itoa(t[0].Nanosecond())
}

// CurrentTimeStampMS get current time with millisecond
func CurrentTimeStampMS() int64 {
	return time.Now().UnixNano() / time.Millisecond.Nanoseconds()
}

// CurrentTimeStampNS get current time with nanoseconds
func CurrentTimeStampNS() int64 {
	return time.Now().UnixNano()
}

// CurrentTimeStamp get current time with unix
func CurrentTimeStamp() int64 {
	return time.Now().Unix()
}

// CurrentTimeStampString get current time to string
func CurrentTimeStampString() string {
	return strconv.FormatInt(CurrentTimeStamp(), 10)
}

// SHA1 transfer string to sha1
func SHA1(s string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(s)))
}

// HmacSha256 ...
func HmacSha256(data []byte, key string) string {
	m := hmac.New(sha256.New, []byte(key))
	m.Write(data)
	return strings.ToUpper(fmt.Sprintf("%x", m.Sum(nil)))
}

//GenerateRandomString2 随机字符串
func GenerateRandomString2(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

//GenerateRandomString 随机字符串
func GenerateRandomString(size int, kind ...RandomKind) string {
	bytes := RandomString[RandomAll]
	if kind != nil {
		if k, b := RandomString[kind[0]]; b == true {
			bytes = k
		}
	}
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
