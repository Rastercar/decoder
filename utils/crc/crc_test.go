package crc

import (
	"encoding/binary"
	"testing"
)

func TestCrc16Itu(t *testing.T) {
	// taken from gt06 login packet on the docs
	loginPacket := []byte{
		0x0D,                                           // packet lenght
		0x01,                                           // protocol number
		0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45, // terminal ID
		0x00, 0x01, // information serial number
	}

	errorCheck := []byte{0x8C, 0xDD}

	if Crc16Itu(loginPacket) != binary.BigEndian.Uint16(errorCheck) {
		t.Fatal("Invalid CRC16")
	}
}
