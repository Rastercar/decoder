package queue

type TrackerEvent struct {
	Imei string `json:"imei"`
	Type string `json:"type"` // <protocol>:<event_name> (eg: h02:LocationMsg)
	Data any    `json:"data"`
}
