package gt06

import (
	"math"
	"reciever-ms/utils/arrays"
	"reflect"
	"testing"
)

const float64EqualityThreshold = 1e-8

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestValidateAndSetMessageSize(t *testing.T) {
	t.Run("returns error if message is too small to contain packet length", func(t *testing.T) {
		testCases := [][]byte{
			nil,
			{1},
			{1, 2},
		}

		for i := 0; i < len(testCases); i++ {
			m := Gt06MsgMeta{Bytes: testCases[i]}
			if err := m.validateAndSetMessageSize(); err == nil {
				t.Fatal("did not throw error for message with less than 3 bytes")
			}
		}
	})

	t.Run("sets the InfoContentCount to the value the third byte - 5", func(t *testing.T) {
		var packetLen uint8 = 10

		m := Gt06MsgMeta{Bytes: []byte{1, 1, packetLen}}
		m.validateAndSetMessageSize()

		if m.InfoContentByteCnt != int(packetLen)-5 {
			t.Fatal("did not read the value in the third byte as the packet lenght")
		}
	})

	t.Run("sets the TotalByteCnt to the value the third byte + 5", func(t *testing.T) {
		var packetLen uint8 = 20

		m := Gt06MsgMeta{Bytes: []byte{1, 1, packetLen}}
		m.validateAndSetMessageSize()

		if m.TotalByteCnt != int(packetLen)+5 {
			t.Fatal("did not read the value in the third byte as the packet lenght")
		}
	})

	t.Run("returns error if message size is less then the informed on TotalByteCnt", func(t *testing.T) {
		var packetLen uint8 = 20

		m := Gt06MsgMeta{Bytes: []byte{1, 1, packetLen}}
		err := m.validateAndSetMessageSize()

		if err == nil {
			t.Fatal("should return error as message has only 3 bytes but total bytes informed are 25")
		}
	})

	t.Run("succeds on correct packetLen", func(t *testing.T) {
		var packetLen uint8 = 20

		m := Gt06MsgMeta{
			Bytes: []byte{1, 2, packetLen, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25},
		}
		err := m.validateAndSetMessageSize()

		if err != nil {
			t.Fatalf("should return no error: got %v", err)
		}
	})
}

func TestGetBytesForCrc(t *testing.T) {
	t.Run("returns all bytes in packetLen, protocolNum, terminalId and infoSerialNum", func(t *testing.T) {
		m := Gt06MsgMeta{Bytes: loginMessage}

		expected := arrays.ConcatAppend([][]byte{
			packetLen,
			protocolNum,
			terminalId,
			infoSerialNum,
		})

		bytes, _ := m.getBytesForCrc()

		if !reflect.DeepEqual(bytes, expected) {
			t.Fatal("returned invalid bytes for crc check")
		}
	})

	t.Run("fails on insuficient bytes", func(t *testing.T) {
		m := Gt06MsgMeta{Bytes: nil}

		if _, err := m.getBytesForCrc(); err == nil {
			t.Fatal("should fail on msg with no bytes")
		}
	})
}

func TestValidateSelf(t *testing.T) {
	t.Run("succeds on valid messages", func(t *testing.T) {
		m := Gt06MsgMeta{Bytes: loginMessage}
		if err := m.validateSelf(); err != nil {
			t.Fatal("validate failed for valid login packets")
		}
	})

	t.Run("fails on invalid bytes", func(t *testing.T) {
		m := Gt06MsgMeta{Bytes: []byte{1, 2, 3}}
		if err := m.validateSelf(); err == nil {
			t.Fatal("validate should fail on nonsensical bytes")
		}
	})

	t.Run("fails on crc failure", func(t *testing.T) {
		invalidTerminalId := make([]byte, len(terminalId))
		copy(invalidTerminalId, terminalId)

		// make the first byte invalid
		invalidTerminalId[0] = terminalId[0] + 1

		invalidLoginMessage := arrays.ConcatAppend([][]byte{
			startBit,
			packetLen,
			protocolNum,
			invalidTerminalId,
			infoSerialNum,
			errCheck,
			stopBit,
		})

		m := Gt06MsgMeta{Bytes: invalidLoginMessage}

		if err := m.validateSelf(); err == nil {
			t.Fatal("validate should fail messages that fail the crc")
		}
	})
}

func TestNewMsg(t *testing.T) {
	msg, err := NewMsg(loginMessage)
	if err != nil {
		t.Fatalf("unknonw error on valid message: %v", err)
	}

	if !reflect.DeepEqual(msg.ErrorCheck, errCheck) {
		t.Fatal("ErrorCheck bytes do not match")
	}

	if !reflect.DeepEqual(msg.InformationContent, terminalId) {
		t.Fatal("InformationContent bytes do not match")
	}

	if msg.PacketLenght[0] != packetLen[0] {
		t.Fatal("PacketLenght bytes do not match")
	}

	if msg.ProtocolNumber[0] != protocolNum[0] {
		t.Fatal("ProtocolNumber bytes do not match")
	}

	if !reflect.DeepEqual(msg.InformationSerialNumber, infoSerialNum) {
		t.Fatal("InformationSerialNumber bytes do not match")
	}
}
