package main

import (
	"strings"
	"testing"
)

var s = strings.Repeat("abcdefghij", 100)

func BenchmarkToString(b *testing.B) {
	d := []byte(s)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = toString(d)
	}
}

func BenchmarkUnsafeToString(b *testing.B) {
	d := []byte(s)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = unsafeToString(d)
	}
}