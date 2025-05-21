package swar

import (
	"math/bits"
	"testing"
)

func TestToBytes(t *testing.T) {
	in := uint64(0x000000bfe5bd580c)
	expected := []byte{0x00, 0x00, 0x00, 0xbf, 0xe5, 0xbd, 0x58, 0x0c}
	out := toBytes(in)
	if len(out) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(out))
	}
	for i := 0; i < len(expected); i++ {
		if out[i] != expected[i] {
			t.Errorf("Expected byte %d to be %x, got %x", i, expected[i], out[i])
		}
	}
}

// helper: big-endian lanes
func toBytes(v uint64) [8]byte {
	var b [8]byte
	for i := 0; i < 8; i++ {
		shift := uint((7 - i) * 8)
		b[i] = byte((v >> shift) & 0xFF)
	}
	return b
}

func fromBytes(b [8]byte) uint64 {
	var v uint64
	for i := 0; i < 8; i++ {
		shift := uint((7 - i) * 8)
		v |= uint64(b[i]) << shift
	}
	return v
}

func byteFromLowBits(b [8]byte) uint8 {
	var out uint8
	if b[0]&1 != 0 {
		out |= 1 << 7
	}
	if b[1]&1 != 0 {
		out |= 1 << 6
	}
	if b[2]&1 != 0 {
		out |= 1 << 5
	}
	if b[3]&1 != 0 {
		out |= 1 << 4
	}
	if b[4]&1 != 0 {
		out |= 1 << 3
	}
	if b[5]&1 != 0 {
		out |= 1 << 2
	}
	if b[6]&1 != 0 {
		out |= 1 << 1
	}
	if b[7]&1 != 0 {
		out |= 1 << 0
	}
	return out
}

func addWrapBytes(ba, bb [8]byte) [8]byte {
	ba[0] = ba[0] + bb[0]
	ba[1] = ba[1] + bb[1]
	ba[2] = ba[2] + bb[2]
	ba[3] = ba[3] + bb[3]
	ba[4] = ba[4] + bb[4]
	ba[5] = ba[5] + bb[5]
	ba[6] = ba[6] + bb[6]
	ba[7] = ba[7] + bb[7]
	return ba
}

func minBytes(ba, bb [8]byte) [8]byte {
	if bb[0] < ba[0] {
		ba[0] = bb[0]
	}
	if bb[1] < ba[1] {
		ba[1] = bb[1]
	}
	if bb[2] < ba[2] {
		ba[2] = bb[2]
	}
	if bb[3] < ba[3] {
		ba[3] = bb[3]
	}
	if bb[4] < ba[4] {
		ba[4] = bb[4]
	}
	if bb[5] < ba[5] {
		ba[5] = bb[5]
	}
	if bb[6] < ba[6] {
		ba[6] = bb[6]
	}
	if bb[7] < ba[7] {
		ba[7] = bb[7]
	}
	return ba
}

func maxBytes(ba, bb [8]byte) [8]byte {
	if bb[0] > ba[0] {
		ba[0] = bb[0]
	}
	if bb[1] > ba[1] {
		ba[1] = bb[1]
	}
	if bb[2] > ba[2] {
		ba[2] = bb[2]
	}
	if bb[3] > ba[3] {
		ba[3] = bb[3]
	}
	if bb[4] > ba[4] {
		ba[4] = bb[4]
	}
	if bb[5] > ba[5] {
		ba[5] = bb[5]
	}
	if bb[6] > ba[6] {
		ba[6] = bb[6]
	}
	if bb[7] > ba[7] {
		ba[7] = bb[7]
	}
	return ba
}

func averageBytes(ba, bb [8]byte) [8]byte {
	ba[0] = byte((int(ba[0]) + int(bb[0])) / 2)
	ba[1] = byte((int(ba[1]) + int(bb[1])) / 2)
	ba[2] = byte((int(ba[2]) + int(bb[2])) / 2)
	ba[3] = byte((int(ba[3]) + int(bb[3])) / 2)
	ba[4] = byte((int(ba[4]) + int(bb[4])) / 2)
	ba[5] = byte((int(ba[5]) + int(bb[5])) / 2)
	ba[6] = byte((int(ba[6]) + int(bb[6])) / 2)
	ba[7] = byte((int(ba[7]) + int(bb[7])) / 2)
	return ba
}

