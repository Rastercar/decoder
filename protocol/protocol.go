package protocol

type DecodeResult struct {
	Res     []byte      // Response to send to the tracker
	Msg     interface{} // The decoded message
	MsgType string      // Name of the struct with the decoded message based on the msg protocol num
}
