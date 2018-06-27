package mcpeproxy

import (
	"github.com/infinitygamers/mcpeproxy/info"
	"github.com/infinitygamers/mcpeproxy/packets"
)

var LastY float64 = 0

func RegisterDefaultHandlers(pkHandler *PacketHandler)  {
	pkHandler.RegisterHandler(info.MovePlayerPacket, func(bytes []byte, host Host, connection *Connection) bool {
		move := packets.NewMovePlayerPacket()
		move.Buffer = bytes
		move.DecodeHeader()
		move.Decode()

		if connection.IsClient(host.GetAddress()) {

			LastY = move.Position.Y

			connection.client.Position = move.Position
		}

		return false
	})
}