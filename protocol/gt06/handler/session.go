package handler

import (
	"context"
	"errors"
	"fmt"
	"reciever-ms/protocol/gt06/decoder"
	"reciever-ms/tracer"
)

var (
	dec                           = decoder.NewDecoder(true)
	MAX_INVALID_MESSAGES_PER_CONN = 10
)

func castDecodeRes[T any](d *decoder.DecodeRes) (*T, error) {
	cast, ok := d.Msg.(T)
	if !ok {
		return nil, fmt.Errorf("failed to cast message of type %s", d.MsgType)
	}

	return &cast, nil
}

// here we should have rastercar bussiness logic, such as publishing recieved positions to rmq, etc...
func handleDecodedMsg(ctx context.Context, d *decoder.DecodeRes) {
	_, span := tracer.NewSpan(ctx, "fn", "handleDecodedMsg")
	defer span.End()

	tracer.AddSpanTags(span, map[string]string{"msg_type": d.MsgType})

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
	ctx, span := tracer.NewSpan(context.TODO(), "fn", "handlePackets")
	defer span.End()

	if s.Imei != "" {
		tracer.AddSpanTags(span, map[string]string{"imei": s.Imei})
	}

	decRes := dec.Decode(packets)
	if decRes.Err != nil {
		s.InvalidMsgCnt++

		if s.InvalidMsgCnt >= MAX_INVALID_MESSAGES_PER_CONN {
			errMsg := fmt.Sprintf("invalid message count above limit of: %d", MAX_INVALID_MESSAGES_PER_CONN)

			tracer.AddSpanErrorAndFail(span, err, errMsg)
			return nil, errors.New(errMsg)
		}

		tracer.AddSpanError(span, err)
		return nil, nil
	}

	tracer.AddSpanTags(span, map[string]string{"msg_type": decRes.MsgType})

	if decRes.MsgType == "LoginRes" {
		r, err := castDecodeRes[decoder.LoginRes](&decRes)
		if err != nil {
			return nil, err
		}

		s.Imei = r.Imei
		tracer.AddSpanTags(span, map[string]string{"imei": s.Imei})

		return decRes.Res, nil
	}

	// we should not handle any other message type without the imei
	// as we cant know what tracker send the message, therefore the
	// information would be pointless to us
	if s.Imei == "" {
		errMsg := fmt.Sprintf("despite ok msg of type: %s the imei has not been set by a login packet previously", decRes.MsgType)

		tracer.AddSpanErrorAndFail(span, err, errMsg)
		return nil, errors.New(errMsg)
	}

	go handleDecodedMsg(ctx, &decRes)

	return decRes.Res, nil
}
