package swar

import "unsafe"

const (
	// LowBits has the lowest bit set in each byte for value duplication
	LowBits  uint64 = 0x0101_0101_0101_0101
	// packMask packs low bits from each byte into a single byte
	packMask uint64 = 0x0102_0408_1020_4080
)

// BytesToLanes converts a []byte to []uint64 for SWAR processing
// Returns uint64 lanes and index where unused bytes begin
func BytesToLanes(b []byte) ([]uint64, int) {
	countChunks := len(b) / 8
	chunks := unsafe.Slice((*uint64)(unsafe.Pointer(&b[0])), countChunks)
	return chunks, countChunks * 8
}

// LanesToBytes converts []uint64 back to []byte
// Zero-copy conversion for optimal performance
func LanesToBytes(lanes []uint64) []byte {
	countBytes := len(lanes) * 8
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(&lanes[0])), countBytes)
	return bytes
}

// Dupe duplicates a byte across all 8 bytes of a uint64
// Creates comparison values for parallel operations
func Dupe(c byte) uint64 {
	return uint64(c) * LowBits
}

// ExtractLowBits packs the low bit from each byte into a single byte
// Compacts 8 comparison results into a single byte
func ExtractLowBits(v uint64) byte {
	return byte((v * packMask) >> 56)
}

// IntToLanes converts a uint64 to an 8-byte array
// Access individual bytes for mixed SWAR/byte-level operations
func IntToLanes(i uint64) [8]byte {
	return *(*[8]byte)(unsafe.Pointer(&i))
}

// LanesToInt converts an 8-byte array to uint64
// Zero-copy conversion from byte-level to SWAR format
func LanesToInt(lanes [8]byte) uint64 {
	return *(*uint64)(unsafe.Pointer(&lanes))
}

// Lookup provides precomputed data for optimized operations
// OnesPositions maps byte values to positions of their set bits
var Lookup = struct {
	OnesPositions [256][]int
}{
	func() (res [256][]int) {
		for b := range res {
			for i := 0; i < 8; i++ {
				if b>>i&1 == 1 {
					res[b] = append(res[b], i)
				}
			}
		}
		return
	}()}
