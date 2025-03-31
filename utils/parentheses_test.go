package utils

import "testing"

// 27.14 ns/op
func BenchmarkValidParenthesesASCII(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidParentheses("{[()()]}")
	}
}

// 61.29 ns/op
func BenchmarkValidParenthesesMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidParenthesesMap("{[()()]}")
	}
}
