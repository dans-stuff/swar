package swar

import (
	"testing"
)

// TestAverageBytes verifies that our parallel averaging algorithm correctly calculates
// the mean of corresponding bytes. This ensures proper data smoothing and interpolation
// behavior when processing multiple values simultaneously.
func TestAverageBytes(t *testing.T) {
	run := func(a, b, want uint64) {
		if got := AverageBytes(a, b); got != want {
			t.Errorf("AverageBytes(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", a, b, got, want)
		}
	}

	run(0x01_10_40_FF, 0xFF_30_80_FD, 0x80_20_60_FE)
	run(0x04, 0x08, 0x06)
	run(0x10_DD, 0x30_FF, 0x20_EE)
	run(0x0004, 0xCDEB, 0x6677)
	run(0x01FE, 0xCC11, 0x6687)
}

// TestAddSatBytes verifies that our saturating addition correctly clamps results to 0xFF
// when overflow occurs. This is crucial for applications like image processing and signal
// manipulation where preventing overflow is necessary for correct results.
func TestAddSatBytes(t *testing.T) {
	run := func(a, b, want uint64) {
		if got := AddBytesWithMaximum(a, b); got != want {
			t.Errorf("AddBytesWithMaximum(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", a, b, got, want)
		}
	}

	run(0xFF_FE_FD, 0x01_01_01, 0xFF_FF_FE)
	run(0xFD_FC_FB, 0x03_03_03, 0xFF_FF_FE)
}

// TestAddBytesWithWrapping ensures that our wrapping addition correctly handles overflow
// by wrapping around to zero. This behavior is essential for certain algorithms like
// checksums and hash functions where wrap-around arithmetic is expected and required.
func TestAddBytesWithWrapping(t *testing.T) {
	run := func(a, b, want uint64) {
		if got := AddBytesWithWrapping(a, b); got != want {
			t.Errorf("AddBytesWithWrapping(0x%016x, 0x%016x) = 0x%016x; want 0x%016x", a, b, got, want)
		}
	}

	run(0xFF_FE_FD, 0x01_01_01, 0x00_FF_FE)
	run(0xFD_FC_FB, 0x03_03_03, 0x00_FF_FE)
	run(0xF4_F9, 0x0F_01, 0x03_FA)
	run(0xFF_0F_FF, 0x01_F0_00, 0x00_FF_FF)
}
