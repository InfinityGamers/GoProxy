package inventoryio

import (
	"github.com/infinitygamers/goproxy/packets/inventoryio/stream"
)

const (
	ContainerSource = iota + 0
	WorldSource = 2
	//CreativeSource = 3
)

type InventoryActionIO struct {
	Source uint32
	WindowId int32
	SourceFlags uint32
	InventorySlot uint32
	OldItem *stream.VirtualItem
	NewItem *stream.VirtualItem
}

func NewInventoryActionIO() InventoryActionIO{
	return InventoryActionIO{}
}

func (IO *InventoryActionIO) WriteToBuffer(bs *stream.InventoryActionStream) {
	bs.PutUnsignedVarInt(IO.Source)

	switch IO.Source {
	case ContainerSource:
		bs.PutVarInt(IO.WindowId)
		break
	case WorldSource:
		bs.PutUnsignedVarInt(IO.SourceFlags)
		break
	}

	bs.PutUnsignedVarInt(IO.InventorySlot)
	bs.PutVirtualItem(IO.OldItem)
	bs.PutVirtualItem(IO.NewItem)
}

func (IO *InventoryActionIO) ReadFromBuffer(bs *stream.InventoryActionStream) InventoryActionIO {
	v := NewInventoryActionIO()
	v.Source = bs.GetUnsignedVarInt()
	switch v.Source {
	case ContainerSource:
		v.WindowId = bs.GetVarInt()
		break
	case WorldSource:
		v.SourceFlags = bs.GetUnsignedVarInt()
		break
	}

	IO.InventorySlot = bs.GetUnsignedVarInt()
	IO.OldItem = bs.GetVirtualItem()
	IO.NewItem = bs.GetVirtualItem()

	return v
}