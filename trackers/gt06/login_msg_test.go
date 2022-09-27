package gt06

import (
	"reciever-ms/utils/arrays"
	"testing"
)

var (
	startBit      = []byte{0x78, 0x78}
	packetLen     = []byte{0x0D}
	protocolNum   = []byte{0x01}
	terminalId    = []byte{0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45}
	infoSerialNum = []byte{0x00, 0x01}
	errCheck      = []byte{0x8C, 0xDD}
	stopBit       = []byte{0x0D, 0x0A}
)

var loginMessage = arrays.ConcatAppend([][]byte{
	startBit,
	packetLen,
	protocolNum,
	terminalId,
	infoSerialNum,
	errCheck,
	stopBit,
})

func TestGetImei(t *testing.T) {
	i := []byte{0x0f, 0xa}

	if getImei(i) != "F0A" {
		t.Fatal()
	}
}
