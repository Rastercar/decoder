package strings

import (
	"encoding/binary"
	"fmt"
	"strings"
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

// GetStringInBetween Returns empty string if no start string found
func GetStringInBetween(str, start, end string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}

	s += len(start)

	e := strings.Index(str[s:], end)
	if e == -1 {
		return ""
	}

	e = s + e

	return str[s:e]
}
