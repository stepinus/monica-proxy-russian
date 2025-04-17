package utils

import (
	"math/rand"
	"time"
)

var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandStringUsingMathRand 生成指定长度的随机字符串
func RandStringUsingMathRand(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	// 创建一个长度为 n 的切片，用来存放随机字符
	result := make([]rune, n)
	for i := range n {
		result[i] = letters[randSource.Intn(len(letters))]
	}

	return string(result)
}
