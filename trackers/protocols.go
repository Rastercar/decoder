package protocols

type PacketDecoder interface {
	Decode(packets []byte)
}
