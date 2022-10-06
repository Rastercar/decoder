# H02 events suport

*Uplink*
- ✅ heartbeat
- ✅ real time location
- ❌ location request
- ❌ blind spots uploading
- ❌ device alarm 

*Downlink*
- ❌ cut-off/recover oil and engine 
- ❌ response to location request
- ❌ fortification (SF)
- ❌ fortification (SF2)
- ❌ disarming (CF)
- ❌ disarming (CF2)
- ❌ main number bind (UR)
- ❌ server setting (IP)
- ❌ terminal password setting
- ❌ interval settings
- ❌ allarm setting
- ❌ device reboot
- ❌ reset to defaults
- ❌ network access point
- ❌ answer mode
- ❌ IMEI setting
- ❌ language setting
- ❌ audiomonitor
- ❌ query device information
- ❌ working mode setting

---

## Events example

### location

```JSON
{
  "imei": "867232051148352",
  "type": "location",
  "protocol": "h02",
  "data": {
    "lat": -20.465548333333334,
    "lng": -54.582398166666664,
    "imei": "867232051148352",
    "speed": 0,
    "direction": 0,
    "date": "2022-07-11T04:46:39Z",
    "status": {
      "temperature_alarm": false,
      "three_times_pass_err_alarm": false,
      "gprs_occlusion_alarm": false,
      "oil_and_engine_cut_off": false,
      "storage_battery_removal_state": false,
      "high_level_sensor_1": false,
      "high_level_sensor_2": false,
      "low_level_sensor_1_bond_strap": false,
      "gps_reciever_fault_alarm": false,
      "analog_quantity_transfinite_alarm": false,
      "sos_alarm": false,
      "host_powered_by_backup_battery": false,
      "storage_battery_removed": false,
      "open_circuit_for_gps_antenna": false,
      "shor_circuit_for_gps_antenna": false,
      "low_level_sensor_2_bond_strap": false,
      "door_open": false,
      "vehicle_fortified": false,
      "acc": false,
      "engine": true,
      "custom_alarm": false,
      "overspeed": false,
      "theft_alarm": false,
      "roberry_alarm": false,
      "overspeed_alarm": false,
      "ilegal_ignition_alarm": false,
      "no_entry_cross_border_alarm_in": false,
      "gps_antenna_open_circuit_alarm": false,
      "gps_antenna_short_circuit_alarm": false,
      "no_entry_cross_border_alarm_out": false
    }
  }
}
```