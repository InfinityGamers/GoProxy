package packets

import (
	"github.com/golang/geo/r3"
	"github.com/irmine/gomine/net/info"
	"github.com/Irmine/GoMine/net/packets"
	"github.com/infinitygamers/mcpeproxy/packets/inventoryio"
	"github.com/infinitygamers/mcpeproxy/packets/inventoryio/stream"
)

/**
 * Transaction Types
 */
const (
	Normal = iota + 0
	Mismatch
	UseItem
	UseItemOnEntity
	ReleaseItem
)

// Action Types
const (
	ItemClickBlock = iota + 0
	ItemClickAir
	ItemBreakBlock

	//CONSUMABLE ITEMS
	ItemRelease = iota + 0
	ItemConsume
)

// Entity Action ztypes
const (
	ItemOnEntityInteract = iota + 0
	ItemOnEntityAttack
)

type InventoryTransactionPacket struct {
	*packets.Packet
	InvStream *stream.InventoryActionStream
	ActionList *inventoryio.InventoryActionIOList
	TransactionType, ActionType uint32
	Face, HotbarSlot int32
	ItemSlot *stream.VirtualItem
	BlockX int32
	BlockY uint32
	BlockZ int32
	PlayerPosition, ClickPosition, HeadPosition r3.Vector
	RuntimeId uint64
}

func NewInventoryTransactionPacket() *InventoryTransactionPacket {
	pk := &InventoryTransactionPacket{Packet: packets.NewPacket(info.PacketIds200[info.InventoryTransactionPacket]),
		ActionList: inventoryio.NewInventoryActionIOList(),
		TransactionType: 0,
		ActionType: 0,
		Face: 0,
		HotbarSlot: 0,
		ItemSlot: &stream.VirtualItem{},
		PlayerPosition: r3.Vector{},
		ClickPosition: r3.Vector{},
		HeadPosition: r3.Vector{},
		RuntimeId: 0,
	}
	pk.InvStream = stream.NewInventoryActionStream(pk.Stream)
	return pk
}


func (pk *InventoryTransactionPacket) GetBlockPos() {
	pk.BlockX = pk.GetVarInt()
	pk.BlockY = pk.GetUnsignedVarInt()
	pk.BlockZ = pk.GetVarInt()
}

func (pk *InventoryTransactionPacket) PutBlockPos() {
	pk.PutVarInt(pk.BlockX)
	pk.PutUnsignedVarInt(pk.BlockY)
	pk.PutVarInt(pk.BlockZ)
}

func (pk *InventoryTransactionPacket) Encode()  {
	pk.PutUnsignedVarInt(pk.TransactionType)
	pk.ActionList.WriteToBuffer(pk.InvStream)

	switch pk.TransactionType {
	case Normal, Mismatch:
		break
	case UseItem:
		pk.PutUnsignedVarInt(pk.ActionType)
		pk.PutBlockPos()
		pk.PutVarInt(pk.Face)
		pk.PutVarInt(pk.HotbarSlot)
		pk.InvStream.PutVirtualItem(pk.ItemSlot)
		pk.PutVector(pk.PlayerPosition)
		pk.PutVector(pk.ClickPosition)
		break
	case UseItemOnEntity:
		pk.PutUnsignedVarLong(pk.RuntimeId)
		pk.PutUnsignedVarInt(pk.ActionType)
		pk.PutVarInt(pk.HotbarSlot)
		pk.InvStream.PutVirtualItem(pk.ItemSlot)
		pk.PutVector(pk.PlayerPosition)
		pk.PutVector(pk.ClickPosition)
		break
	case ReleaseItem:
		pk.PutUnsignedVarInt(pk.ActionType)
		pk.PutVarInt(pk.HotbarSlot)
		pk.InvStream.PutVirtualItem(pk.ItemSlot)
		pk.PutVector(pk.HeadPosition)
		break
	default:
		panic("Unknown transaction type passed: " + string(pk.TransactionType))
	}
}

func (pk *InventoryTransactionPacket) Decode() {
	pk.TransactionType = pk.GetUnsignedVarInt()

	pk.ActionList.ReadFromBuffer(pk.InvStream)

	switch pk.TransactionType{
	case Normal, Mismatch:
		break
	case UseItem:
		pk.ActionType = pk.GetUnsignedVarInt()
		pk.GetBlockPos()
		pk.Face = pk.GetVarInt()
		pk.HotbarSlot = pk.GetVarInt()
		pk.ItemSlot = pk.InvStream.GetVirtualItem()
		pk.PlayerPosition = pk.GetVector()
		pk.ClickPosition = pk.GetVector()
	case UseItemOnEntity:
		pk.RuntimeId = pk.GetUnsignedVarLong()
		pk.ActionType = pk.GetUnsignedVarInt()
		pk.HotbarSlot = pk.GetVarInt()
		pk.ItemSlot = pk.InvStream.GetVirtualItem()
		pk.PlayerPosition = pk.GetVector()
		pk.ClickPosition = pk.GetVector()
	case ReleaseItem:
		pk.ActionType = pk.GetUnsignedVarInt()
		pk.HotbarSlot = pk.GetVarInt()
		pk.ItemSlot = pk.InvStream.GetVirtualItem()
		pk.HeadPosition = pk.GetVector()
	default:
		panic("Error: Unknown transaction type received: " + string(pk.TransactionType))
	}
}
