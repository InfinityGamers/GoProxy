package goproxy

import (
	"strings"
	"strconv"
)

// this struct holds the Server Data
type ServerData struct {
	// raw Server Data as mixed slice
	RawData []interface{}
	// Server name as string
	serverName string
	// Server protocol version as int
	protocolVersion int
	// Server game version as int
	gameVersion int
	// Server online player count as int
	onlinePlayerCount int //ignore
	// Server nax player count as int
	maxPlayerCount int //ignore
}

func NewServerData() ServerData {
	return ServerData{}
}

// This parses Data sent from the Server
// the Data is send in this format:
// MCPE;Server name;protocol version;MCPE version;online player count;max player count
func (svData *ServerData) ParseFromString(data string){
	split := strings.Split(data, ";")[1:]
	svData.serverName = split[0]
	svData.protocolVersion, _ = strconv.Atoi(split[1])
	svData.gameVersion, _ = strconv.Atoi(split[2])
	svData.onlinePlayerCount, _ = strconv.Atoi(split[3])
	svData.maxPlayerCount, _ = strconv.Atoi(split[4])
}

// returns the Server name sent from the Server
func (svData *ServerData) GetServerName() string {
	return svData.serverName
}

// returns the protocol version sent from the Server
func (svData *ServerData) GetProtocolVersion() int {
	return svData.protocolVersion
}

// returns the game version sent from the Server
func (svData *ServerData) GetGameVersion() int {
	return svData.gameVersion
}