package mcpeproxy

import (
	"net"
	"github.com/Irmine/GoMine/net/packets"
)

type Server struct {
	Host // server extends a host that has writable packets
	// server's udp address
	addr *net.UDPAddr
	// the main proxy
	proxy *Proxy
	// the connection between the server
	// and the client
	conn *Connection
	// server data contains the server's query info such as:
	// server name, version, protocol, players, software, etc...
	data ServerData
	// the server's uid
	serverId int64
}


// returns new server host
// that has writable packets
func NewServer(proxy *Proxy, conn *Connection) *Server  {
	server := Server{}
	server.proxy = proxy
	server.conn = conn
	server.data = NewServerData()
	return &server
}

// this function is from the host interface
// it writes a packet buffer to the server
func (server Server) WritePacket(buffer []byte) {
	_, err := server.proxy.WriteToUDP(buffer, server.addr)
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
	datagram := server.conn.pkHandler.datagramBuilder.BuildFromPackets(packets)
	server.WritePacket(datagram.Buffer)
}

// this set's the udp address
// which is used to communicate with the server
func (server *Server) SetAddress(addr *net.UDPAddr) {
	server.addr = addr
}

// returns the server's address as net.UDPAddr
func (server Server) GetAddress() net.UDPAddr {
	return *server.addr
}

// this returns the server's data struct
func (server *Server) GetServerData() ServerData {
	return server.GetServerData()
}