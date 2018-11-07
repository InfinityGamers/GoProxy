package packets

import (
	"github.com/Irmine/GoMine/net/packets"
	"github.com/infinitygamers/goproxy/info"
)

type ModalFormResponsePacket struct {
	*packets.Packet
	FormId uint32
	FormData string
}

func NewModalFormResponsePacket() *ModalFormResponsePacket {
	return &ModalFormResponsePacket{packets.NewPacket(info.ModalFormResponsePacket), 0, ""}
}

func (pk *ModalFormResponsePacket) Encode() {
	pk.PutUnsignedVarInt(pk.FormId)
	pk.PutString(pk.FormData)
}

func (pk *ModalFormResponsePacket) Decode() {
	pk.FormId = pk.GetUnsignedVarInt()
	pk.FormData = pk.GetString()
}