func swapNibbles(b [8]byte) [8]byte {
	b[0] = (b[0]&0x0F)<<4 | (b[0]&0xF0)>>4
	b[1] = (b[1]&0x0F)<<4 | (b[1]&0xF0)>>4
	b[2] = (b[2]&0x0F)<<4 | (b[2]&0xF0)>>4
	b[3] = (b[3]&0x0F)<<4 | (b[3]&0xF0)>>4
	b[4] = (b[4]&0x0F)<<4 | (b[4]&0xF0)>>4
	b[5] = (b[5]&0x0F)<<4 | (b[5]&0xF0)>>4
	b[6] = (b[6]&0x0F)<<4 | (b[6]&0xF0)>>4
	b[7] = (b[7]&0x0F)<<4 | (b[7]&0xF0)>>4
	return b
}

func reverseBits(b [8]byte) [8]byte {
	b[0] = bits.Reverse8(b[0])
	b[1] = bits.Reverse8(b[1])
	b[2] = bits.Reverse8(b[2])
	b[3] = bits.Reverse8(b[3])
	b[4] = bits.Reverse8(b[4])
	b[5] = bits.Reverse8(b[5])
	b[6] = bits.Reverse8(b[6])
	b[7] = bits.Reverse8(b[7])
	return b
}

func popcountPerByte(b [8]byte) [8]byte {
	b[0] = byte(bits.OnesCount8(b[0]))
	b[1] = byte(bits.OnesCount8(b[1]))
	b[2] = byte(bits.OnesCount8(b[2]))
	b[3] = byte(bits.OnesCount8(b[3]))
	b[4] = byte(bits.OnesCount8(b[4]))
	b[5] = byte(bits.OnesCount8(b[5]))
	b[6] = byte(bits.OnesCount8(b[6]))
	b[7] = byte(bits.OnesCount8(b[7]))
	return b
}

func addSatBytes(ba, bb [8]byte) [8]byte {
	if bb[0] > 0xFF-ba[0] {
		ba[0] = 0xFF
	} else {
		ba[0] += bb[0]
	}
	if bb[1] > 0xFF-ba[1] {
		ba[1] = 0xFF
	} else {
		ba[1] += bb[1]
	}
	if bb[2] > 0xFF-ba[2] {
		ba[2] = 0xFF
	} else {
		ba[2] += bb[2]
	}
	if bb[3] > 0xFF-ba[3] {
		ba[3] = 0xFF
	} else {
		ba[3] += bb[3]
	}
	if bb[4] > 0xFF-ba[4] {
		ba[4] = 0xFF
	} else {
		ba[4] += bb[4]
	}
	if bb[5] > 0xFF-ba[5] {
		ba[5] = 0xFF
	} else {
		ba[5] += bb[5]
	}
	if bb[6] > 0xFF-ba[6] {
		ba[6] = 0xFF
	} else {
		ba[6] += bb[6]
	}
	if bb[7] > 0xFF-ba[7] {
		ba[7] = 0xFF
	} else {
		ba[7] += bb[7]
	}
	return ba
}

func selectByLowBits(ba, bb, ma [8]byte) [8]byte {
	var out [8]byte
	if ma[0] != 0 {
		out[0] = ba[0]
	} else {
		out[0] = bb[0]
	}
	if ma[1] != 0 {
		out[1] = ba[1]
	} else {
		out[1] = bb[1]
	}
	if ma[2] != 0 {
		out[2] = ba[2]
	} else {
		out[2] = bb[2]
	}
	if ma[3] != 0 {
		out[3] = ba[3]
	} else {
		out[3] = bb[3]
	}
	if ma[4] != 0 {
		out[4] = ba[4]
	} else {
		out[4] = bb[4]
	}
	if ma[5] != 0 {
		out[5] = ba[5]
	} else {
		out[5] = bb[5]
	}
	if ma[6] != 0 {
		out[6] = ba[6]
	} else {
		out[6] = bb[6]
	}
	if ma[7] != 0 {
		out[7] = ba[7]
	} else {
		out[7] = bb[7]
	}
	return out
}

