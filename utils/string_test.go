package utils

import (
	"crypto/rand"
	"math"
	"math/big"
	"strings"
	"testing"
)

func BenchmarkRandomgString(b *testing.B) {
	type testSpec struct {
		name string
		size int
	}
	testCases := []testSpec{
		{name: "extra-small", size: 1},
		{name: "small", size: 1 << 3},  // 8
		{name: "medium", size: 1 << 6}, // 64
		{name: "large", size: 1 << 9},  // 512
	}

	type randomStringFunc func(n int) string
	randomers := map[string]randomStringFunc{
		"random_runes": func() randomStringFunc {
			letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			return func(n int) string {
				ret := make([]rune, n)
				for i := 0; i < n; i++ {
					num := MustNotError(rand.Int(rand.Reader, big.NewInt(int64(len(letters)))))
					ret[i] = letters[num.Int64()]
				}
				return string(ret)
			}
		}(),
		"random_bytes": func() randomStringFunc {
			letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			return func(n int) string {
				ret := make([]byte, n)
				for i := 0; i < n; i++ {
					num := MustNotError(rand.Int(rand.Reader, big.NewInt(int64(len(letters)))))
					ret[i] = letters[num.Int64()]
				}
				return string(ret)
			}
		}(),
		"random_bytes_optimized_1": func() randomStringFunc {
			// use strings.Builder to reduce allocations.
			letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			return func(n int) string {
				var b strings.Builder
				for i := 0; i < n; i++ {
					num := MustNotError(rand.Int(rand.Reader, big.NewInt(int64(len(letters)))))
					b.WriteByte(letters[num.Int64()%int64(len(letters))])
				}
				return b.String()
			}
		}(),
		"random_bytes_optimized_2": func() randomStringFunc {
			// declare letters as a const to reduce len() calls to O(1).
			const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
			return func(n int) string {
				ret := make([]byte, n)
				for i := 0; i < n; i++ {
					num := MustNotError(rand.Int(rand.Reader, big.NewInt(int64(len(letters)))))
					ret[i] = letters[num.Int64()]
				}
				return string(ret)
			}
		}(),
		"random_bytes_optimized_3": func() randomStringFunc {
			letters := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			return func(n int) string {
				ret := make([]byte, n)
				for i := 0; i < n; i++ {
					num := MustNotError(rand.Int(rand.Reader, big.NewInt(math.MaxInt64)))
					ret[i] = letters[num.Int64()%int64(len(letters))]
				}
				return string(ret)
			}
		}(),
		"random_bytes_optimized_all": func() randomStringFunc {
			const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
			return func(n int) string {
				var b strings.Builder
				b.Grow(n)
				for i := 0; i < n; i++ {
					num := MustNotError(rand.Int(rand.Reader, big.NewInt(math.MaxInt64)))
					b.WriteByte(letters[num.Int64()%int64(len(letters))])
				}
				return b.String()
			}
		}(),
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for name, random := range randomers {
				b.Run(name, func(b *testing.B) {
					for i := 0; i < b.N; i++ {
						random(tc.size)
						b.SetBytes(int64(tc.size))
					}
				})
			}
		})
	}
}
