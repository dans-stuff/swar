package swar

const (
	mEven       uint64 = 0x00FF_00FF_00FF_00FF
	mOdd        uint64 = 0xFF00_FF00_FF00_FF00
	laneNotHigh uint64 = 0x7F7F_7F7F_7F7F_7F7F
)

func SubtractBytesWithWrapping(a, b uint64) uint64 {
	return ((a | HighBits) - (b &^ HighBits)) ^ ((a ^ ^b) & HighBits)
}

func SubtractBytesWithMinimum(a, b uint64) uint64 {
	diff := ((a | HighBits) - (b &^ HighBits)) ^ ((a ^ ^b) & HighBits)
	bo := ((^a & b) | ((^a | b) & diff)) & HighBits
	return diff &^ ((bo >> 7) * 0xFF)
}

func AddBytesWithWrapping(a, b uint64) uint64 {
	sum := (a & laneNotHigh) + (b & laneNotHigh)
	return sum ^ ((a ^ b) & HighBits)
}

func AddBytesWithMaximum(a, b uint64) uint64 {
	preSum := (a & laneNotHigh) + (b & laneNotHigh)
	sum := preSum ^ ((a ^ b) & HighBits)
	carry := ((a & b) | ((a | b) & ^sum)) & HighBits
	return sum | (carry>>7)*0xFF
}

func AbsoluteDifferenceBetweenBytes(a, b uint64) uint64 {
	d := a - b
	borrow := ((^a & b) | ((^a | b) & d)) & HighBits
	mask := (borrow >> 7) * 0xFF
	n := (a &^ mask) | (b & mask)
	m := (a & mask) | (b &^ mask)
	return ((n | HighBits) - (m &^ HighBits)) ^ ((n ^ ^m) & HighBits)
}

func SelectSmallerBytes(a, b uint64) uint64 {
	d := a - b
	borrow := ((^a & b) | ((^a | b) & d)) & HighBits
	mask := (borrow >> 7) * 0xFF
	return (a & mask) | (b &^ mask)
}

func SelectLargerBytes(a, b uint64) uint64 {
	d := a - b
	borrow := ((^a & b) | ((^a | b) & d)) & HighBits
	mask := (borrow >> 7) * 0xFF
	return (a &^ mask) | (b & mask)
}

func AverageBytes(a, b uint64) uint64 {
	common := a & b
	diff := (a ^ b) & 0xFEFE_FEFE_FEFE_FEFE
	return common + (diff >> 1)
}
