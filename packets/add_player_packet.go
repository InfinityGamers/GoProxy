package packets

import (
	"github.com/google/uuid"
	"github.com/golang/geo/r3"
	"github.com/irmine/worlds/entities/data"
	"github.com/Irmine/GoMine/net/packets"
	"github.com/infinitygamers/mcpeproxy/info"
)

type AddPlayerPacket struct {
	*packets.Packet
	UUID            uuid.UUID
	Username        string
	DisplayName     string
	Platform        int32
	UnknownString   string
	EntityUniqueId  int64
	EntityRuntimeId uint64
	Position        r3.Vector
	Motion          r3.Vector
	Rotation        data.Rotation
	// HandItem TODO: Items.
	Metadata          map[uint32][]interface{}
	Flags             uint32
	CommandPermission uint32
	Flags2            uint32
	PlayerPermission  uint32
	CustomFlags       uint32
	Long1             int64
}

func NewAddPlayerPacket() *AddPlayerPacket {
	return &AddPlayerPacket{Packet: packets.NewPacket(info.AddPlayerPacket), Metadata: make(map[uint32][]interface{}), Motion: r3.Vector{}}
}

func (pk *AddPlayerPacket) Encode() {
	pk.PutUUID(pk.UUID)
	pk.PutString(pk.Username)
	pk.PutString(pk.DisplayName)
	pk.PutVarInt(pk.Platform)

	pk.PutEntityUniqueId(pk.EntityUniqueId)
	pk.PutEntityRuntimeId(pk.EntityRuntimeId)
	pk.PutString(pk.UnknownString)

	pk.PutVector(pk.Position)
	pk.PutVector(pk.Motion)
	pk.PutPlayerRotation(pk.Rotation)

	pk.PutVarInt(0) // TODO
	pk.PutEntityData(pk.Metadata)

	pk.PutUnsignedVarInt(pk.Flags)
	pk.PutUnsignedVarInt(pk.CommandPermission)
	pk.PutUnsignedVarInt(pk.Flags2)
	pk.PutUnsignedVarInt(pk.PlayerPermission)
	pk.PutUnsignedVarInt(pk.CustomFlags)

	pk.PutVarLong(pk.Long1)

	pk.PutUnsignedVarInt(0) // TODO
}

func (pk *AddPlayerPacket) Decode() {
	pk.UUID = pk.GetUUID()
	pk.Username = pk.GetString()
	pk.DisplayName = pk.GetString()
	pk.Platform = pk.GetVarInt()

	pk.EntityUniqueId = pk.GetEntityUniqueId()
	pk.EntityRuntimeId = pk.GetEntityRuntimeId()
	pk.UnknownString = pk.GetString()

	pk.Position = pk.GetVector()
	pk.Motion = pk.GetVector()
	pk.Rotation = pk.GetPlayerRotation()

	pk.GetVarInt()
	pk.Metadata = pk.GetEntityData()

	pk.Flags = pk.GetUnsignedVarInt()
	pk.CommandPermission = pk.GetUnsignedVarInt()
	pk.Flags2 = pk.GetUnsignedVarInt()
	pk.PlayerPermission = pk.GetUnsignedVarInt()
	pk.CustomFlags = pk.GetUnsignedVarInt()

	pk.Long1 = pk.GetVarLong()

	pk.GetUnsignedVarInt()
}
