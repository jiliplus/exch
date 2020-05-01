package exch

import (
	"bytes"
	"encoding/gob"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_maxFloat64(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{

		{
			"a 大",
			args{
				a: 1,
				b: 0,
			},
			1,
		},

		{
			"b 大",
			args{
				a: 1,
				b: 2,
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxFloat64(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("maxFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_minFloat64(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"b 大",
			args{
				a: 1,
				b: 2,
			},
			1,
		},

		{
			"b 小",
			args{
				a: 1,
				b: 0,
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minFloat64(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("minFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_EncFunc(t *testing.T) {
	Convey("生成 enc 以便于对变量进行编码", t, func() {
		enc := EncFunc()
		res := enc(true)
		So(res, ShouldResemble, []byte{3, 2, 0, 1})
		res = enc(false)
		So(res, ShouldResemble, []byte{3, 2, 0, 0})
		Convey("如果对 nil 编码会 panic", func() {
			So(func() { enc(nil) }, ShouldPanic)
		})
	})
}

func Benchmark_EncFunc(b *testing.B) {
	enc := EncFunc()
	now := time.Now()
	for i := 1; i < b.N; i++ {
		enc(now)
	}
}

// enc 返回的函数能够将输入转换成 []byte
func enc(e interface{}) []byte {
	var bb bytes.Buffer
	enc := gob.NewEncoder(&bb)
	err := enc.Encode(e)
	if err != nil {
		panic("gob encode error:" + err.Error())
	}
	return bb.Bytes()
}

func Benchmark_enc(b *testing.B) {
	now := time.Now()
	for i := 1; i < b.N; i++ {
		enc(now)
	}
}
