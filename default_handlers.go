package mcpeproxy

import (
	"github.com/infinitygamers/mcpeproxy/info"
	"github.com/infinitygamers/mcpeproxy/packets"
)

func RegisterDefaultHandlers(pkHandler *PacketHandler)  {
	pkHandler.RegisterHandler(info.MovePlayerPacket, func(bytes []byte, host Host, connection *Connection) bool {
		move := packets.NewMovePlayerPacket()
		move.Buffer = bytes
		move.DecodeHeader()
		move.Decode()

		connection.client.Position = move.Position

		return false
	})
}