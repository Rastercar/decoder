package decoder

import (
	"encoding/binary"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

type StatusPacket struct {
	TemperatureAlarm              bool `json:"temperature_alarm"`
	ThreeTimesPassErrorAlarm      bool `json:"three_times_pass_err_alarm"`
	GprsOcclusionAlarm            bool `json:"gprs_occlusion_alarm"`
	OilAndEngineCutOff            bool `json:"oil_and_engine_cut_off"`
	StorageBatteryRemovalState    bool `json:"storage_battery_removal_state"`
	HighLevelSensor1              bool `json:"high_level_sensor_1"`
	HighLevelSensor2              bool `json:"high_level_sensor_2"`
	LowLevelSensor1BondStrap      bool `json:"low_level_sensor_1_bond_strap"`
	GpsRecieverFaultAlarm         bool `json:"gps_reciever_fault_alarm"`
	AnalogQuantityTransfinitAlarm bool `json:"analog_quantity_transfinite_alarm"`
	SosAlarm                      bool `json:"sos_alarm"`
	HostPoweredByBackupBattery    bool `json:"host_powered_by_backup_battery"`
	StorageBatteryRemoved         bool `json:"storage_battery_removed"`
	OpenCircuitForGpsAntenna      bool `json:"open_circuit_for_gps_antenna"`
	ShortCircuitForGpsAntenna     bool `json:"shor_circuit_for_gps_antenna"`
	LowLevelSensor2BondStrap      bool `json:"low_level_sensor_2_bond_strap"`
	DoorOpen                      bool `json:"door_open"`
	VehicleFortified              bool `json:"vehicle_fortified"`
	Acc                           bool `json:"acc"`
	Engine                        bool `json:"engine"`
	CustomAlarm                   bool `json:"custom_alarm"`
	Overspeed                     bool `json:"overspeed"`
	TheftAlarm                    bool `json:"theft_alarm"`
	RoberryAlarm                  bool `json:"roberry_alarm"`
	OverspeedAlarm                bool `json:"overspeed_alarm"`
	IlegalIgnitionAlarm           bool `json:"ilegal_ignition_alarm"`
	NoEntryCrossBorderAlarmIn     bool `json:"no_entry_cross_border_alarm_in"`
	GpsAntennaOpenCircuitAlarm    bool `json:"gps_antenna_open_circuit_alarm"`
	GpsAntennaShortCircuitAlarm   bool `json:"gps_antenna_short_circuit_alarm"`
	NoEntryCrossBorderAlarmOut    bool `json:"no_entry_cross_border_alarm_out"`
}

// Refer to the "3. Packet definition" of H02_protocol.pdf
type Packet struct {
	Imei             string `json:"imei"`
	Cmd              string `json:"cmd"`
	Time             string `json:"time"`
	DataValidBit     string `json:"data_valid_bit"`
	Lat              string `json:"lat"`
	LatSymbol        string `json:"lat_symbol"`
	Lng              string `json:"lng"`
	LngSymbol        string `json:"lng_symbol"`
	Speed            string `json:"speed"`
	DirectionDegrees string `json:"direction_degrees"`
	DayMonthYear     string `json:"day_month_year"`
	Status           string `json:"status"`
}

func decodeStatus(b []byte) (*StatusPacket, error) {
	spew.Dump(b)

	if len(b) < 4 {
		return nil, fmt.Errorf("cannot decode status of %d bytes", len(b))
	}

	i := binary.BigEndian.Uint32(b)
	binStr := fmt.Sprintf("%32b", i)

	println(binStr)

	var binBool = func(idx int) bool {
		return string(binStr[idx]) == "0"
	}

	// See H02_protocol.pdf "4. Terminal Status (alarm) analysis"
	return &StatusPacket{
		// byte 1
		TemperatureAlarm:           binBool(0),
		ThreeTimesPassErrorAlarm:   binBool(1),
		GprsOcclusionAlarm:         binBool(2),
		OilAndEngineCutOff:         binBool(3),
		StorageBatteryRemovalState: binBool(4),
		HighLevelSensor1:           binBool(5),
		HighLevelSensor2:           binBool(6),
		LowLevelSensor1BondStrap:   binBool(7),

		// byte 2
		GpsRecieverFaultAlarm:         binBool(8),
		AnalogQuantityTransfinitAlarm: binBool(9),
		SosAlarm:                      binBool(10),
		HostPoweredByBackupBattery:    binBool(11),
		StorageBatteryRemoved:         binBool(12),
		OpenCircuitForGpsAntenna:      binBool(13),
		ShortCircuitForGpsAntenna:     binBool(14),
		LowLevelSensor2BondStrap:      binBool(15),

		// byte 3
		DoorOpen:         binBool(16),
		VehicleFortified: binBool(17),
		Acc:              binBool(18),
		// reserved
		// reserved
		Engine:      binBool(21),
		CustomAlarm: binBool(22),
		Overspeed:   binBool(23),

		// byte 4
		TheftAlarm:                  binBool(24),
		RoberryAlarm:                binBool(25),
		OverspeedAlarm:              binBool(26),
		IlegalIgnitionAlarm:         binBool(27),
		NoEntryCrossBorderAlarmIn:   binBool(28),
		GpsAntennaOpenCircuitAlarm:  binBool(29),
		GpsAntennaShortCircuitAlarm: binBool(30),
		NoEntryCrossBorderAlarmOut:  binBool(31),
	}, nil
}

// Create a packet string array containing the parts of the package
// definition (see: "3. Packet definition" of H02_protocol.pdf)
func PacketFromParts(parts []string) (*Packet, error) {
	if len(parts) < 12 {
		return nil, fmt.Errorf("cant create packets with less than 12 parts (got: %d)", len(parts))
	}

	return &Packet{
		Imei:             parts[0],
		Cmd:              parts[1],
		Time:             parts[2],
		DataValidBit:     parts[3],
		Lat:              parts[4],
		LatSymbol:        parts[5],
		Lng:              parts[6],
		LngSymbol:        parts[7],
		Speed:            parts[8],
		DirectionDegrees: parts[9],
		DayMonthYear:     parts[10],
		Status:           parts[11],
	}, nil
}
