package stream

import (
	"github.com/irmine/binutils"
)

type VirtualItem struct {
	Id int32
	Damage int32
	Count int32
	Nbt string
}

type InventoryActionStream struct {
	*binutils.Stream
}

func NewInventoryActionStream(stream *binutils.Stream) *InventoryActionStream  {
	return &InventoryActionStream{stream}
}

func (pk *InventoryActionStream) PutVirtualItem(item *VirtualItem) {
	pk.PutVarInt(item.Id)
	pk.PutVarInt(((item.Damage & 0x7fff) << 8) | item.Count)
	pk.PutLittleShort(int16(len(item.Nbt)))
	pk.PutString(item.Nbt)
	pk.PutVarInt(0) //todo
	pk.PutVarInt(0) //todo
}

func (pk *InventoryActionStream) GetVirtualItem() *VirtualItem {
	id := pk.GetVarInt()

	if id == 0 {
		return &VirtualItem{Id: 0, Damage: 0, Count: 0, Nbt: ""}
	}

	aux := pk.GetVarInt()
	data := aux >> 8

	if data == 0x7fff {
		data = -1
	}

	count := aux & 0xff

	nbtLen := pk.GetLittleShort()
	nbt := ""
	if nbtLen > 0 {
		//nbt = string(pk.Get(int(nbtLen))) //todo
	}

	return &VirtualItem{Id: id, Damage: data, Count: count, Nbt: nbt}
}