package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"reciever-ms/config"
	"reciever-ms/protocol"
	"reciever-ms/protocol/h02/decoder"
	"reciever-ms/queue"
	"reciever-ms/tracer"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type Handler struct {
	cfg *config.Config
	dec decoder.Decoder
	rmq *queue.Server
}

func New(cfg *config.Config, rmq *queue.Server) Handler {
	return Handler{
		cfg: cfg,
		rmq: rmq,
		dec: decoder.New(cfg),
	}
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

		r, err := h.handlePackets(ctx, buf[:n])
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

func (h *Handler) handlePackets(ctx context.Context, packets []byte) (*protocol.DecodeResult, error) {
	ctx, span := tracer.NewSpan(ctx, "handler", "handlePackets")
	defer span.End()

	res, err := h.dec.Decode(ctx, packets)
	if err != nil {
		tracer.AddSpanErrorAndFail(span, err, "decode failed")
	}

	if res != nil && res.Evt != nil {
		go h.sendTrackerEvent(ctx, *res.Evt)
	}

	return res, err
}

func (h *Handler) sendTrackerEvent(ctx context.Context, evt queue.TrackerEvent) {
	ctx, span := tracer.NewSpan(ctx, "handler", "handleDecodedMessage")
	defer span.End()

	switch evt.Type {
	case "location":
		h.rmq.PublishTrackerEvent(ctx, evt)

	case "heartbeat":
		span.AddEvent(fmt.Sprintf("skipping publishing for event of type: %s", evt.Type))
		return

	default:
		span.SetStatus(codes.Error, fmt.Sprintf("unknown event type: %s", evt.Type))
	}
}
