package utils

// CopyMapOfArray high performance copy map of array
// to new map of array with shared backing array.
// Reference: https://cs.opensource.google/go/go/+/master:src/net/http/header.go;l=94
func CopyMapOfArray[K comparable, V any](src map[K][]V) map[K][]V {
	if src == nil {
		return nil
	}

	// Find total number of values.
	totalValue := 0
	for _, val := range src {
		totalValue += len(val)
	}

	tmp := make([]V, totalValue) // use shared backing array for reduce memory allocation.
	dst := make(map[K][]V, len(src))
	for k, val := range src {
		if val == nil {
			dst[k] = nil
			continue
		}
		n := copy(tmp, val) // copy values to shared array.
		dst[k] = tmp[:n:n]  // point to specific length and capacity of shared backing array.
		tmp = tmp[n:]       // move pointer to next position.
	}

	return dst
}
