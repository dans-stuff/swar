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
