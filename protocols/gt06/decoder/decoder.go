package decoder

import (
	"errors"
	"fmt"
)

var (
	START_BIT = []byte{0x78, 0x78}
	STOP_BIT  = []byte{0x0D, 0x0A}
)

// GT06 Protocol numbers

var (
	LOGIN_MESSAGE                     = uint8(0x01)
	LOCATION_DATA                     = uint8(0x12)
	STATUS_INFORMATION                = uint8(0x13)
	STRING_INFORMATION                = uint8(0x15)
	ALARM_DATA                        = uint8(0x16)
	GPS_QUERY_ADDRESS_BY_PHONE_NUMBER = uint8(0x1A)
	CMD_INFORMATION                   = uint8(0x80)
)

type decoder struct {
	debugMode bool
}

type DecodeRes struct {
	Err     error       // the decoding error
	Res     []byte      // bytes to respond to the tracker through tcp/udp
	Msg     interface{} // nil or a serializable struct with data about the decode response
	MsgType string      // name of the struct in msg
}

func NewDecoder(debug bool) decoder {
	return decoder{debugMode: debug}
}

func (d *decoder) printfIfDebug(s string, a ...any) {
	if d.debugMode {
		fmt.Printf(s, a...)
	}
}

func (d *decoder) err(s string, a ...any) error {
	ss := fmt.Sprintf(s, a...)
	d.printfIfDebug(ss)
	return errors.New(ss)
}

func (d *decoder) Decode(msg []byte) DecodeRes {
	d.printfIfDebug("\ndecoding GT06 message")

	m, err := NewMsg(msg)
	if err != nil {
		return DecodeRes{Err: d.err("decoding failed: %v", err)}
	}

	var res DecodeRes

	switch m.ProtocolNumber[0] {
	case LOGIN_MESSAGE:
		res = m.DecodeLogin()
	case LOCATION_DATA:
		res = m.DecodeLocation()
	default:
		res = DecodeRes{Err: d.err("cannot decode msg, unkown protocol %d", m.ProtocolNumber)}
	}

	if res.Err != nil {
		d.printfIfDebug("\ndecoding failed: %v", err)
	}

	return res
}
