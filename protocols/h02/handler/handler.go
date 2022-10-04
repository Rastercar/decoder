package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"reciever-ms/config"
	"reciever-ms/protocols/h02/decoder"
	"reciever-ms/tracer"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Handler struct {
	cfg *config.Config
	dec decoder.Decoder
}

func New(cfg *config.Config) Handler {
	return Handler{cfg: cfg, dec: decoder.New(cfg)}
}

// deals with the connection between tracker and decoder, listening to tracker packets
// until the connection is dropped or the too many invalid packets are recieved.
func (h *Handler) HandleRequest(c net.Conn) {
	ctx, span := tracer.NewSpan(context.TODO(), "handler", "HandleRequest")
	span.SetAttributes(attribute.String("protocol", "h02"))

	defer span.End()

	invalidPacketsCnt := 0

	for {
		buf := make([]byte, 1024)

		n, err := c.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				span.SetStatus(codes.Ok, "connection closed")
			} else {
				tracer.AddSpanErrorAndFail(span, err, "connection error")
			}

			c.Close()
			return
		}

		req := buf[:n]
		r, err := h.handlePackets(ctx, req)

		if r != nil {
			if r.Res != nil {
				c.Write(r.Res)
			}

			continue
		}

		if err != nil {
			invalidPacketsCnt++

			if invalidPacketsCnt == h.cfg.App.MaxInvalidPackets {
				span.SetStatus(codes.Error, fmt.Sprintf("max invalid packets (%d) reached", h.cfg.App.MaxInvalidPackets))

				c.Close()
				return
			}
		}
	}
}

func (h *Handler) handlePackets(ctx context.Context, packets []byte) (*decoder.DecodeResult, error) {
	ctx, span := tracer.NewSpan(ctx, "handler", "handlePackets")
	defer span.End()

	res, err := h.dec.Decode(ctx, packets)
	if err != nil {
		tracer.AddSpanErrorAndFail(span, err, "decode failed")
	}

	if res != nil {
		go h.handleDecodedMessage(ctx, res)
	}

	return res, err
}

func (h *Handler) handleDecodedMessage(ctx context.Context, res *decoder.DecodeResult) {
	_, span := tracer.NewSpan(ctx, "handler", "handleDecodedMessage")
	span.AddEvent("some event !")

	defer span.End()
}
