package utils

import "testing"

// 68.45 ns/op
func BenchmarkIsPrimeOld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrime(97)
	}
}

// 1.256 ns/op
func BenchmarkIsPrimeNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPrimeNew(97)
	}
}
