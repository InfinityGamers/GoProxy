package mcpeproxy

import (
	"strings"
	"strconv"
)

// this struct holds the server data
type ServerData struct {
	// raw server data as mixed slice
	RawData []interface{}
	// server name as string
	serverName string
	// server protocol version as int
	protocolVersion int
	// server game version as int
	gameVersion int
	// server online player count as int
	onlinePlayerCount int //ignore
	// server nax player count as int
	maxPlayerCount int //ignore
}

func NewServerData() ServerData {
	return ServerData{}
}

// This parses data sent from the server
// the data is send in this format:
// MCPE;server name;protocol version;MCPE version;online player count;max player count
func (svData *ServerData) ParseFromString(data string){
	split := strings.Split(data, ";")[1:]
	svData.serverName = split[0]
	svData.protocolVersion, _ = strconv.Atoi(split[1])
	svData.gameVersion, _ = strconv.Atoi(split[2])
	svData.onlinePlayerCount, _ = strconv.Atoi(split[3])
	svData.maxPlayerCount, _ = strconv.Atoi(split[4])
}

// returns the server name sent from the server
func (svData *ServerData) GetServerName() string {
	return svData.serverName
}

// returns the protocol version sent from the server
func (svData *ServerData) GetProtocolVersion() int {
	return svData.protocolVersion
}

// returns the game version sent from the server
func (svData *ServerData) GetGameVersion() int {
	return svData.gameVersion
}