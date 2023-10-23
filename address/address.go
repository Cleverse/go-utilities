package address

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/Cleverse/go-utilities/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// AddressLength is the expected length of the address
	AddressLength = common.AddressLength
)

var (
	// Ether is the default address of the Ethereum's native currency.
	Ether = common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")

	// Zero is zero value of common.Address.
	Zero = common.Address{}

	// DeadWalletAddress is the defailt address of the dead wallet.
	DeadWalletAddress = common.HexToAddress("0x000000000000000000000000000000000000dead")

	// DeadAddresseses is the list of dead addresses.
	DeadAddresseses = map[common.Address]interface{}{
		Zero:              nil, // 0x0000000000000000000000000000000000000000
		DeadWalletAddress: nil, // 0x000000000000000000000000000000000000dead
		common.HexToAddress("0x0000000000000000000000000000000000000001"): nil, // Safemoon, GALA, AWC, NOW
		common.HexToAddress("0xdEAD000000000000000042069420694206942069"): nil, // SHIB, BOHR, HEGIC, etc in Ethereum
		common.HexToAddress("0xdead000000000000000000000000000000000000"): nil,
	}

	// RandomGenerator is the address random generator function and used by address.Random() function.
	//
	// Default address random generator is cryptographically secure random number generator.
	RandomGenerator func() (addr common.Address) = RandomFromBytes
)

// FromHex safely converts a hex string to a common.Address.
func FromHex(address string) (addr common.Address) {
	address = strings.TrimSpace(address)
	if len(address) == 0 {
		return Zero
	}
	return common.HexToAddress(address)
}

// FromString safely converts a string to a common.Address.
func FromString(address string) (addr common.Address) {
	return FromHex(address)
}

// FromHexes safely converts hex string slice to a common.Address slice.
func FromHexes(src []string) (dst []common.Address) {
	dst = make([]common.Address, 0, len(src))
	for _, addr := range src {
		dst = append(dst, FromHex(addr))
	}
	return dst
}

// FromStrings alias of FromHexes. safely converts string slice to a common.Address slice.
func FromStrings(src []string) (dst []common.Address) {
	return FromHexes(src)
}

// ToLower converts an address to a lower-case string withouth checksum.
func ToLower(a common.Address) string {
	buf := [common.AddressLength*2 + 2]byte{}
	copy(buf[:2], []byte("0x"))
	_ = hex.Encode(buf[2:], a[:])
	return string(buf[:])
}

// ToString alias of ToLower. converts an address to a lower-case string withouth checksum.
func ToString(a common.Address) string {
	return ToLower(a)
}

// ToLowers converts an addresses to a lower-case strings withouth checksum.
func ToLowers(src []common.Address) []string {
	dst := make([]string, 0, len(src))
	for _, addr := range src {
		dst = append(dst, ToLower(addr))
	}
	return dst
}

// ToStrings alias of ToLowers. converts an addresses to a lower-case strings withouth checksum.
func ToStrings(src []common.Address) []string {
	return ToLowers(src)
}

// IsZero returns `true` if the address is zero value.
func IsZero(a common.Address) bool {
	return a == Zero
}

// IsEmpty returns `true` if the address is empty (alias of IsZero)
func IsEmpty(a common.Address) bool {
	return a == Zero
}

// IsDead returns `true` if the address is dead address.
func IsDead(a common.Address) bool {
	_, ok := DeadAddresseses[a]
	return ok
}

// IsZero returns `true` if the address is can't be used. (zero value or dead address)
func IsValid(a common.Address) bool {
	return !IsZero(a) && !IsDead(a)
}

// Random returns a random common.Address. can be changed address random generator by RandomGenerator variable.
//
// Default address random generator is use cryptographically secure random number generator.
func Random() (addr common.Address) {
	if RandomGenerator == nil {
		RandomGenerator = RandomFromBytes
	}
	return RandomGenerator()
}

// RandomFromPrivateKey returns a random address from a random private key
func RandomFromPrivateKey() common.Address {
	return crypto.PubkeyToAddress(utils.Must(crypto.GenerateKey()).PublicKey)
}

// RandomFromBytes returns a random address from a random byte slice (via crypto/rand)
func RandomFromBytes() (addr common.Address) {
	_, _ = rand.Read(addr[:])
	return addr
}
