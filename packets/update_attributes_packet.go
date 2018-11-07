package packets

import (
	"github.com/Irmine/GoMine/net/packets"
	"github.com/irmine/worlds/entities/data"
	info2 "github.com/infinitygamers/goproxy/info"
)

type UpdateAttributesPacket struct {
	*packets.Packet
	RuntimeId  uint64
	Attributes data.AttributeMap
}

func NewUpdateAttributesPacket() *UpdateAttributesPacket {
	return &UpdateAttributesPacket{packets.NewPacket(info2.UpdateAttributesPacket), 0, data.NewAttributeMap()}
}

func (pk *UpdateAttributesPacket) Encode() {
	pk.PutEntityRuntimeId(pk.RuntimeId)
	pk.PutAttributeMap(pk.Attributes)
}

func (pk *UpdateAttributesPacket) Decode() {
	pk.RuntimeId = pk.GetEntityRuntimeId()
	pk.Attributes = pk.GetAttributeMap()
}
