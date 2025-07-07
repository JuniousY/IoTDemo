package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

func SecureRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

func HmacSha256Pass(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Ip2binary(ip string) string {
	str := strings.Split(ip, ".")
	var ipstr string
	for _, s := range str {
		i, err := strconv.ParseUint(s, 10, 8)
		if err != nil {
			fmt.Println(err)
		}
		ipstr = ipstr + fmt.Sprintf("%08b", i)
	}
	return ipstr
}

// MatchIP 测试IP地址和地址端是否匹配 变量ip为字符串，例子"192.168.56.4" iprange为地址端"192.168.56.64/26"
func MatchIP(ip, iprange string) bool {
	ipb := Ip2binary(ip)
	if strings.Contains(iprange, "/") { //如果是ip段
		ipr := strings.Split(iprange, "/")
		masklen, err := strconv.ParseUint(ipr[1], 10, 32)
		if err != nil {
			return false
		}
		iprb := Ip2binary(ipr[0])
		return strings.EqualFold(ipb[0:masklen], iprb[0:masklen])
	} else {
		return ip == iprange
	}
}
