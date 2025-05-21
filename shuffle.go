package swar

// SwapByteHalves swaps the high and low nibbles of each byte in a 64-bit integer.
func SwapByteHalves(v uint64) uint64 {
	lo := v & 0x0F0F_0F0F_0F0F_0F0F
	hi := v & 0xF0F0_F0F0_F0F0_F0F0
	return (lo << 4) | (hi >> 4)
}

// ReverseByteHalves reverses each byte of a 64-bit integer.
func ReverseEachByte(v uint64) uint64 {
	x := ((v >> 1) & 0x5555555555555555) | ((v & 0x5555555555555555) << 1)
	x = ((x >> 2) & 0x3333333333333333) | ((x & 0x3333333333333333) << 2)
	x = ((x >> 4) & 0x0F0F0F0F0F0F0F0F) | ((x & 0x0F0F0F0F0F0F0F0F) << 4)
	return x
}

func SelectByLowBit(a, b, mask uint64) uint64 {
	byteMask := mask * 0xFF
	return (a & byteMask) | (b &^ byteMask)
}

func CountOnesPerByte(v uint64) uint64 {
	m1 := v - ((v >> 1) & 0x5555_5555_5555_5555)
	m2 := (m1 & 0x3333_3333_3333_3333) + ((m1 >> 2) & 0x3333_3333_3333_3333)
	return (m2 + (m2 >> 4)) & 0x0F0F_0F0F_0F0F_0F0F
}
