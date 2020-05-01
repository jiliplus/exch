package exch

import (
	"bytes"
	"encoding/gob"
)

func maxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minFloat64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// EncFunc 返回的函数能够将输入转换成 []byte
// 这些写成闭包的形式，而不是 enc
// 是为了提高速度
// 运行压力测试即可知道，提速了 6 倍
func EncFunc() func(interface{}) []byte {
	var bb bytes.Buffer
	enc := gob.NewEncoder(&bb)
	return func(e interface{}) []byte {
		bb.Reset()
		err := enc.Encode(e)
		if err != nil {
			panic("gob encode error:" + err.Error())
		}
		return bb.Bytes()
	}
}
