package queue

type TrackerEvent struct {
	Imei     string `json:"imei"`
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
	Data     any    `json:"data"`
}
