package mcpeproxy

import (
	"net"
	"github.com/Irmine/GoMine/net/packets"
)

type Host interface {
	WritePacket(buffer []byte)
	SendPacket(packets.IPacket)
	SendBatchPacket([]packets.IPacket)
	GetAddress() net.UDPAddr
	IsClient() bool
	IsServer() bool
}