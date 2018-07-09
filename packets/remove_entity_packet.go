package packets

import (
	"github.com/Irmine/GoMine/net/packets"
	"github.com/infinitygamers/mcpeproxy/info"
)

type RemoveEntityPacket struct {
	*packets.Packet
	RuntimeId uint64
}

func NewRemoveEntityPacket() *RemoveEntityPacket {
	return &RemoveEntityPacket{packets.NewPacket(info.RemoveEntityPacket), 0}
}

func (pk *RemoveEntityPacket) Encode() {
	pk.PutEntityRuntimeId(pk.RuntimeId)
}

func (pk *RemoveEntityPacket) Decode() {
	pk.RuntimeId = pk.GetEntityRuntimeId()
}
