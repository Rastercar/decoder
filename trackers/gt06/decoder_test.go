package gt06

import (
	"bytes"
	"testing"
)

func TestDecodeLogin(t *testing.T) {
	// Taken from Gt06 protocol docs
	loginMsg := []byte{
		0x78, 0x78,
		0x0D,
		0x01,
		0x01, 0x23, 0x45, 0x67, 0x89, 0x01, 0x23, 0x45,
		0x00, 0x01,
		0x8C, 0xDD,
		0x0D, 0x0A,
	}

	expectedRes := []byte{0x78, 0x78, 0x05, 0x01, 0x00, 0x01, 0xD9, 0xDC, 0x0D, 0x0A}

	m, err := NewMsg(loginMsg)
	if err != nil {
		t.Fatal(err)
	}

	res := m.DecodeLogin()
	if res.Err != nil {
		t.Fatal(err)
	}

	if res.MsgType != "LoginRes" {
		t.Fatal("invalid msg type")
	}

	_, ok := res.Msg.(LoginRes)
	if !ok {
		t.Fatal("decoded login message did not cast to its type")
	}

	if !bytes.Equal(res.Res, expectedRes) {
		t.Fatal("response does not match expected !")
	}
}
