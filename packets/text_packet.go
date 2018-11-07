package packets

import (
	"github.com/Irmine/GoMine/net/packets"
	"github.com/infinitygamers/goproxy/info"
)

const (
	Raw          = 0x00
	Chat         = 0x01
	Translation  = 0x02
	Popup        = 0x03
	JukeboxPopup = 0x04
	Tip          = 0x05
	System       = 0x06
	Whisper      = 0x07
	Announcement = 0x08
)

type TextPacket struct {
	*packets.Packet
	TextType byte
	Translation bool
	Source string
	SourceThirdParty string
	SourcePlatform int32
	Message string
	Xuid string
	PlatformChatId string

	Params []string
}

func NewTextPacket() *TextPacket {
	return &TextPacket{
		packets.NewPacket(info.TextPacket),
		Raw,
		false,
		"",
		"",
		0,
		"",
		"",
		"",
		[]string{},
	}
}

func (pk *TextPacket) Encode()  {
	pk.PutByte(pk.TextType)
	pk.PutBool(pk.Translation)

	switch pk.TextType {
	case Raw, Tip, System:
		pk.PutString(pk.Message)
		break
	case Chat, Whisper, Announcement:
		pk.PutString(pk.Source)
		pk.PutString(pk.SourceThirdParty)
		pk.PutVarInt(pk.SourcePlatform)
		pk.PutString(pk.Message)
		break
	case Translation, Popup, JukeboxPopup:
		pk.PutString(pk.Message)
		pk.PutUnsignedVarInt(uint32(len(pk.Params)))
		for _, v := range pk.Params {
			pk.PutString(v)
		}
		break
	}

	pk.PutString(pk.Xuid)
	pk.PutString(pk.PlatformChatId)
}

func (pk *TextPacket) Decode() {
	pk.TextType = pk.GetByte()
	pk.Translation = pk.GetBool()

	switch pk.TextType {
	case Raw, Tip, System:
		pk.Message = pk.GetString()
		break
	case Chat, Whisper, Announcement:
		pk.Source = pk.GetString()
		pk.SourceThirdParty = pk.GetString()
		pk.SourcePlatform = pk.GetVarInt()
		pk.Message = pk.GetString()
		break
	case Translation, Popup, JukeboxPopup:
		pk.Message = pk.GetString()
		c := pk.GetUnsignedVarInt()
		for i := uint32(0); i < c; i++ {
			pk.Params = append(pk.Params, pk.GetString())
		}
		break
	}

	pk.Xuid = pk.GetString()
	pk.PlatformChatId = pk.GetString()
}