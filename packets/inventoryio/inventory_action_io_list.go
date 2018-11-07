package inventoryio

import (
	"github.com/infinitygamers/goproxy/packets/inventoryio/stream"
)

type InventoryActionIOList struct {
	List []InventoryActionIO
}

func NewInventoryActionIOList() *InventoryActionIOList{
	return &InventoryActionIOList{}
}

func (IOList *InventoryActionIOList) GetCount() int {
	return len(IOList.List)
}

func (IOList *InventoryActionIOList) PutAction(io InventoryActionIO) {
	IOList.List = append(IOList.List, io)
}

func (IOList *InventoryActionIOList) WriteToBuffer(bs *stream.InventoryActionStream) {
	c := len(IOList.List)
	bs.PutUnsignedVarInt(uint32(c))
	for i := 0; i < c; i++ {
		IOList.List[i].WriteToBuffer(bs)
	}
}

func (IOList *InventoryActionIOList) ReadFromBuffer(bs *stream.InventoryActionStream) *InventoryActionIOList{
	c := bs.GetUnsignedVarInt()
	for i := uint32(0); i < c; i ++{
		a := NewInventoryActionIO()
		a.ReadFromBuffer(bs)
		IOList.PutAction(a)
	}
	return IOList
}