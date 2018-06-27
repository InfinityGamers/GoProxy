package mcpeproxy

import (
	"github.com/infinitygamers/mcpeproxy/info"
	"github.com/infinitygamers/mcpeproxy/packets"
	"strings"
	"time"
	"strconv"
	packets2 "github.com/Irmine/GoMine/net/packets"
)

var GameMode int32 = 0
var LastY float64 = 0

var Hit bool

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
				return true
			}

			if strings.ToLower(text.Message) == "hit" {

				if Hit {
					Hit = false
					connection.client.SendMessage(Green + Prefix + Orange + "Starting kill aura.")
				}else{
					Hit = true
					connection.client.SendMessage(Green + Prefix + Orange + "Stopping kill aura.")
				}

				if Hit {
					go func() {
						for Hit {
							time.Sleep(time.Millisecond * 250)
							pk := packets.NewInteractPacket()
							pk.Action = packets.LeftClick
							d := pkHandler.datagramBuilder.BuildFromPacketsWithLatestSequence([]packets2.IPacket{pk})
							connection.client.WritePacket(d.Buffer)
						}
					}()
				}

				return true
			}
		}
		return false
	})
}