func subBytesWrap(aa, bb [8]byte) [8]byte {
	aa[0] = aa[0] - bb[0]
	aa[1] = aa[1] - bb[1]
	aa[2] = aa[2] - bb[2]
	aa[3] = aa[3] - bb[3]
	aa[4] = aa[4] - bb[4]
	aa[5] = aa[5] - bb[5]
	aa[6] = aa[6] - bb[6]
	aa[7] = aa[7] - bb[7]
	return aa
}

func subBytesSat(aa, bb [8]byte) [8]byte {
	if bb[0] > aa[0] {
		aa[0] = 0
	} else {
		aa[0] -= bb[0]
	}
	if bb[1] > aa[1] {
		aa[1] = 0
	} else {
		aa[1] -= bb[1]
	}
	if bb[2] > aa[2] {
		aa[2] = 0
	} else {
		aa[2] -= bb[2]
	}
	if bb[3] > aa[3] {
		aa[3] = 0
	} else {
		aa[3] -= bb[3]
	}
	if bb[4] > aa[4] {
		aa[4] = 0
	} else {
		aa[4] -= bb[4]
	}
	if bb[5] > aa[5] {
		aa[5] = 0
	} else {
		aa[5] -= bb[5]
	}
	if bb[6] > aa[6] {
		aa[6] = 0
	} else {
		aa[6] -= bb[6]
	}
	if bb[7] > aa[7] {
		aa[7] = 0
	} else {
		aa[7] -= bb[7]
	}
	return aa
}

func addBytesSat(aa, bb [8]byte) [8]byte {
	s0 := uint16(aa[0]) + uint16(bb[0])
	if s0 > 255 {
		aa[0] = 0xFF
	} else {
		aa[0] = byte(s0)
	}
	s1 := uint16(aa[1]) + uint16(bb[1])
	if s1 > 255 {
		aa[1] = 0xFF
	} else {
		aa[1] = byte(s1)
	}
	s2 := uint16(aa[2]) + uint16(bb[2])
	if s2 > 255 {
		aa[2] = 0xFF
	} else {
		aa[2] = byte(s2)
	}
	s3 := uint16(aa[3]) + uint16(bb[3])
	if s3 > 255 {
		aa[3] = 0xFF
	} else {
		aa[3] = byte(s3)
	}
	s4 := uint16(aa[4]) + uint16(bb[4])
	if s4 > 255 {
		aa[4] = 0xFF
	} else {
		aa[4] = byte(s4)
	}
	s5 := uint16(aa[5]) + uint16(bb[5])
	if s5 > 255 {
		aa[5] = 0xFF
	} else {
		aa[5] = byte(s5)
	}
	s6 := uint16(aa[6]) + uint16(bb[6])
	if s6 > 255 {
		aa[6] = 0xFF
	} else {
		aa[6] = byte(s6)
	}
	s7 := uint16(aa[7]) + uint16(bb[7])
	if s7 > 255 {
		aa[7] = 0xFF
	} else {
		aa[7] = byte(s7)
	}
	return aa
}

