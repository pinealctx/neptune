package tex

import (
	"strconv"
)

//I64Hex convert int64 to hex(hexadecimal 16-base) string
func I64Hex(i int64) string {
	return strconv.FormatInt(i, 16)
}

//U64Hex convert uint64 to hex(hexadecimal 16-base) string
func U64Hex(u uint64) string {
	return strconv.FormatUint(u, 16)
}

//I64HexV2 convert int64 to 32 base string
func I64HexV2(i int64) string {
	return strconv.FormatInt(i, 32)
}

//U64HexV2 convert uint64 to 32 base string
func U64HexV2(u uint64) string {
	return strconv.FormatUint(u, 32)
}

//HexI64 convert hex string(16-base) to int64
func HexI64(s string) (int64, error) {
	return strconv.ParseInt(s, 16, 64)
}

//HexU64 convert hex string(16-base) to uint64
func HexU64(s string) (uint64, error) {
	return strconv.ParseUint(s, 16, 64)
}

//HexI64V2 convert 32 base string to int64
func HexI64V2(s string) (int64, error) {
	return strconv.ParseInt(s, 32, 64)
}

//HexU64V2 convert 32 base string to uint64
func HexU64V2(s string) (uint64, error) {
	return strconv.ParseUint(s, 32, 64)
}
