package utils

import (
	"errors"
	"testing"
)

func BenchmarkIsRetryable(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IsRetryable(ErrNotConnected)
		}
	})
}

func BenchmarkIsTemporary(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IsTemporary(ErrConnectionTimeout)
		}
	})
}

func BenchmarkNewError(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			NewError(testKindConnection, testOpConnect, ErrNotConnected)
		}
	})
}

func BenchmarkErrorChain(b *testing.B) {
	wrapped := WrapError(WrapError(errors.New(testMsgRootErr), testMsgLevel1), testMsgLevel2)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ErrorChain(wrapped)
		}
	})
}
