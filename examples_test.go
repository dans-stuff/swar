package swar

import (
	"bytes"
	"math/bits"
	"testing"
)

// Sample text for benchmark tests - contains spaces and mixed case for testing
var lotsOfBytes = []byte("Allo Zorld! I am NOT yelling, but I am using SWAR!")

// BenchmarkUsageCount compares the performance of counting spaces using traditional
// byte-by-byte scanning versus SWAR-based parallel comparison. This benchmark 
// demonstrates how SIMD-within-a-register can accelerate simple character counting,
// which is useful in text processing applications.
func BenchmarkUsageCount(b *testing.B) {
	b.Run("BestNaive", func(b *testing.B) {
		count := 0
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, c := range lotsOfBytes {
				if c == ' ' {
					count++
				}
			}
		}
		if count != 10*b.N {
			b.Errorf("Expected %d, got %d", 10*b.N, count)
		}
	})

	b.Run("SWAR", func(b *testing.B) {
		spaces := Dupe(' ')
		count := 0
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			chunks, unused := BytesToLanes(lotsOfBytes)
			for _, chunk := range chunks {
				count += bits.OnesCount64(HighBitWhereEqual(chunk, spaces))
			}
			for _, c := range lotsOfBytes[unused:] {
				if c == ' ' {
					count++
				}
			}
		}
		if count != 10*b.N {
			b.Errorf("Expected %d, got %d", 10*b.N, count)
		}
	})
}

// BenchmarkUsageVisitCaps compares traditional and SWAR approaches for finding and 
// processing uppercase letters in text. This benchmark demonstrates how SWAR enables 
// efficient filtering and position tracking in parallel, which is valuable for 
// text analysis and pattern matching applications.
func BenchmarkUsageVisitCaps(b *testing.B) {
	b.Run("BestNaive", func(b *testing.B) {
		sum := 0
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for i, c := range lotsOfBytes {
				if c >= 'A' && c <= 'Z' {
					sum += i
				}
			}
		}
		if sum/b.N != 291 {
			b.Errorf("Expected 291, got %d", sum/b.N)
		}
	})

	b.Run("SWAR", func(b *testing.B) {
		firstCapital, lastCapital := Dupe('A'-1), Dupe('Z'+1)
		sum := 0
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			chunks, unused := BytesToLanes(lotsOfBytes)

			for idx, chunk := range chunks {
				caps := HighBitWhereGreater(chunk, firstCapital) & HighBitWhereLess(chunk, lastCapital)
				matches := ExtractLowBits(caps >> 7)
				offsets := Lookup.OnesPositions[matches]
				for _, v := range offsets {
					sum += v + idx*8
				}
			}

			for i, c := range lotsOfBytes[unused:] {
				if c >= 'A' && c <= 'Z' {
					sum += unused + i
				}
			}
		}

		if sum/b.N != 291 {
			b.Errorf("Expected 291, got %d", sum/b.N)
		}
	})
}

// BenchmarkUsageUppercase compares standard library and SWAR approaches to converting 
// text to uppercase. This benchmark shows how SWAR enables high-performance text 
// transformation by applying character-level changes to multiple bytes in parallel,
// which is important for text processing pipelines.
func BenchmarkUsageUppercase(b *testing.B) {
	b.Run("BestNaive", func(b *testing.B) {

		out := []byte{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			out = bytes.ToUpper(lotsOfBytes)
		}
		if string(out) != "ALLO ZORLD! I AM NOT YELLING, BUT I AM USING SWAR!" {
			b.Errorf("Expected 'ALLO ZORLD! I AM NOT YELLING, BUT I AM USING SWAR!', got '%s'", string(out))
		}
	})

	b.Run("SWAR", func(b *testing.B) {
		out := make([]byte, len(lotsOfBytes))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			inc := Dupe(32)
			firstLower, lastLower := Dupe('a'-1), Dupe('z'+1)

			out = make([]byte, len(lotsOfBytes))
			outLanes, _ := BytesToLanes(out)
			chunks, unused := BytesToLanes(lotsOfBytes)
			for idx, chunk := range chunks {
				lowercases := HighBitWhereGreater(chunk, firstLower) & HighBitWhereLess(chunk, lastLower)
				allUpper := SubtractBytesWithWrapping(chunk, inc)
				outLanes[idx] = SelectByLowBit(allUpper, chunk, lowercases>>7)
			}
			for i, c := range lotsOfBytes[unused:] {
				if c >= 'a' && c <= 'z' {
					out[unused+i] = c - 32
				} else {
					out[unused+i] = c
				}
			}
		}
		if string(out) != "ALLO ZORLD! I AM NOT YELLING, BUT I AM USING SWAR!" {
			b.Errorf("Expected 'ALLO ZORLD! I AM NOT YELLING, BUT I AM USING SWAR!', got '%s'", string(out))
		}
	})
}

// BenchmarkUsageAnomalies demonstrates using SWAR for anomaly detection in time series data.
// This benchmark shows how SWAR enables efficient detection of unusual patterns or outliers
// by processing multiple values simultaneously and using parallel threshold comparison,
// which is critical for real-time monitoring and alerting systems.
func BenchmarkUsageAnomalies(b *testing.B) {

	b.Run("BestNaive", func(b *testing.B) {
		currentTemps := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		averageTemps := []byte{0, 0, 0, 0, 0, 0, 0, 0}
		anomalies := 0
		threshold := int(2) // going above 2 in one step is an anomaly

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			currentTemps[i%8] += 1 // normal behaviour
			if i%81 == 0 {
				currentTemps[i%8] = 0 // simulate an anomaly
			}

			for j, avg := range averageTemps {
				averageTemps[j] = byte((int(avg) + int(currentTemps[j])) / 2)
				delta := int(currentTemps[j]) - int(avg)
				if delta > threshold || delta < -threshold {
					averageTemps[j] = currentTemps[j]
					anomalies++
				}
			}
		}
		if anomalies != b.N/81 {
			b.Errorf("Expected %d (%d/81) anomalies, got %d", b.N/81, b.N, anomalies)
		}
	})

	b.Run("SWAR", func(b *testing.B) {
		currentTemps := []byte{0, 0, 0, 0, 0, 0, 0, 0}

		currentLane, _ := BytesToLanes(currentTemps)
		averageTemps, _ := BytesToLanes([]byte{0, 0, 0, 0, 0, 0, 0, 0})
		anomalies := 0
		threshold := Dupe(2) // going above 2 in one step is an anomaly

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			currentTemps[i%8] += 1 // normal behaviour
			if i%81 == 0 {
				currentTemps[i%8] = 0 // simulate an anomaly
			}

			averageTemps[0] = AverageBytes(currentLane[0], averageTemps[0])
			delta := AbsoluteDifferenceBetweenBytes(currentLane[0], averageTemps[0])
			overThreshold := HighBitWhereGreater(delta, threshold)
			if overThreshold != 0 {
				averageTemps[0] = currentLane[0]
				anomalies++
			}
		}
		if anomalies != b.N/81 {
			b.Errorf("Expected %d (%d/81) anomalies, got %d", b.N/81, b.N, anomalies)
		}
	})
}
