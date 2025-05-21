package swar

import "unsafe"

const (
	LowBits  uint64 = 0x0101_0101_0101_0101
	packMask uint64 = 0x0102_0408_1020_4080
)

// BytesToLanes casts a []byte to uint64's for SWAR processing, and an index of unused bytes.
func BytesToLanes(b []byte) ([]uint64, int) {
	countChunks := len(b) / 8
	chunks := unsafe.Slice((*uint64)(unsafe.Pointer(&b[0])), countChunks)
	return chunks, countChunks * 8
}

// LanesToBytes converts a slice of uint64 back to a byte slice.
func LanesToBytes(lanes []uint64) []byte {
	countBytes := len(lanes) * 8
	bytes := unsafe.Slice((*byte)(unsafe.Pointer(&lanes[0])), countBytes)
	return bytes
}

// Dupe duplicates a byte into a uint64, so that each bit of the byte is represented in each byte of the uint64.
func Dupe(c byte) uint64 {
	return uint64(c) * LowBits
}

func ExtractLowBits(v uint64) byte {
	return byte((v * packMask) >> 56)
}

func IntToLanes(i uint64) [8]byte {
	return *(*[8]byte)(unsafe.Pointer(&i))
}

func LanesToInt(lanes [8]byte) uint64 {
	return *(*uint64)(unsafe.Pointer(&lanes))
}

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
