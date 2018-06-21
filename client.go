package mcpeproxy

import (
	"net"
	"github.com/infinitygamers/mcpeproxy/packets"
	"github.com/google/uuid"
	packets2 "github.com/Irmine/GoMine/net/packets"
	"time"
	"github.com/golang/geo/r3"
)

type Client struct {
	Host // server extends a host that has writable packets
	// server's udp address
	addr net.UDPAddr
	// the main proxy
	proxy *Proxy
	// the connection between the server
	// and the client
	conn *Connection
	// a bool that is true if the client
	// is connected with the server/proxy
	connection bool
	// the UUID of the client's player
	uuid uuid.UUID
	// the entity runtime id of
	// the client's player
	//entityRuntimeID uint64
	// the client's position in the level
	Position r3.Vector
}

// returns new client host
// that has writable packets
func NewClient(proxy *Proxy, conn *Connection) *Client  {
	client := Client{}
	client.proxy = proxy
	client.conn = conn
	client.connection = false
	return &client
}

// this function is from the host interface
// it writes a packet buffer to the client
func (client Client) WritePacket(buffer []byte) {
	_, err := client.proxy.UDPConn.WriteTo(buffer, &client.addr)
	if err != nil {
		Alert(err.Error())
	}
}

// this function is from the host interface
// This sends a single packet
func (client Client) SendPacket(packet packets2.IPacket) {
	client.SendBatchPacket([]packets2.IPacket{packet})
}

// this function is from the host interface
// it sends a batch packet:
// a packet with multiple packets inside
func (client Client) SendBatchPacket(packets []packets2.IPacket) {
	datagram := client.conn.pkHandler.datagramBuilder.BuildFromPackets(packets)
	client.WritePacket(datagram.Buffer)
}

// this set's the udp address
// which is used to communicate with the client
func (client *Client) SetAddress(addr net.UDPAddr) {
	client.addr = addr
}

// returns the client's address as net.UDPAddr
func (client Client) GetAddress() net.UDPAddr {
	return client.addr
}

// set true if the client is connected with the
// server/proxy
func (client *Client) SetConnected(b bool) {
	client.connection = b
}

// returns true if the client is connected with the
// server/proxy
func (client *Client) IsConnected() bool {
	return client.connection
}

// sends the client's player a string message
func (client *Client) SendMessage(m string) {
	text := packets.NewTextPacket()
	text.Message = m
	client.SendPacket(text)
}

// changes the client's player game mode
// although it updates for the client, it will
// not on the server's side
func (client *Client) SetGameMode(g int32) {
	gm := packets.NewSetGamemodePacket()
	gm.GameMode = g
	client.conn.server.SendPacket(gm)
	client.SendPacket(gm)
}

// Set a screen title and subtitle to the client
func (client *Client) SetTitle(title, subtitle string, fadeInTime, stayTime, fadeOutTime int32) {
	// title
	t := packets.NewSetTitlePacket()
	t.TitleType = packets.SetTitle
	t.Text = title
	t.FadeInTime = fadeInTime
	t.StayTime = stayTime
	t.FadeOutTime = fadeOutTime

	// subtitle
	t2 := packets.NewSetTitlePacket()
	t2.TitleType = packets.SetSubtitle
	t2.Text = subtitle
	t2.FadeInTime = fadeInTime
	t2.StayTime = stayTime
	t2.FadeOutTime = fadeOutTime

	client.SendBatchPacket([]packets2.IPacket{t, t2})
}

func (client *Client) SendJoinMessage() {
	go func() {
		time.Sleep(5 * time.Second)
		client.SendMessage(BrightBlue + "==============================")
		client.SendMessage(BrightGreen + Prefix + Orange + "You are using " + Author + "'s proxy version " + Version)
		client.SendMessage(BrightBlue + "==============================")
		client.SetTitle(BrightGreen + Author + "'s Proxy", Orange + Author + "'s proxy version " + Version, 1, 1, 1)
	}()
}