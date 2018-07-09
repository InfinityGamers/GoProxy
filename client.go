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
	Host // Server extends a host that has writable packets
	// Server's udp address
	addr net.UDPAddr
	// the main Proxy
	Proxy *Proxy
	// the Connection between the Server
	// and the Client
	Conn *Connection
	// a bool that is true if the Client
	// is connected with the Server/Proxy
	Connection bool
	// the UUID of the Client's player
	uuid uuid.UUID
	// the entity runtime id of
	// the Client's player
	//entityRuntimeID uint64
	// the Client's position in the level
	Position r3.Vector
}

// returns new Client host
// that has writable packets
func NewClient(proxy *Proxy, conn *Connection) *Client  {
	client := Client{}
	client.Proxy = proxy
	client.Conn = conn
	client.Connection = false
	return &client
}

// returns if this host is the client
// this is the client struct to it returns true
func (client Client) IsClient() bool {
	return true
}

// returns if this host is the client
// this is the client struct to it returns false
func (client Client) IsServer() bool {
	return false
}

// this function is from the host interface
// it writes a packet buffer to the Client
func (client Client) WritePacket(buffer []byte) {
	_, err := client.Proxy.UDPConn.WriteTo(buffer, &client.addr)
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
	datagram := client.Conn.pkHandler.DatagramBuilder.BuildFromPackets(packets)
	client.WritePacket(datagram.Buffer)
}

// this set's the udp address
// which is used to communicate with the Client
func (client *Client) SetAddress(addr net.UDPAddr) {
	client.addr = addr
}

// returns the Client's address as net.UDPAddr
func (client Client) GetAddress() net.UDPAddr {
	return client.addr
}

// set true if the Client is connected with the
// Server/Proxy
func (client *Client) SetConnected(b bool) {
	client.Connection = b
}

// returns true if the Client is connected with the
// Server/Proxy
func (client *Client) IsConnected() bool {
	return client.Connection
}

// sends the Client's player a string message
func (client *Client) SendMessage(m string) {
	text := packets.NewTextPacket()
	text.Message = m
	client.SendPacket(text)
}

// changes the Client's player game mode
// although it updates for the Client, it will
// not on the Server's side
func (client *Client) SetGameMode(g int32) {
	gm := packets.NewSetGamemodePacket()
	gm.GameMode = g
	client.Conn.Server.SendPacket(gm)
	client.SendPacket(gm)
}

// Set a screen title and subtitle to the Client
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
		client.SendMessage(BrightGreen + Prefix + Orange + "You are using " + Author + "'s Proxy version " + Version)
		client.SendMessage(BrightBlue + "==============================")
		client.SetTitle(BrightGreen + Author + "'s Proxy", Orange + Author + "'s Proxy version " + Version, 1, 1, 1)
	}()
}