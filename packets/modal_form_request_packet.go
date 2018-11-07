package packets

import (
	"github.com/Irmine/GoMine/net/packets"
	"github.com/infinitygamers/goproxy/info"
)

type ModalFormRequestPacket struct {
	*packets.Packet
	FormId uint32
	FormData string
}

func NewModalFormRequestPacket() *ModalFormRequestPacket {
	return &ModalFormRequestPacket{packets.NewPacket(info.ModalFormRequestPacket), 0, ""}
}

func (pk *ModalFormRequestPacket) Encode() {
	pk.PutUnsignedVarInt(pk.FormId)
	pk.PutString(pk.FormData)
}

func (pk *ModalFormRequestPacket) Decode() {
	pk.FormId = pk.GetUnsignedVarInt()
	pk.FormData = pk.GetString()
}