package decoder

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reciever-ms/utils/strings"
	"strconv"
	"time"
)

type CourseAndStatus struct {
	GpsRealTime      bool   `json:"gps_real_time"`     // if gps real time is active (docs dont help much)
	GpsPositioned    bool   `json:"gps_positioned"`    // if the gps has been positioned (thats all the docs say...)
	LatHemisphere    string `json:"lat_hemisphere"`    // N (north) or S (south)
	LngHemisphere    string `json:"lng_hemisphere"`    // E (east) or W (west)
	DirectionDegrees int    `json:"direction_degrees"` // direction as clockwise degrees where north is 0 and south is 180
}

type LocationRes struct {
	CourseAndStatus
	Lat               float64   `json:"lat"`                 //
	Lng               float64   `json:"lng"`                 //
	Speed             int       `json:"speed"`               // km/h
	DateTime          time.Time `json:"time"`                // location sending datetime
	SateliteCnt       int       `json:"satelite_cnt"`        // amount of satellites connected to the tracker at that moment
	TowerCellId       string    `json:"tower_cell_id"`       //
	LocationAreaCode  string    `json:"location_area_code"`  // see GSM specification 03.03, 04.08 and 11.11
	MobileCountryCode string    `json:"mobile_country_code"` //
	MobileNetworkCode string    `json:"mobile_network_code"` //
}

// see section 5.2.1.5. of the gt06 protocol docs
func decodeSateliteCnt(b byte) (int, error) {
	hexString := fmt.Sprintf("%02x", b)
	secondHexDigit := hexString[1:]

	i, err := strconv.ParseUint("0"+secondHexDigit, 16, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}

// see section 5.2.1.4. of the gt06 protocol docs
func decodeDateTime(b []byte) (time.Time, error) {
	if len(b) < 6 {
		return time.Time{}, errors.New("cant decode dateTime of len < 6")
	}

	s := fmt.Sprintf("20%02d-%02d-%02dT%02d:%02d:%02d", b[0], b[1], b[2], b[3], b[4], b[5])
	return time.Parse("2006-01-02T15:04:05", s)
}

// decodes a 4byte array to a decimal based coordinate
// see section 5.2.1.6. of the gt06 protocol docs
func decodeCoord(b []byte) float64 {
	i := binary.BigEndian.Uint32(b)

	val := float64(i) / 30000.0

	degrees := int(val / 60.0)
	minutes := math.Mod(val, 60.0)

	mindecimal := (minutes * 100) / 60

	result := float64(degrees) + mindecimal/100

	// truncate to 10 digits precision
	return float64(int(result*10000000000)) / 10000000000
}

// see section 5.2.1.9. of the gt06 protocol docs
func decodeCourseAndStatus(b []byte) (*CourseAndStatus, error) {
	i := binary.BigEndian.Uint16(b)
	bitsStr := fmt.Sprintf("%016b", i)

	lngHemisphere := "S"
	latHemisphere := "E"

	if string(bitsStr[4]) == "1" {
		lngHemisphere = "N"
	}

	if string(bitsStr[5]) == "1" {
		latHemisphere = "W"
	}

	directionDegrees, err := strconv.ParseInt(bitsStr[6:16], 2, 64)
	if err != nil {
		return nil, err
	}

	return &CourseAndStatus{
		GpsRealTime:      string(bitsStr[2]) == "1",
		GpsPositioned:    string(bitsStr[3]) == "1",
		LngHemisphere:    lngHemisphere,
		LatHemisphere:    latHemisphere,
		DirectionDegrees: int(directionDegrees),
	}, nil
}

func (m *Gt06Msg) DecodeLocation() DecodeRes {
	if len(m.InformationContent) != 26 {
		return DecodeRes{Err: errors.New("invalid location information content, size must be 26 bytes")}
	}

	dateTime, err := decodeDateTime(m.InformationContent[0:6])
	if err != nil {
		return DecodeRes{Err: errors.New("failed to parse location date")}
	}

	sateliteCnt, err := decodeSateliteCnt(m.InformationContent[6])
	if err != nil {
		return DecodeRes{Err: errors.New("failed to parse satelite count")}
	}

	courseAndStatus, err := decodeCourseAndStatus(m.InformationContent[16:18])
	if err != nil {
		return DecodeRes{Err: err}
	}

	locationRes := LocationRes{
		Lat:               decodeCoord(m.InformationContent[7:11]),
		Lng:               decodeCoord(m.InformationContent[11:15]),
		Speed:             int(m.InformationContent[15]),
		DateTime:          dateTime,
		SateliteCnt:       sateliteCnt,
		TowerCellId:       strings.BytesToDecimalString(m.InformationContent[23:26]),
		LocationAreaCode:  strings.BytesToDecimalString(m.InformationContent[21:23]),
		MobileCountryCode: strings.BytesToDecimalString(m.InformationContent[18:20]),
		MobileNetworkCode: fmt.Sprintf("%d", m.InformationContent[20]),
		CourseAndStatus:   *courseAndStatus,
	}

	return DecodeRes{Err: nil, Res: nil, Msg: locationRes, MsgType: "LocationRes"}
}
