package tcp

import (
	"context"
	"fmt"
	"reciever-ms/tracer"
	"reciever-ms/trackers/gt06"
)

func castDecodeRes[T any](d *gt06.DecodeRes) (*T, error) {
	cast, ok := d.Msg.(T)
	if !ok {
		return nil, fmt.Errorf("failed to cast message of type %s", d.MsgType)
	}

	return &cast, nil
}

// here we should have rastercar bussiness logic, such as publishing recieved positions to rmq, etc...
func handleDecodedMsg(d *gt06.DecodeRes) {
	switch d.MsgType {
	case "LocationRes":
		// TODO:
	}
}

type Session struct {
	// the imei of the tracker in the current connection
	// set once the tracker sends a login packet with it
	Imei          string
	InvalidMsgCnt int
}

func (s *Session) handlePackets(packets []byte) (res []byte, err error) {
	// TODO: fixme (do not use jhoy code)
	_, span := tracer.NewSpan(context.Background(), "tcp", "meh")
	defer span.End()

	decRes := decoder.Decode(packets)
	if decRes.Err != nil {
		// TODO-JAEGER: log the decode error
		s.InvalidMsgCnt++

		if s.InvalidMsgCnt >= MAX_INVALID_MESSAGES_PER_CONN {
			return nil, fmt.Errorf("invalid message count above limit of: %d", MAX_INVALID_MESSAGES_PER_CONN)
		}

		return nil, nil
	}

	if decRes.MsgType == "LoginRes" {
		r, err := castDecodeRes[gt06.LoginRes](&decRes)
		if err != nil {
			return nil, err
		}

		s.Imei = r.Imei

		return decRes.Res, nil
	}

	// we should not handle any other message type without the imei
	// as we cant know what tracker send the message, therefore the
	// information would be pointless to us
	if s.Imei == "" {
		// TODO-JAEGER: log bellow
		return nil, fmt.Errorf("despite ok msg of type: %s the imei has not been set by a login packet previously", decRes.MsgType)
	}

	go handleDecodedMsg(&decRes)

	return decRes.Res, nil
}
