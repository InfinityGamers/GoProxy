package packets

import (
	"github.com/Irmine/GoMine/net/packets"
	"github.com/infinitygamers/goproxy/info"
)

const (
	ClearTitle = iota
	ResetTitle
	SetTitle
	SetSubtitle
	SetActionBarMessage
	SetAnimationTimes
)

type SetTitlePacket struct {
	*packets.Packet
	TitleType int32
	Text string
	FadeInTime int32
	StayTime int32
	FadeOutTime int32
}

func NewSetTitlePacket() *SetTitlePacket {
	return &SetTitlePacket{
		packets.NewPacket(info.SetTitlePacket),
		0,
		"",
		0,
		0,
		0,
	}
}

func (pk *SetTitlePacket) Encode() {
	pk.PutVarInt(pk.TitleType)
	pk.PutString(pk.Text)
	pk.PutVarInt(pk.FadeInTime)
	pk.PutVarInt(pk.StayTime)
	pk.PutVarInt(pk.FadeOutTime)
}

func (pk *SetTitlePacket) Decode() {
	pk.TitleType = pk.GetVarInt()
	pk.Text = pk.GetString()
	pk.FadeInTime = pk.GetVarInt()
	pk.StayTime = pk.GetVarInt()
	pk.FadeOutTime = pk.GetVarInt()
}