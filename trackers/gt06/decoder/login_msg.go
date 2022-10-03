package decoder

import (
	"reciever-ms/utils/arrays"
	"reciever-ms/utils/crc"
	dstrings "reciever-ms/utils/strings"
	"strings"
)

type LoginRes struct {
	Imei string `json:"imei"`
}

func getImei(b []byte) string {
	s := dstrings.BytesAsLiteralString(b)
	return strings.TrimLeft(s, "0")
}

func (m *Gt06Msg) DecodeLogin() DecodeRes {
	// always 5 in the login response
	PACKET_LEN := []byte{0x05}

	crcBytes := arrays.ConcatAppend([][]byte{
		PACKET_LEN,
		m.ProtocolNumber,
		m.InformationSerialNumber,
	})

	errorCheck := crc.Crc16ItuToByteArr(crcBytes)

	response := arrays.ConcatAppend([][]byte{
		START_BIT,
		PACKET_LEN,
		m.ProtocolNumber,
		m.InformationSerialNumber,
		errorCheck,
		STOP_BIT,
	})

	return DecodeRes{
		Err:     nil,
		Res:     response,
		Msg:     LoginRes{Imei: getImei(m.InformationContent)},
		MsgType: "LoginRes",
	}
}
