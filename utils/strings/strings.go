package strings

import (
	"encoding/binary"
	"fmt"
)

// eg: []byte{1, 0xA2, 0x33, 10} -> 01A2330A
func BytesAsLiteralString(b []byte) string {
	s := ""

	for i := 0; i < len(b); i++ {
		s += fmt.Sprintf("%02X", b[i])
	}

	return s
}

func BytesToDecimalString(b []byte) string {
	return fmt.Sprintf("%d", binary.BigEndian.Uint16(b))
}
