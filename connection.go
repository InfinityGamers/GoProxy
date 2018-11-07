package goproxy

import (
	"net"
	"github.com/Irmine/GoRakLib/protocol"
)

const (
	IdUnconnectedPingOpenConnection = 0x01 //Client to Server
	IdUnconnectedPongOpenConnection = 0x1c //Server to Client as response to IdConnectedPingOpenConnection (0x01)
	IdOpenConnectionRequest1 = 0x05 //Client to Server
	IdOpenConnectionReply1 = 0x06 //Server to Client
	IdOpenConnectionRequest2 = 0x07 //Client to Server
	IdOpenConnectionReply2 = 0x08 //Server to Client
)

type Connection struct {
	Proxy      *Proxy
	Client     *Client
	Server     *Server
	pkHandler  *PacketHandler
}

func NewConnection(proxy *Proxy) Connection {
	conn := Connection{}
	conn.Proxy = proxy

	conn.Client = NewClient(proxy, &conn)
	conn.Server = NewServer(proxy, &conn)
	conn.pkHandler = NewPacketHandler()

	conn.pkHandler.Conn = &conn

	RegisterDefaultHandlers(conn.pkHandler)

	return conn
}

func (conn *Connection) GetClient() *Client {
	return conn.Client
}

func (conn *Connection) GetServer() *Server {
	return conn.Server
}

func (conn *Connection) GetPacketHandler() *PacketHandler {
	return conn.pkHandler
}

func (conn *Connection) IsClient(addr net.UDPAddr) bool {
	return conn.Client.addr.IP.Equal(addr.IP)
}

func (conn *Connection) IsServer(addr net.UDPAddr) bool {
	return conn.Server.GetAddress().IP.Equal(addr.IP)
}

func (conn *Connection) handleUnconnectedPing(addr net.UDPAddr, buffer []byte) {
	conn.Client.SetAddress(addr)
	conn.Server.WritePacket(buffer)

	//stringAddr := addr.IP.String()
	//Info("Received unconnected ping from Client address: " + stringAddr)
}

func (conn *Connection) handleUnconnectedPong(addr net.UDPAddr, buffer []byte) {
	pk := protocol.NewUnconnectedPong()
	pk.SetBuffer(buffer)
	pk.Decode()
	pk.Encode()

	conn.Client.WritePacket(pk.Buffer)
	conn.Server.Data.ParseFromString(pk.PongData)

	//stringAddr := addr.IP.String()
	//Info("Received unconnected pong from Server address: " + stringAddr)
}

func (conn *Connection) handleConnectionRequest1(addr net.UDPAddr, buffer []byte) {
	conn.Client.SetAddress(addr)
	conn.Server.WritePacket(buffer)

	//stringAddr := addr.IP.String()
	//
	//Info("Received Connection request 1 from Client address: " + stringAddr)
}

func (conn *Connection) handleConnectionReply1(addr net.UDPAddr, buffer []byte) {
	//stringAddr := addr.IP.String()
	//Info("Received Connection reply 1 from Server address: " + stringAddr)

	pk := protocol.NewOpenConnectionReply1()
	pk.SetBuffer(buffer)
	pk.Decode()
	pk.Encode()

	conn.Client.WritePacket(pk.Buffer)
}

func (conn *Connection) handleConnectionRequest2(addr net.UDPAddr, buffer []byte) {
	conn.Client.SetAddress(addr)
	conn.Server.WritePacket(buffer)

	//stringAddr := addr.IP.String()
	//
	//Info("Received Connection request 2 from Client address: " + stringAddr)
}

func (conn *Connection) handleConnectionReply2(addr net.UDPAddr, buffer []byte) {
	conn.Client.SetConnected(true)

	//stringAddr := addr.IP.String()
	//Info("Received Connection reply 2 from Server address: " + stringAddr)

	pk := protocol.NewOpenConnectionReply2()
	pk.SetBuffer(buffer)
	pk.Decode()

	pk.Encode()

	conn.Client.WritePacket(pk.Buffer)
}

func (conn *Connection) HandleIncomingPackets() {
	for {
		buffer := make([]byte, 2048)
		_, addr, err := conn.Proxy.UDPConn.ReadFromUDP(buffer)

		if err != nil {
			Alert(err.Error())
			continue
		}

		MessageId := buffer[0]

		if conn.Client.IsConnected() {
			if conn.IsServer(*addr) {
				conn.pkHandler.HandleIncomingPacket(buffer, conn.Client, *addr)
			} else {
				conn.pkHandler.HandleIncomingPacket(buffer, conn.Server, *addr)
			}
			continue
		}

		switch MessageId {
		case byte(IdUnconnectedPingOpenConnection):
			conn.handleUnconnectedPing(*addr, buffer)
			break
		case byte(IdUnconnectedPongOpenConnection):
			conn.handleUnconnectedPong(*addr, buffer)
			break
		case byte(IdOpenConnectionRequest1):
			conn.handleConnectionRequest1(*addr, buffer)
			break
		case byte(IdOpenConnectionReply1):
			conn.handleConnectionReply1(*addr, buffer)
			break
		case byte(IdOpenConnectionRequest2):
			conn.handleConnectionRequest2(*addr, buffer)
			break
		case byte(IdOpenConnectionReply2):
			conn.handleConnectionReply2(*addr, buffer)
			break
		}
	}
}