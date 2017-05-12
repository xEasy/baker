package app

import (
	"testing"
)

func TestGenQrcodeImg(t *testing.T) {
	if _, e := GenQrcodeImg("hehehe", 320); e != nil { //try a unit test on function
		t.Error("TestGenQrcodeImg did not work as expected.", e) // 如果不是如预期的那么就报错
	} else {
		t.Log("TestGenQrcodeImg test passed.") //记录一些你期望记录的信息
	}
}

func BenchmarkGenQrcodeImg(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		GenQrcodeImg("hehehe", 320)
	}
}
