package swar

import (
	"testing"
)

// TestHighBitWhereEqual verifies that the HighBitWhereEqual function correctly
// identifies bytes that match a comparison value. These tests are important because
// the SWAR technique uses non-intuitive bit manipulation that needs proper verification.
func TestHighBitWhereEqual(t *testing.T) {
	run := func(v, c, want uint64) {
		if got := HighBitWhereEqual(v, c); got != want {
			t.Errorf("HighBitWhereEqual(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", v, c, got, want)
		}
	}

	run(0x05, Dupe(5), 0x80)
	run(0x04, Dupe(5), 0x00)
	run(0x05_04, Dupe(5), 0x80_00)
	run(0xFF_00, Dupe(0), 0x80_80_80_80_80_80_00_80)
}

// TestHighBitWhereLess verifies that the HighBitWhereLess function correctly identifies
// bytes less than a comparison value. This is crucial for threshold-based processing
// and range checks operating on multiple bytes in parallel.
func TestHighBitWhereLess(t *testing.T) {
	run := func(v, c, want uint64) {
		if got := HighBitWhereLess(v, c); got != want {
			t.Errorf("HighBitWhereLess(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", v, c, got, want)
		}
	}

	run(0x06, Dupe(5), 0x80_80_80_80_80_80_80_00)
	run(0x04, Dupe(5), 0x80_80_80_80_80_80_80_80)
	run(0x01_02_03_04_05_06_07_08, Dupe(5), 0x80_80_80_80_00_00_00_00)
}

// TestHighBitWhereGreater verifies that the HighBitWhereGreater function correctly
// identifies bytes greater than a comparison value. This functionality is essential for
// detecting outliers, anomalies, and values exceeding specified thresholds.
func TestHighBitWhereGreater(t *testing.T) {
	run := func(v, c, want uint64) {
		if got := HighBitWhereGreater(v, c); got != want {
			t.Errorf("HighBitWhereGreater(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", v, c, got, want)
		}
	}

	run(0x05, Dupe(5), 0x00)
	run(0x06, Dupe(5), 0x80)
	run(0xFF_04_05_06_00, Dupe(5), 0x80_00_00_80_00)
}

// TestSelectByLowBit verifies that values are correctly selected from a or b based on
// the corresponding mask bit. This branchless selection is critical for data-dependent
// operations where conditional logic would otherwise harm performance.
func TestSelectByLowBit(t *testing.T) {
	run := func(a, b, mask, want uint64) {
		if got := SelectByLowBit(a, b, mask); got != want {
			t.Errorf("SelectByLowBit(0x%016x, 0x%016x, 0x%016x) = 0x%016x; want 0x%016x", a, b, mask, got, want)
		}
	}

	run(0x11_11_11_11, 0x22_22_22_22, 0x01_00_01_00, 0x11_22_11_22)
}

// TestMinBytes verifies that our parallel minimum function correctly selects the smaller
// of two values for each byte position. This is essential for applications like image processing
// where per-pixel minimum operations affect visual outcomes.
func TestMinBytes(t *testing.T) {
	run := func(a, b, want uint64) {
		if got := SelectSmallerBytes(a, b); got != want {
			t.Errorf("SelectSmallerBytes(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", a, b, got, want)
		}
	}

	run(0x01_02_03_04_05_06_07_08, 0x05_04_03_02_01_00_09_0A, 0x01_02_03_02_01_00_07_08)
	run(0x0000_0000_0000_0004, 0x1234_5678_90AB_CDEB, 0x0000_0000_0000_0004)
}

// TestMaxBytes verifies that our parallel maximum function correctly selects the larger
// of two values for each byte position. This is critical for algorithms like feature extraction
// and signal peak detection where maintaining maximum values is required.
func TestMaxBytes(t *testing.T) {
	run := func(a, b, want uint64) {
		if got := SelectLargerBytes(a, b); got != want {
			t.Errorf("SelectLargerBytes(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", a, b, got, want)
		}
	}

	run(0x01_02_03_04_05_06_07_08, 0x05_04_03_02_01_00_09_0A, 0x05_04_03_04_05_06_09_0A)
	run(0x04, 0xEB, 0xEB)
	run(0x01, 0x02, 0x02)
}

// TestSwapNibbles verifies that our nibble-swapping function correctly exchanges the high
// and low 4 bits of each byte. This transformation is important for BCD encoding/decoding
// and certain data format conversions that rely on nibble-level manipulations.
func TestSwapNibbles(t *testing.T) {
	run := func(s, want uint64) {
		if got := SwapByteHalves(s); got != want {
			t.Errorf("SwapByteHalves(0x%016x) = 0x%016x; want 0x%016x", s, got, want)
		}
	}

	run(0xF0_F0_F0_F0_F0_F0_F0_F0, 0x0F_0F_0F_0F_0F_0F_0F_0F)
}

// TestReverseBits verifies that our bit-reversal function correctly reverses the order
// of bits within each byte. This is crucial for operations like endianness conversion
// and certain data transformations that depend on bit-level mirroring.
func TestReverseBits(t *testing.T) {
	run := func(v, want uint64) {
		if got := ReverseEachByte(v); got != want {
			t.Errorf("ReverseEachByte(0x%016x) = 0x%016x; want 0x%016x", v, got, want)
		}
	}

	run(0x01_02_04_08_10_20_40_80, 0x80_40_20_10_08_04_02_01)
	run(0b01001000_11100001_11000011_11110000, 0b00010010_10000111_11000011_00001111)
}

// TestPopcountPerByte verifies that our parallel population count correctly counts the
// set bits in each byte. This functionality is essential for feature extraction, hamming
// distance calculation, and statistical analysis of binary data.
func TestPopcountPerByte(t *testing.T) {
	run := func(p, want uint64) {
		if got := CountOnesPerByte(p); got != want {
			t.Errorf("CountOnesPerByte(0x%016x) = 0x%016x; want 0x%016x", p, got, want)
		}
	}

	run(0x0F_F0_55_AA_00_FF_33_CC, 0x04_04_04_04_00_08_04_04)
}
