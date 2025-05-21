package swar

const (
	// HighBits is a mask with the high bit set in all 8 bytes of a uint64
	HighBits uint64 = 0x8080_8080_8080_8080
)

// HighBitWhereLess sets the high bit (0x80) in each byte where v < cm
// Enables parallel comparison of 8 bytes simultaneously
func HighBitWhereLess(v, cm uint64) uint64 {
	d := (v | HighBits) - (cm &^ HighBits)
	sel := ((v & (v ^ cm)) | (d &^ (v ^ cm))) & HighBits
	hbit := sel ^ HighBits // 0x80 in each byte where v < cm
	return hbit & HighBits // 0x80 or 0x00 per lane
}

// HighBitWhereGreater sets the high bit (0x80) in each byte where v > cm
// Perfect for threshold detection across multiple values
func HighBitWhereGreater(v, cm uint64) uint64 {
	d := (cm | HighBits) - (v &^ HighBits)
	sel := ((cm & (cm ^ v)) | (d &^ (cm ^ v))) & HighBits
	hbit := sel ^ HighBits // 0x80 in each byte where v > cm
	return hbit & HighBits // 0x80 or 0x00 per lane
}

// HighBitWhereEqual sets the high bit (0x80) in each byte where v == cm
// Ideal for pattern matching and finding specific values in data
func HighBitWhereEqual(v, cm uint64) uint64 {
	x := v ^ cm
	y := ((x & 0x7F7F7F7F7F7F7F7F) + 0x7F7F7F7F7F7F7F7F) | x
	hi := ^y & HighBits  // 0x80 where x==0 (v==cm)
	return hi & HighBits // mask off other bits
}
