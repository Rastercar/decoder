package protocol

import "reciever-ms/queue"

type DecodeResult struct {
	Res []byte              // Response to send to the tracker
	Evt *queue.TrackerEvent // Event to publish on the tracker events queue
}
