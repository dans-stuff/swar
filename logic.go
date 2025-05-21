package swar

const (
	HighBits uint64 = 0x8080_8080_8080_8080
)

// HighBitWhereLess:  0x80 in each lane where v < cm
func HighBitWhereLess(v, cm uint64) uint64 {
	d := (v | HighBits) - (cm &^ HighBits)
	sel := ((v & (v ^ cm)) | (d &^ (v ^ cm))) & HighBits
	hbit := sel ^ HighBits // 0x80 in each byte where v < cm
	return hbit & HighBits // already 0x80 or 0x00 per lane
}

// HighBitWhereGreater: 0x80 in each lane where v > cm
func HighBitWhereGreater(v, cm uint64) uint64 {
	d := (cm | HighBits) - (v &^ HighBits)
	sel := ((cm & (cm ^ v)) | (d &^ (cm ^ v))) & HighBits
	hbit := sel ^ HighBits // 0x80 in each byte where v > cm
	return hbit & HighBits // 0x80 or 0x00 per lane
}

// HighBitWhereEqual: 0x80 in each lane where v == cm
func HighBitWhereEqual(v, cm uint64) uint64 {
	x := v ^ cm
	y := ((x & 0x7F7F7F7F7F7F7F7F) + 0x7F7F7F7F7F7F7F7F) | x
	hi := ^y & HighBits  // 0x80 where x==0 ⇔ v==cm
	return hi & HighBits // mask off any other bits
}
