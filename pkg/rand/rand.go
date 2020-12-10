package rand

import (
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().Unix()))
var source = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// From 以 source 作为元字符生成 bits 位随机字符
func From(bits int, source []byte) []byte {
	bytes := make([]byte, bits)
	n := len(source)
	for i := 0; i < bits; i++ {
		bytes[i] = source[r.Intn(n)]
	}
	return bytes
}

// Str 生成 bits 位随机字符串
func Str(bits int) string {
	return string(From(bits, source))
}