func absDiffBytes(aa, bb [8]byte) [8]byte {
	if aa[0] >= bb[0] {
		aa[0] = aa[0] - bb[0]
	} else {
		aa[0] = bb[0] - aa[0]
	}
	if aa[1] >= bb[1] {
		aa[1] = aa[1] - bb[1]
	} else {
		aa[1] = bb[1] - aa[1]
	}
	if aa[2] >= bb[2] {
		aa[2] = aa[2] - bb[2]
	} else {
		aa[2] = bb[2] - aa[2]
	}
	if aa[3] >= bb[3] {
		aa[3] = aa[3] - bb[3]
	} else {
		aa[3] = bb[3] - aa[3]
	}
	if aa[4] >= bb[4] {
		aa[4] = aa[4] - bb[4]
	} else {
		aa[4] = bb[4] - aa[4]
	}
	if aa[5] >= bb[5] {
		aa[5] = aa[5] - bb[5]
	} else {
		aa[5] = bb[5] - aa[5]
	}
	if aa[6] >= bb[6] {
		aa[6] = aa[6] - bb[6]
	} else {
		aa[6] = bb[6] - aa[6]
	}
	if aa[7] >= bb[7] {
		aa[7] = aa[7] - bb[7]
	} else {
		aa[7] = bb[7] - aa[7]
	}
	return aa
}

func highBitWhereLess(b [8]byte, c [8]byte) [8]byte {
	if b[0] < c[0] {
		b[0] = 0x80
	} else {
		b[0] = 0
	}
	if b[1] < c[1] {
		b[1] = 0x80
	} else {
		b[1] = 0
	}
	if b[2] < c[2] {
		b[2] = 0x80
	} else {
		b[2] = 0
	}
	if b[3] < c[3] {
		b[3] = 0x80
	} else {
		b[3] = 0
	}
	if b[4] < c[4] {
		b[4] = 0x80
	} else {
		b[4] = 0
	}
	if b[5] < c[5] {
		b[5] = 0x80
	} else {
		b[5] = 0
	}
	if b[6] < c[6] {
		b[6] = 0x80
	} else {
		b[6] = 0
	}
	if b[7] < c[7] {
		b[7] = 0x80
	} else {
		b[7] = 0
	}
	return b
}

func highBitWhereGreater(b [8]byte, c [8]byte) [8]byte {
	if b[0] > c[0] {
		b[0] = 0x80
	} else {
		b[0] = 0
	}
	if b[1] > c[1] {
		b[1] = 0x80
	} else {
		b[1] = 0
	}
	if b[2] > c[2] {
		b[2] = 0x80
	} else {
		b[2] = 0
	}
	if b[3] > c[3] {
		b[3] = 0x80
	} else {
		b[3] = 0
	}
	if b[4] > c[4] {
		b[4] = 0x80
	} else {
		b[4] = 0
	}
	if b[5] > c[5] {
		b[5] = 0x80
	} else {
		b[5] = 0
	}
	if b[6] > c[6] {
		b[6] = 0x80
	} else {
		b[6] = 0
	}
	if b[7] > c[7] {
		b[7] = 0x80
	} else {
		b[7] = 0
	}
	return b
}

func highBitWhereEqual(b [8]byte, c [8]byte) [8]byte {
	if b[0] == c[0] {
		b[0] = 0x80
	} else {
		b[0] = 0
	}
	if b[1] == c[1] {
		b[1] = 0x80
	} else {
		b[1] = 0
	}
	if b[2] == c[2] {
		b[2] = 0x80
	} else {
		b[2] = 0
	}
	if b[3] == c[3] {
		b[3] = 0x80
	} else {
		b[3] = 0
	}
	if b[4] == c[4] {
		b[4] = 0x80
	} else {
		b[4] = 0
	}
	if b[5] == c[5] {
		b[5] = 0x80
	} else {
		b[5] = 0
	}
	if b[6] == c[6] {
		b[6] = 0x80
	} else {
		b[6] = 0
	}
	if b[7] == c[7] {
		b[7] = 0x80
	} else {
		b[7] = 0
	}
	return b
}

