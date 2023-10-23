package address

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func BenchmarkRandomAddress(b *testing.B) {
	b.Run("rand/privatekey", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = RandomFromPrivateKey()
			b.SetBytes(AddressLength)
		}
	})
	b.Run("crypto/rand", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			buf := make([]byte, AddressLength)
			_, _ = rand.Read(buf)
			_ = common.HexToAddress(hex.EncodeToString(buf))
			b.SetBytes(AddressLength)
		}
	})
	b.Run("crypto/rand /w optimize", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = RandomFromBytes()
			b.SetBytes(AddressLength)
		}
	})
}

func BenchmarkToLower(b *testing.B) {
	address := common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	var sizeAddressString int64 = 42
	toStrings := map[string]func(common.Address) string{
		"Strings.ToLower": func(a common.Address) string {
			return strings.ToLower(a.String())
		},
		"Old": func(a common.Address) string {
			hexEncoding := hex.EncodeToString(a.Bytes())
			return fmt.Sprintf("0x%s", hexEncoding)
		},
		"StringBuilder_1": func(a common.Address) string {
			var s strings.Builder
			s.WriteString("0x")
			s.WriteString(hex.EncodeToString(a.Bytes()))
			return s.String()
		},
		"StringBuilder_2": func(a common.Address) string {
			var s strings.Builder
			b := make([]byte, 40)
			_ = hex.Encode(b, a.Bytes()) // Avoid conversion between []byte and string to reduce memory allocation.
			s.WriteString("0x")
			s.Write(b)
			return s.String()
		},
		"StringBuilder_3": func(a common.Address) string {
			var s strings.Builder
			b := make([]byte, 40)
			_ = hex.Encode(b, a.Bytes()) // Avoid conversion between []byte and string to reduce memory allocation.
			s.Write([]byte{48, 120})     // use Write instead of WriteString
			s.Write(b)
			return s.String()
		},
		"Now": ToLower,
	}
	ordered := []string{"Strings.ToLower", "Old", "StringBuilder_1", "StringBuilder_2", "StringBuilder_3", "Now"}
	for _, name := range ordered {
		toString := toStrings[name]
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = toString(address)
				b.SetBytes(sizeAddressString)
			}
		})
	}
	/*
		goos: darwin
		goarch: arm64
		BenchmarkToLower/Strings.ToLower-8            	 1529884	       808.3 ns/op	  51.96 MB/s	    1104 B/op	       7 allocs/op
		BenchmarkToLower/Old-8               			10714300	       104.1 ns/op	 403.51 MB/s	     160 B/op	       4 allocs/op
		BenchmarkToLower/StringBuilder_1-8            	14235339	        88.43 ns/op	 474.93 MB/s	     152 B/op	       4 allocs/op
		BenchmarkToLower/StringBuilder_2-8            	24213360	        49.09 ns/op	 855.51 MB/s	      56 B/op	       2 allocs/op
		BenchmarkToLower/StringBuilder_3-8            	23675196	        49.30 ns/op	 851.92 MB/s	      56 B/op	       2 allocs/op
		BenchmarkToLower/Now-8      					39285494	        30.81 ns/op	1363.39 MB/s	      48 B/op	       1 allocs/op
	*/
}

func BenchmarkToString(b *testing.B) {
	address := common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	var sizeAddressString int64 = 42

	toStrings := map[string]func(common.Address) string{
		"Address.String()": func(a common.Address) string {
			return a.String()
		},
		"Old": func(a common.Address) string {
			hexEncoding := hex.EncodeToString(a.Bytes())
			return fmt.Sprintf("0x%s", hexEncoding)
		},
		"StringBuilder": func(a common.Address) string {
			var s strings.Builder
			s.WriteString("0x")
			s.WriteString(hex.EncodeToString(a.Bytes()))
			return s.String()
		},
		"Now": ToString,
	}

	for name, toString := range toStrings {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = toString(address)
				b.SetBytes(sizeAddressString)
			}
		})
	}
	/*
		goos: darwin
		goarch: arm64
		BenchmarkToString/Address.String()-8     	 1906382	       632.1 ns/op	  66.44 MB/s	    1056 B/op	       6 allocs/op
		BenchmarkToString/Old-8         			10145156	       103.6 ns/op	 405.57 MB/s	     160 B/op	       4 allocs/op
		BenchmarkToString/StringBuilder-8        	14861654	        81.47 ns/op	 515.50 MB/s	     152 B/op	       4 allocs/op
		BenchmarkToString/Now-8                   	39218355	        30.88 ns/op	1359.94 MB/s	      48 B/op	       1 allocs/op
	*/
}

func TestToLower(t *testing.T) {
	address := "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
	expected := strings.ToLower(address)
	actual := ToLower(common.HexToAddress(expected))
	assert.Equal(t, expected, actual)
}

func TestToLowers(t *testing.T) {
	addresses := []common.Address{
		common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"),
		common.HexToAddress("0xABCDEFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		common.HexToAddress("0xdEAD000000000000000042069420694206942069"),
	}

	actuals := ToLowers(addresses)
	for i, actual := range actuals {
		assert.Equal(t, strings.ToLower(addresses[i].String()), actual)
	}
}

func TestToString(t *testing.T) {
	expected := "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"
	actual := ToString(common.HexToAddress(expected))
	assert.Equal(t, strings.ToLower(expected), actual)
}

func TestToStrings(t *testing.T) {
	addresses := []common.Address{
		common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"),
		common.HexToAddress("0xABCDEFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		common.HexToAddress("0xdEAD000000000000000042069420694206942069"),
	}

	actuals := ToStrings(addresses)
	for i, actual := range actuals {
		assert.Equal(t, strings.ToLower(addresses[i].String()), actual)
	}
}

func TestFromStrings(t *testing.T) {
	expected := []common.Address{
		common.HexToAddress("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"),
		common.HexToAddress("0xABCDEFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		common.HexToAddress("0xdEAD000000000000000042069420694206942069"),
	}
	expectedString := ToStrings(expected)

	actuals := FromStrings(expectedString)
	for i, actual := range actuals {
		assert.Equal(t, expected[i], actual)
	}
}
