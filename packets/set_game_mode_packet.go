package packets

import (
	"github.com/Irmine/GoMine/net/packets"
	"github.com/infinitygamers/mcpeproxy/info"
)

type SetGameModePacket struct {
	*packets.Packet
	GameMode int32
}

func NewSetGamemodePacket() *SetGameModePacket {
	return &SetGameModePacket{ packets.NewPacket(info.SetPlayerGameTypePacket), 0}
}

func (pk *SetGameModePacket) Encode() {
	pk.PutVarInt(pk.GameMode)
}

func (pk *SetGameModePacket) Decode() {
	pk.GameMode = pk.GetVarInt()
}