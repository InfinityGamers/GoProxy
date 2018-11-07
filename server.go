package goproxy

import (
	"net"
	"github.com/Irmine/GoMine/net/packets"
)

type Server struct {
	Host // Server extends a host that has writable packets
	// Server's udp address
	addr *net.UDPAddr
	// the main Proxy
	Proxy *Proxy
	// the Connection between the Server
	// and the Client
	Conn *Connection
	// Server Data contains the Server's query info such as:
	// Server name, version, protocol, players, software, etc...
	Data ServerData
	// the Server's uid
	serverId int64
}


// returns new Server host
// that has writable packets
func NewServer(proxy *Proxy, conn *Connection) *Server  {
	server := Server{}
	server.Proxy = proxy
	server.Conn = conn
	server.Data = NewServerData()
	return &server
}

// returns if this host is the client
// this is the server struct to it returns false
func (server Server) IsClient() bool {
	return false
}

// returns if this host is the client
// this is the server struct to it returns true
func (server Server) IsServer() bool {
	return true
}

// this function is from the host interface
// it writes a packet buffer to the Server
func (server Server) WritePacket(buffer []byte) {
	_, err := server.Proxy.WriteToUDP(buffer, server.addr)
	if err != nil {
		Alert(err.Error())
	}
}

// this function is from the host interface
// it sends a single packet
func (server Server) SendPacket(packet packets.IPacket) {
	server.SendBatchPacket([]packets.IPacket{packet})
}

// this function is from the host interface
// it sends a batch packet:
// a packet with multiple packets inside
func (server Server) SendBatchPacket(packets []packets.IPacket) {
	datagram := server.Conn.pkHandler.DatagramBuilder.BuildFromPackets(packets)
	server.WritePacket(datagram.Buffer)
}

// this set's the udp address
// which is used to communicate with the Server
func (server *Server) SetAddress(addr *net.UDPAddr) {
	server.addr = addr
}

// returns the Server's address as net.UDPAddr
func (server Server) GetAddress() net.UDPAddr {
	return *server.addr
}

// this returns the Server's Data struct
func (server *Server) GetServerData() ServerData {
	return server.GetServerData()
}