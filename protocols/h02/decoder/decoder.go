package decoder

import (
	"errors"
	"fmt"
	"log"
	dstrings "reciever-ms/utils/strings"
	"strings"
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

type DecodeResult struct {
	Res     []byte      // Response to send to the tracker
	Msg     interface{} // The decoded message
	MsgType string      // Name of the struct with the decoded message based on the msg protocol num
}

type decoder struct {
	debug bool
}

func (d *decoder) logIfDebug(s string, a ...any) {
	if d.debug {
		log.Printf(s, a...)
	}
}

func New(debug bool) decoder {
	return decoder{debug}
}

func (d *decoder) Decode(b []byte) (*DecodeResult, error) {
	s := string(b)
	d.logIfDebug(s)

	s = dstrings.GetStringInBetween(s, "*HQ,", "#")
	if s == "" {
		return nil, errors.New("preffix (*HQ) and suffix (#) not found")
	}

	parts := strings.Split(s, ",")
	if len(parts) < 2 {
		return nil, errors.New("msg size is too low to get operation/command name")
	}

	cmd := parts[1]

	switch cmd {
	case MSG_HEARTBEAT:
	case MSG_HEARTBEAT_V2:
		return d.decodeHeartbeat(parts)

	case MSG_LOCATION:
		return d.decodeLocation(parts)
	}

	return nil, fmt.Errorf("unknown command/msg type: %s", cmd)
}

func (d *decoder) decodeHeartbeat(parts []string) (*DecodeResult, error) {
	if len(parts) == 0 {
		return nil, errors.New("cant decode heartbeat packet with no parts")
	}

	return &DecodeResult{
		Res:     nil,
		Msg:     HeartbeatMsg{Imei: parts[0]},
		MsgType: "HeartbeatMsg",
	}, nil
}

func (d *decoder) decodeLocation(parts []string) (*DecodeResult, error) {
	packet, err := PacketFromParts(parts)
	if err != nil {
		return nil, err
	}

	return packet.ParseToLocationMsg()
}
