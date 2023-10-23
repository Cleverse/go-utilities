package address

import (
	"bytes"

	"github.com/ethereum/go-ethereum/common"
)

// SortAddress
// for sorting two addresses by hex value.
func SortAddress(address0 common.Address, address1 common.Address) (common.Address, common.Address) {
	if CompareAddress(address0, address1) > 0 {
		return address1, address0
	}
	return address0, address1
}

// CompareAddress
// returns -1 if address0 < address1
// returns 0 if address0 == address1
// returns 1 if address0 > address1
func CompareAddress(address0 common.Address, address1 common.Address) int {
	return bytes.Compare(address0.Bytes(), address1.Bytes())
}
