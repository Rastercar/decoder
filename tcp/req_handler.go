package tcp

import (
	"fmt"
	"net"
	"reciever-ms/trackers/gt06"
)

var decoder = gt06.NewDecoder(true)

var MAX_INVALID_MESSAGES_PER_CONN = 10

func castDecodeRes[T any](d *gt06.DecodeRes) (*T, error) {
	cast, ok := d.Msg.(T)
	if !ok {
		return nil, fmt.Errorf("failed to cast message of type %s", d.MsgType)
	}

	return &cast, nil
}

func handleDecodedMsg(d *gt06.DecodeRes) {
	switch d.MsgType {
	case "LocationRes":
		// TODO:
	}
}

func HandleRequest(c net.Conn) {
	// the imei of the tracker in the current connection
	// set once the tracker sends a login packet with it
	imei := ""
	invalidMsgCnt := 0

	for {
		buf := make([]byte, 1024)
		n, err := c.Read(buf)

		if err != nil {
			c.Close()
			return
		}

		req := buf[:n]

		decRes := decoder.Decode(req)
		if decRes.Err != nil {
			// TODO-JAEGER: log the decode error
			invalidMsgCnt++

			if invalidMsgCnt >= MAX_INVALID_MESSAGES_PER_CONN {
				c.Close()
				return
			} else {
				continue
			}
		}

		if decRes.MsgType == "LoginRes" {
			l, err := castDecodeRes[gt06.LoginRes](&decRes)
			if err != nil {
				c.Close()
				return
			}

			// set the imei to associate it with incoming messages
			imei = l.Imei
			c.Write(decRes.Res)

			continue
		}

		// we should not handle any other message type without the imei
		// as we cant know what tracker send the message, therefore the
		// information would be pointless to us
		if imei == "" {
			c.Close()
			// TODO-JAEGER: log that a successfull msg was recieved but
			// imei did not got set
		}

		if decRes.Res != nil {
			c.Write(decRes.Res)
		}

		handleDecodedMsg(&decRes)
	}
}
