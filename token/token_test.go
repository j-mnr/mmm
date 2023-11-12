package token_test

import (
	"mmm/token"
	"runtime"
	"testing"
)

func BenchmarkNew(b *testing.B) {
	loop := int(token.TypeReturn + 1)
	var t token.Token
	for i := 0; i < b.N; i++ {
		t = token.New(token.Type(i%loop), "XXXX")
	}
	runtime.KeepAlive(t)
}
