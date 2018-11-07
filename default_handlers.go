package goproxy

import (
	"github.com/infinitygamers/goproxy/info"
	"github.com/infinitygamers/goproxy/packets"
)

func RegisterDefaultHandlers(pkHandler *PacketHandler)  {
	pkHandler.RegisterHandler(info.MovePlayerPacket, func(bytes []byte, host Host, connection *Connection) bool {
		move := packets.NewMovePlayerPacket()
		move.Buffer = bytes
		move.DecodeHeader()
		move.Decode()

		if host.IsServer() {
			connection.Client.Position = move.Position
		}

		return false
	})
}