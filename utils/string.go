package utils

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	var s strings.Builder
	s.Grow(n)
	for i := 0; i < n; i++ {
		num := MustNotError(rand.Int(rand.Reader, big.NewInt(math.MaxInt64)))
		s.WriteByte(letters[num.Int64()%int64(len(letters))])
	}
	return s.String()
}

// CopyBytes copies a slice to make it immutable
func CopyBytes(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

const (
	uByte = 1 << (10 * iota)
	uKilobyte
	uMegabyte
	uGigabyte
	uTerabyte
	uPetabyte
	uExabyte
)

// ByteSize returns a human-readable byte string, eg. 10M, 12.5K.
func ByteSize(bytes uint64) string {
	unit := ""
	value := float64(bytes)
	switch {
	case bytes >= uExabyte:
		unit = "EB"
		value /= uExabyte
	case bytes >= uPetabyte:
		unit = "PB"
		value /= uPetabyte
	case bytes >= uTerabyte:
		unit = "TB"
		value /= uTerabyte
	case bytes >= uGigabyte:
		unit = "GB"
		value /= uGigabyte
	case bytes >= uMegabyte:
		unit = "MB"
		value /= uMegabyte
	case bytes >= uKilobyte:
		unit = "KB"
		value /= uKilobyte
	case bytes >= uByte:
		unit = "B"
	default:
		return "0B"
	}
	result := strconv.FormatFloat(value, 'f', 1, 64)
	result = strings.TrimSuffix(result, ".0")
	return result + unit
}

// ToString Change any supports types to string
func ToString(arg interface{}, args ...any) string {
	tmp := reflect.Indirect(reflect.ValueOf(arg)).Interface()
	switch v := tmp.(type) {
	case string:
		return fmt.Sprintf(v, args...)
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.Itoa(int(v))
	case uint8:
		return strconv.FormatInt(int64(v), 10)
	case uint16:
		return strconv.FormatInt(int64(v), 10)
	case uint32:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatInt(int64(v), 10)
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		if len(args) > 0 {
			if format, ok := args[0].(string); ok {
				return v.Format(format)
			}
		}
		return v.Format("2006-01-02 15:04:05")
	case reflect.Value:
		return ToString(v.Interface(), args...)
	case fmt.Stringer:
		return v.String()
	default:
		return ""
	}
}

const twoAsterisks = "**"

// Match matches patterns as gitignore pattern.
// Reference https://git-scm.com/docs/gitignore, https://gist.github.com/jstnlvns/ebaa046fae16543cc9efc7f24bcd0e31
func Match(pattern, value string) bool {
	if pattern == "" {
		return false
	}

	// Code Comment
	if strings.HasPrefix(pattern, "#") {
		return false
	}

	pattern = strings.TrimSuffix(pattern, " ")
	pattern = strings.TrimSuffix(pattern, string(os.PathSeparator))

	// Reverse result if pattern starts with "!"
	neg := strings.HasPrefix(pattern, "!")
	if neg {
		pattern = strings.TrimPrefix(pattern, "!")
	}

	// Two Consecutive Asterisks ("**")
	if strings.Contains(pattern, twoAsterisks) {
		result := MatchTwoAsterisk(pattern, value)
		if neg {
			result = !result
		}
		return result
	}

	// Shell-style Pattern Matching
	matched, err := path.Match(pattern, value)
	if err != nil {
		return false
	}

	if neg {
		return !matched
	}

	return matched
}

func MatchTwoAsterisk(pattern, value string) bool {
	// **.openapi.json == fund-api.openapi.json
	if strings.HasPrefix(pattern, twoAsterisks) {
		pattern = strings.TrimPrefix(pattern, twoAsterisks)
		return strings.HasSuffix(value, pattern)
	}

	// docs/** == docs/README.md or index.** == index.json, index.yaml
	if strings.HasSuffix(pattern, twoAsterisks) {
		pattern = strings.TrimSuffix(pattern, twoAsterisks)
		return strings.HasPrefix(value, pattern)
	}

	// "a/**/b" == "a/b", /"a/x/b", "a/x/y/b"
	parts := strings.Split(pattern, twoAsterisks)
	for i, part := range parts {
		switch i {
		case 0: // first part
			if !strings.HasPrefix(value, part) {
				return false
			}
		case len(parts) - 1: // last part
			part = strings.TrimPrefix(part, string(os.PathSeparator))
			return strings.HasSuffix(value, part)
		default:
			if !strings.Contains(value, part) {
				return false
			}
		}

		index := strings.Index(value, part) + len(part)
		value = value[index:]
	}

	return false
}
