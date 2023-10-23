package utils

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"
)

func BenchmarkCopySlice(b *testing.B) {
	sizes := []int{8, 16, 64}

	for _, mapSize := range sizes {
		b.Run(fmt.Sprintf("map_size_%02d", mapSize), func(b *testing.B) {
			sliceSize := 32
			org := make(map[int][]string, mapSize)
			for i := 0; i < mapSize; i++ {
				org[i] = strings.Split(RandomString(sliceSize), "")
			}

			stringByte := int(unsafe.Sizeof(""))
			sliceBytes := int(unsafe.Sizeof(make([]string, sliceSize)))
			mapBytes := int(unsafe.Sizeof(org))
			totalSize := (sliceBytes + (stringByte*sliceSize)*mapSize) + mapBytes

			b.Run("NormalCopy", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					clone := copyNormal(org)
					_ = clone
					b.SetBytes(int64(totalSize))
				}
			})

			b.Run("Optimized", func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					clone := CopyMapOfArray(org)
					_ = clone
					b.SetBytes(int64(totalSize))
				}
			})
		})
	}
}

func copyNormal(org map[int][]string) map[int][]string {
	clone := make(map[int][]string, len(org))
	for k, val := range org {
		if val == nil {
			clone[k] = nil
			continue
		}
		clone[k] = make([]string, len(val))
		copy(clone[k], val)
	}

	return clone
}
