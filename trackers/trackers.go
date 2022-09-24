package trackers

type PacketDecoder interface {
	Decode(packets []byte)
}
