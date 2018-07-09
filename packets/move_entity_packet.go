package packets

import (
	"github.com/golang/geo/r3"
	"github.com/irmine/gomine/net/info"
	"github.com/Irmine/GoMine/net/packets"
	data2 "github.com/irmine/worlds/entities/data"
)

type MoveEntityPacket struct {
	*packets.Packet
	RuntimeId  uint64
	Position   r3.Vector
	Rotation   data2.Rotation
	Mode       byte
	OnGround   bool
	Teleported bool
}

func NewMoveEntityPacket() *MoveEntityPacket {
	return &MoveEntityPacket{Packet: packets.NewPacket(info.PacketIds200[info.MoveEntityPacket]), Position: r3.Vector{}, Rotation: data2.Rotation{}}
}

func (pk *MoveEntityPacket) Encode() {
	pk.PutEntityRuntimeId(pk.RuntimeId)
	pk.PutVector(pk.Position)
	pk.PutEntityRotation(pk.Rotation)
	pk.PutByte(pk.Mode)
	pk.PutBool(pk.OnGround)
	pk.PutBool(pk.Teleported)
}

func (pk *MoveEntityPacket) Decode() {
	pk.RuntimeId = pk.GetEntityRuntimeId()
	pk.Position = pk.GetVector()
	pk.Rotation = pk.GetEntityRotation()
	pk.Mode = pk.GetByte()
	pk.OnGround = pk.GetBool()
	pk.Teleported = pk.GetBool()
}
