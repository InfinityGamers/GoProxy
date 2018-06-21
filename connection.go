package mcpeproxy

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
	proxy *Proxy
	client *Client
	server *Server
	pkHandler *PacketHandler
}

func NewConnection(proxy *Proxy) Connection {
	conn := Connection{}
	conn.proxy = proxy

	conn.client = NewClient(proxy, &conn)
	conn.server = NewServer(proxy, &conn)
	conn.pkHandler = NewPacketHandler()

	conn.pkHandler.conn = &conn

	RegisterDefaultHandlers(conn.pkHandler)

	return conn
}

func (conn *Connection) GetClient() *Client {
	return conn.client
}

func (conn *Connection) GetServer() *Server {
	return conn.server
}

func (conn *Connection) GetPacketHandler() *PacketHandler {
	return conn.pkHandler
}

func (conn *Connection) IsClient(addr net.UDPAddr) bool {
	return conn.client.addr.IP.Equal(addr.IP)
}

func (conn *Connection) IsServer(addr net.UDPAddr) bool {
	return conn.server.GetAddress().IP.Equal(addr.IP)
}

func (conn *Connection) handleUnconnectedPing(addr net.UDPAddr, buffer []byte) {
	conn.client.SetAddress(addr)
	conn.server.WritePacket(buffer)

	//stringAddr := addr.IP.String()
	//Info("Received unconnected ping from client address: " + stringAddr)
}

func (conn *Connection) handleUnconnectedPong(addr net.UDPAddr, buffer []byte) {
	pk := protocol.NewUnconnectedPong()
	pk.SetBuffer(buffer)
	pk.Decode()
	pk.Encode()

	conn.client.WritePacket(pk.Buffer)
	conn.server.data.ParseFromString(pk.PongData)

	//stringAddr := addr.IP.String()
	//Info("Received unconnected pong from server address: " + stringAddr)
}

func (conn *Connection) handleConnectionRequest1(addr net.UDPAddr, buffer []byte) {
	conn.client.SetAddress(addr)
	conn.server.WritePacket(buffer)

	//stringAddr := addr.IP.String()
	//
	//Info("Received connection request 1 from client address: " + stringAddr)
}

func (conn *Connection) handleConnectionReply1(addr net.UDPAddr, buffer []byte) {
	//stringAddr := addr.IP.String()
	//Info("Received connection reply 1 from server address: " + stringAddr)

	pk := protocol.NewOpenConnectionReply1()
	pk.SetBuffer(buffer)
	pk.Decode()
	pk.Encode()

	conn.client.WritePacket(pk.Buffer)
}

func (conn *Connection) handleConnectionRequest2(addr net.UDPAddr, buffer []byte) {
	conn.client.SetAddress(addr)
	conn.server.WritePacket(buffer)

	//stringAddr := addr.IP.String()
	//
	//Info("Received connection request 2 from client address: " + stringAddr)
}

func (conn *Connection) handleConnectionReply2(addr net.UDPAddr, buffer []byte) {
	conn.client.SetConnected(true)

	//stringAddr := addr.IP.String()
	//Info("Received connection reply 2 from server address: " + stringAddr)

	pk := protocol.NewOpenConnectionReply2()
	pk.SetBuffer(buffer)
	pk.Decode()

	pk.Encode()

	conn.client.WritePacket(pk.Buffer)
}

func (conn *Connection) HandleIncomingPackets() {
	for true {
		buffer := make([]byte, 2048)
		_, addr, err := conn.proxy.UDPConn.ReadFromUDP(buffer)

		if err != nil {
			Alert(err.Error())
			continue
		}

		MessageId := buffer[0]

		if conn.client.IsConnected() {
			if conn.IsServer(*addr) {
				conn.pkHandler.HandleIncomingPacket(buffer, conn.client, *addr)
			} else {
				conn.pkHandler.HandleIncomingPacket(buffer, conn.server, *addr)
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