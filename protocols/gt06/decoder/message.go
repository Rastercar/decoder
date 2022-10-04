package decoder

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reciever-ms/utils/crc"
)

var BYTES_NEEDED_TO_READ_PACKET_LEN = 3

type Gt06MsgMeta struct {
	Bytes              []byte // all bytes of the original message
	TotalByteCnt       int    // 10 + InfoContentByteCount
	InfoContentByteCnt int    //
}

type Gt06Msg struct {
	PacketLenght            []byte // 1 byte
	ProtocolNumber          []byte // 1 byte
	InformationContent      []byte // N bytes (N is PacketLenght)
	InformationSerialNumber []byte // 2 bytes
	ErrorCheck              []byte // 2 bytes
}

func (d *Gt06MsgMeta) validateAndSetMessageSize() error {
	originalMsgSize := len(d.Bytes)

	if originalMsgSize < BYTES_NEEDED_TO_READ_PACKET_LEN {
		return errors.New("message size is not enough to read packet lenght")
	}

	infoContentCnt := int(d.Bytes[2] - 5)

	// sum of every part of the protocol except infoContent is 10 bytes
	d.TotalByteCnt = 10 + infoContentCnt
	d.InfoContentByteCnt = infoContentCnt

	if originalMsgSize < d.TotalByteCnt {
		return fmt.Errorf("msg size is lower then the expected total bytes, got: %d - expected %d", originalMsgSize, d.TotalByteCnt)
	}

	return nil
}

// Returns the bytes of the message that are used to generate the checksum for the CRC
func (m *Gt06MsgMeta) getBytesForCrc() ([]byte, error) {
	if err := m.validateAndSetMessageSize(); err != nil {
		return nil, err
	}

	// start at 2 since bytes 0 and 1 are the "start bit"
	packetLenStartIdx := 2

	// packetLen + protocolNum + infoContentCnt + infoSerialNumber = 6
	infoSerialNumberEndIdx := 6 + m.InfoContentByteCnt

	return m.Bytes[packetLenStartIdx:infoSerialNumberEndIdx], nil
}

func (m *Gt06MsgMeta) validateSelf() error {
	crcBytes, err := m.getBytesForCrc()
	if err != nil {
		return err
	}

	errorCheckStartIdx := 6 + m.InfoContentByteCnt

	errorCheckBytes := m.Bytes[errorCheckStartIdx : errorCheckStartIdx+2]

	checkSum := crc.Crc16Itu(crcBytes)
	expected := binary.BigEndian.Uint16(errorCheckBytes)

	if checkSum != expected {
		return fmt.Errorf("checksum missmatch, got: %d - expected: %d", checkSum, expected)
	}

	return nil
}

func NewMsg(packets []byte) (*Gt06Msg, error) {
	m := Gt06MsgMeta{Bytes: packets}

	if err := m.validateSelf(); err != nil {
		return nil, err
	}

	n := m.InfoContentByteCnt

	return &Gt06Msg{
		PacketLenght:            []byte{packets[2]},
		ProtocolNumber:          []byte{packets[3]},
		InformationContent:      packets[4 : 4+n],
		InformationSerialNumber: packets[4+n : 6+n],
		ErrorCheck:              packets[6+n : 8+n],
	}, nil
}
