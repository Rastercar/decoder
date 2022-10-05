package decoder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reciever-ms/config"
	"reciever-ms/protocol"
	"reciever-ms/tracer"
	dstrings "reciever-ms/utils/strings"
	"strings"

	"go.opentelemetry.io/otel/attribute"
)

var (
	MSG_LOCATION                      = "V1"
	MSG_HEARTBEAT                     = "XT"
	MSG_HEARTBEAT_V2                  = "HTBT"
	MSG_LOCATION_REQ                  = "VI1"
	MSG_LOCATION_RES                  = "VI"
	MSG_CUT_OFF_OIL_ENGINE            = "S20"
	MSG_INTRUCTION_ACK                = "V4"
	MSG_FORTIFICATION                 = "SF"
	MSG_FORTIFICATION_V2              = "SF2"
	MSG_DISARM                        = "CF"
	MSG_DISARM_V2                     = "CF2"
	MSG_PLATAFORM_DISTRIBUTED_SMS     = "TG"
	MSG_MAIN_NUMBER_BIND              = "UR"
	MSG_MODIFY_IP                     = "IP"
	MSG_SETTING_SMS_INTERCEPTOR_NUM   = "ST"
	MSG_TERMINAL_PASSWORD_SETTING     = "MP"
	MSG_UPLOADING_INTERVAL_SETTING    = "XT"
	MSG_UPLOADING_INTERVAL_SETTING_V2 = "NXT"
	MSG_ALARM_SETTING                 = "KC"
	MSG_DEVICE_REBOOT                 = "CQ"
	MSG_RESET_TO_DEFALTS              = "RESET"
	MSG_APN_SETTING                   = "APN"
	MSG_FAMILY_NUMBER_SETTING         = "SQQ"
	MSG_ANSWER_MODE_SETTING           = "ACPC"
	MSG_IMEI_SETTING                  = "SIMEI"
	MSG_LANGUAGE_SETTING              = "SLAN"
	MSG_MONITOR                       = "CALB"
	MSG_POWER_SAVING_MODE_SETTING     = "PWM"
	MSG_OVERSPEED_SETTING             = "OVSP"
	MSG_QUERY_DEVICE_STATUS           = "INFO"
	MSG_ALARM                         = "ALRM"
)

type Decoder struct {
	cfg *config.Config
}

func (d *Decoder) logIfDebug(s string, a ...any) {
	if d.cfg.App.Debug {
		log.Printf(s, a...)
	}
}

func New(cfg *config.Config) Decoder {
	return Decoder{cfg}
}

func (d *Decoder) Decode(ctx context.Context, b []byte) (*protocol.DecodeResult, error) {
	_, span := tracer.NewSpan(ctx, "decoder", "Decode")
	defer span.End()

	s := string(b)
	d.logIfDebug(s)

	span.SetAttributes(attribute.Key("msg_str").String(s))
	span.SetAttributes(attribute.Key("msg_hex").String(fmt.Sprintf("%x", b)))

	s = dstrings.GetStringInBetween(s, "*HQ,", "#")
	if s == "" {
		return nil, errors.New("preffix (*HQ) and suffix (#) not found")
	}

	parts := strings.Split(s, ",")
	if len(parts) < 2 {
		return nil, errors.New("msg size is too low to get operation/command name")
	}

	cmd := parts[1]

	span.SetAttributes(attribute.Key("msg_type").String(cmd))

	switch cmd {
	case MSG_HEARTBEAT:
	case MSG_HEARTBEAT_V2:
		return d.decodeHeartbeat(parts)

	case MSG_LOCATION:
		return d.decodeLocation(parts)
	}

	return nil, fmt.Errorf("unknown command/msg type: %s", cmd)
}

func (d *Decoder) decodeHeartbeat(parts []string) (*protocol.DecodeResult, error) {
	if len(parts) == 0 {
		return nil, errors.New("cant decode heartbeat packet with no parts")
	}

	return &protocol.DecodeResult{
		Res:     nil,
		Msg:     HeartbeatMsg{Imei: parts[0]},
		MsgType: "HeartbeatMsg",
	}, nil
}

func (d *Decoder) decodeLocation(parts []string) (*protocol.DecodeResult, error) {
	packet, err := PacketFromParts(parts)
	if err != nil {
		return nil, err
	}

	return packet.ParseToLocationMsg()
}
