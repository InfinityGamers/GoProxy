package mcpeproxy

import (
	"github.com/infinitygamers/mcpeproxy/info"
	"github.com/infinitygamers/mcpeproxy/packets"
	"strconv"
	"strings"
)

var GameMode int32 = 0

//var Geometry = `{"geometry.polo128":{"texturewidth":128,"textureheight":128,"bones":[{"name":"head","pivot":[0,4,0],"cubes":[{"origin":[0,0,0],"size":[128,128,1],"uv":[0,0]}]}]}}`

func GetNextGameMode(toggle bool) int32 {
	if toggle {
		if GameMode == 1 {
			GameMode = 0
		} else {
			GameMode = 1
		}
	}else{
		GameMode++
		if GameMode > 3 {
			GameMode = 0
		}
	}
	return GameMode
}

func RegisterDefaultHandlers(pkHandler *PacketHandler)  {
	pkHandler.RegisterHandler(info.MovePlayerPacket, func(bytes []byte, host Host, connection *Connection) bool {
		if connection.IsServer(host.GetAddress()) {
			move := packets.NewMovePlayerPacket()
			move.Buffer = bytes
			move.DecodeHeader()
			move.Decode()

			connection.client.Position = move.Position
		}
		return false
	})
	pkHandler.RegisterHandler(info.TextPacket, func(bytes []byte, host Host, connection *Connection) bool {
		if connection.IsServer(host.GetAddress()) {
			text := packets.NewTextPacket()
			text.Buffer = bytes
			text.DecodeHeader()
			text.Decode()

			if strings.ToLower(text.Message) == "fly" {
				gm := GetNextGameMode(true)
				connection.client.SetGameMode(gm)
				connection.client.SendMessage(Green + Prefix + Orange + "Set game mode to " + strconv.Itoa(int(gm)))
			}
		}
		return false
	})
}