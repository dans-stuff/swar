package swar

const (
	// mEven selects even bytes in a uint64
	mEven uint64 = 0x00FF_00FF_00FF_00FF
	// mOdd selects odd bytes in a uint64
	mOdd uint64 = 0xFF00_FF00_FF00_FF00
	// laneNotHigh masks all bits except the high bit in each byte
	laneNotHigh uint64 = 0x7F7F_7F7F_7F7F_7F7F
)

// SubtractBytesWithWrapping performs byte-wise subtraction with wrapping
// Parallel subtraction across all 8 bytes with wrap-around behavior
func SubtractBytesWithWrapping(a, b uint64) uint64 {
	return ((a | HighBits) - (b &^ HighBits)) ^ ((a ^ ^b) & HighBits)
}

// SubtractBytesWithMinimum performs byte-wise subtraction clamped at zero
// Provides saturating subtraction to prevent underflow in all 8 bytes
func SubtractBytesWithMinimum(a, b uint64) uint64 {
	diff := ((a | HighBits) - (b &^ HighBits)) ^ ((a ^ ^b) & HighBits)
	bo := ((^a & b) | ((^a | b) & diff)) & HighBits
	return diff &^ ((bo >> 7) * 0xFF)
}

// AddBytesWithWrapping performs byte-wise addition with wrap-around
// Parallel addition across all 8 bytes with overflow wrapping to zero
func AddBytesWithWrapping(a, b uint64) uint64 {
	sum := (a & laneNotHigh) + (b & laneNotHigh)
	return sum ^ ((a ^ b) & HighBits)
}

// AddBytesWithMaximum performs byte-wise addition clamped at 255
// Saturating addition to prevent overflow in all 8 bytes
func AddBytesWithMaximum(a, b uint64) uint64 {
	preSum := (a & laneNotHigh) + (b & laneNotHigh)
	sum := preSum ^ ((a ^ b) & HighBits)
	carry := ((a & b) | ((a | b) & ^sum)) & HighBits
	return sum | (carry>>7)*0xFF
}

// AbsoluteDifferenceBetweenBytes calculates |a-b| for each byte
// Computes unsigned distances for metrics and signal processing
func AbsoluteDifferenceBetweenBytes(a, b uint64) uint64 {
	d := a - b
	borrow := ((^a & b) | ((^a | b) & d)) & HighBits
	mask := (borrow >> 7) * 0xFF
	n := (a &^ mask) | (b & mask)
	m := (a & mask) | (b &^ mask)
	return ((n | HighBits) - (m &^ HighBits)) ^ ((n ^ ^m) & HighBits)
}

// SelectSmallerBytes returns min(a,b) for each byte
// Efficient for clipping, filtering, and data preprocessing
func SelectSmallerBytes(a, b uint64) uint64 {
	d := a - b
	borrow := ((^a & b) | ((^a | b) & d)) & HighBits
	mask := (borrow >> 7) * 0xFF
	return (a & mask) | (b &^ mask)
}

// SelectLargerBytes returns max(a,b) for each byte
// Ideal for peak detection, ceiling operations, and filtering
func SelectLargerBytes(a, b uint64) uint64 {
	d := a - b
	borrow := ((^a & b) | ((^a | b) & d)) & HighBits
	mask := (borrow >> 7) * 0xFF
	return (a &^ mask) | (b & mask)
}

// AverageBytes calculates (a+b)/2 for each byte without overflow
// Perfect for signal processing, image manipulation, and smoothing
func AverageBytes(a, b uint64) uint64 {
	common := a & b
	diff := (a ^ b) & 0xFEFE_FEFE_FEFE_FEFE
	return common + (diff >> 1)
}

// SwapByteHalves swaps the high and low nibbles in each byte
// Useful for BCD encoding/decoding and nibble-level transforms
func SwapByteHalves(v uint64) uint64 {
	lo := v & 0x0F0F_0F0F_0F0F_0F0F
	hi := v & 0xF0F0_F0F0_F0F0_F0F0
	return (lo << 4) | (hi >> 4)
}

// ReverseEachByte reverses the bit order within each byte
// Useful for endianness conversion and bit-level manipulations
func ReverseEachByte(v uint64) uint64 {
	x := ((v >> 1) & 0x5555555555555555) | ((v & 0x5555555555555555) << 1)
	x = ((x >> 2) & 0x3333333333333333) | ((x & 0x3333333333333333) << 2)
	x = ((x >> 4) & 0x0F0F0F0F0F0F0F0F) | ((x & 0x0F0F0F0F0F0F0F0F) << 4)
	return x
}

// SelectByLowBit selects values from a or b based on mask bits
// Branchless selection between values based on conditions
func SelectByLowBit(a, b, mask uint64) uint64 {
	byteMask := mask * 0xFF
	return (a & byteMask) | (b &^ byteMask)
}

// CountOnesPerByte counts set bits in each byte
// Parallel population count for hamming distance and feature extraction
func CountOnesPerByte(v uint64) uint64 {
	m1 := v - ((v >> 1) & 0x5555_5555_5555_5555)
	m2 := (m1 & 0x3333_3333_3333_3333) + ((m1 >> 2) & 0x3333_3333_3333_3333)
	return (m2 + (m2 >> 4)) & 0x0F0F_0F0F_0F0F_0F0F
}