func TestSWARFunctionsRef(t *testing.T) {
	for n := uint64(0); n < 0x_FF_FF_FF_FF_FF; n = (n*12 + 13) / 11 {
		nA := toBytes(n)
		if a, b := SwapByteHalves(n), swapNibbles(nA); a != fromBytes(b) {
			t.Errorf("SwapByteHalves(0x%016x) = 0x%016x; want 0x%016x", n, a, fromBytes(b))
		}
		if a, b := ReverseEachByte(n), reverseBits(nA); a != fromBytes(b) {
			t.Errorf("ReverseEachByte(0x%016x) = 0x%016x; want 0x%016x", n, a, fromBytes(b))
		}
		if a, b := CountOnesPerByte(n), popcountPerByte(nA); a != fromBytes(b) {
			t.Errorf("CountOnesPerByte(0x%016x) = 0x%016x; want 0x%016x", n, a, fromBytes(b))
		}
		if a, b := ExtractLowBits(n&LowBits), byteFromLowBits(toBytes(n&LowBits)); a != b {
			t.Errorf("ExtractLowBits(0x%016x) = 0b%08b; want 0b%08b", n, a, b)
		}

		m := n ^ 0x0000005351952b76
		mA := toBytes(m)
		if a, b := SelectSmallerBytes(n, m), minBytes(nA, mA); a != fromBytes(b) {
			t.Errorf("SelectSmallerBytes(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", n, m, a, fromBytes(b))
		}
		if a, b := SelectLargerBytes(n, m), maxBytes(nA, mA); a != fromBytes(b) {
			t.Errorf("SelectLargerBytes(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", n, m, a, fromBytes(b))
		}
		if a, b := AverageBytes(n, m), averageBytes(nA, mA); a != fromBytes(b) {
			t.Errorf("AverageBytes(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", n, m, a, fromBytes(b))
		}
		if a, b := AbsoluteDifferenceBetweenBytes(n, m), absDiffBytes(nA, mA); a != fromBytes(b) {
			t.Errorf("AbsoluteDifferenceBetweenBytes(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", n, m, a, fromBytes(b))
		}
		if a, b := AddBytesWithWrapping(n, m), addWrapBytes(nA, mA); a != fromBytes(b) {
			t.Errorf("AddWrapBytes(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", n, m, a, fromBytes(b))
		}
		if a, b := SubtractBytesWithWrapping(n, m), subBytesWrap(nA, mA); a != fromBytes(b) {
			t.Errorf("SubtractBytesWithWrapping(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", n, m, a, fromBytes(b))
		}
		if a, b := AddBytesWithMaximum(n, m), addSatBytes(nA, mA); a != fromBytes(b) {
			t.Errorf("AddBytesWithMaximum(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", n, m, a, fromBytes(b))
		}
		if a, b := SubtractBytesWithMinimum(n, m), subBytesSat(nA, mA); a != fromBytes(b) {
			t.Errorf("SubtractBytesWithMinimum(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", n, m, a, fromBytes(b))
		}

		c := Dupe(byte(m % 0x_FE))
		if a, b := HighBitWhereLess(n, c), highBitWhereLess(nA, toBytes(c)); a != fromBytes(b) {
			t.Errorf("HighBitWhereLess(0x%016x, %2x) = 0x%016x; want 0x%016x", n, c, a, fromBytes(b))
		}
		if a, b := HighBitWhereGreater(n, c), highBitWhereGreater(nA, toBytes(c)); a != fromBytes(b) {
			t.Errorf("HighBitWhereGreater(0x%016x, %2x) = 0x%016x; want 0x%016x", n, c, a, fromBytes(b))
		}
		if a, b := HighBitWhereEqual(n, c), highBitWhereEqual(nA, toBytes(c)); a != fromBytes(b) {
			t.Errorf("HighBitWhereEqual(0x%016x, %2x) = 0x%016x; want 0x%016x", n, c, a, fromBytes(b))
		}

		d := uint64(0x_01_00_01_01_00_00_01_00)
		dA := toBytes(d)
		if a, b := SelectByLowBit(n, m, d), selectByLowBits(nA, mA, dA); a != fromBytes(b) {
			t.Errorf("SelectByLowBit(0x%016x, 0x%016x, 0x%016x) = 0x%016x; want 0x%016x (%v)", n, m, d, a, fromBytes(b), dA)
		}

		// t.Logf("Tested with n=0x%016x, m=0x%016x, c=%02x, d=0x%016x", n, m, c, d)
		// t.Logf("As arrays: n=%v, m=%v, d=%v", nA, mA, dA)
	}
}
