package decoder

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reciever-ms/protocol"
	"reciever-ms/queue"
	"strconv"
	"time"
)

type LocationMsg struct {
	Lat       float64      `json:"lat"`       // latitude (90 to -90) in decimal degrees
	Lng       float64      `json:"lng"`       // longitude (180 to -180) in decimal degrees
	Imei      string       `json:"imei"`      // tracker imei
	Speed     int          `json:"speed"`     // speed in km/h
	Status    StatusPacket `json:"status"`    // info about vehicle / tracker status
	Direction int          `json:"direction"` // direction in degrees (0 degrees = North, 180 = S)
	Timestamp time.Time    `json:"date"`      // date and time sent by the tracker
}

// Latitude，format DDFF.FFFF
// DD      Latitude Degree（00 ~ 90）
// FF.FFFF Latitude points（00.0000 ~ 59.9999) reserved four decimals
//
// Longitude，format DDDFF.FFFF
// DDD     Longitude Degree（000 ~ 180)
// FF.FFFF Longitude points（00.0000 ~ 59.9999) reserved four decimals
func parseLatLng(s string, degreeDigits int) (float64, error) {
	if len(s) < 9 {
		return 0, errors.New("lat/lng does not have 8 digits")
	}

	degrees, err := strconv.Atoi(s[:degreeDigits])
	if err != nil {
		return 0, fmt.Errorf("failed to converse lat/lng degrees: %s", s[:degreeDigits])
	}

	minutes, err := strconv.ParseFloat(s[degreeDigits:], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to converse lat/lng points: %s", s[degreeDigits:])
	}

	return float64(degrees) + minutes/60, nil
}

func parseDate(d string) (time.Time, error) {
	if len(d) < 12 {
		return time.Time{}, errors.New("failed to parse date as it does not follow the ddmmyyhhmmss format")
	}

	layout := "2006-01-02T15:04:05"
	dateIso := "20" + d[4:6] + "-" + d[2:4] + "-" + d[0:2] + "T" + d[6:8] + ":" + d[8:10] + ":" + d[10:12]

	return time.Parse(layout, dateIso)
}

func (m *Packet) ParseToLocationMsg() (*protocol.DecodeResult, error) {
	if m.DataValidBit != "A" {
		return nil, errors.New("invalid location data (data valid bit != A)")
	}

	direction, err := strconv.Atoi(m.DirectionDegrees)
	if err != nil {
		return nil, fmt.Errorf("cant convert direction (%s) to int", m.DirectionDegrees)
	}

	speedF, err := strconv.ParseFloat(m.Speed, 64)
	if err != nil {
		return nil, fmt.Errorf("cant convert speed (%s) to float64", m.Speed)
	}

	// convert knots/h to km/h
	speed := int(speedF * 1.852)

	lat, err := parseLatLng(m.Lat, 2)
	if err != nil {
		return nil, err
	}

	if m.LatSymbol == "S" || m.LatSymbol == "s" {
		lat = lat * -1
	}

	lng, err := parseLatLng(m.Lng, 3)
	if err != nil {
		return nil, err
	}

	if m.LngSymbol == "W" || m.LngSymbol == "w" {
		lng = lng * -1
	}

	timestamp, err := parseDate(m.DayMonthYear + m.Time)
	if err != nil {
		return nil, err
	}

	// https://github.com/traccar/traccar/blob/master/src/main/java/org/traccar/protocol/H02ProtocolDecoder.java
	statusBytes, err := hex.DecodeString(m.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to decode status - %v", err)
	}

	status, err := decodeStatus(statusBytes)
	if err != nil {
		return nil, err
	}

	return &protocol.DecodeResult{
		Res: nil,
		Evt: &queue.TrackerEvent{
			Type: "h02:LocationMsg",
			Data: LocationMsg{
				Lat:       lat,
				Lng:       lng,
				Imei:      m.Imei,
				Speed:     speed,
				Direction: direction,
				Timestamp: timestamp,
				Status:    *status,
			},
			Imei: m.Imei,
		},
	}, nil
}
