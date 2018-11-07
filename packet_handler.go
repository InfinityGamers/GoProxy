package goproxy

import (
	"net"
	"github.com/Irmine/GoRakLib/protocol"
	"github.com/Irmine/GoMine/net/packets"
)

const (
	ClientHandshakeId = 0x13
	ClientCancelConnectId = 0x15
)

// Indexes is used for the collection of indexes related to datagrams and encapsulated packets.
// It uses several maps and is therefore protected by a mutex.
type Indexes struct {
	splits       map[int16][]*protocol.EncapsulatedPacket
	splitCounts  map[int16]uint
	splitId 	 int16
}

// Orders is used to re-order the indexes of encapsulated packets
// this is to avoid message order conflicts with the Client and Server
type Orders struct {
	//SplitIndex uint

	// the last encapsulated order
	// index number from both hosts
	///OrderIndex uint32

	// the last encapsulated
	// message index number from both hosts
	//MessageIndex uint32

	// the last datagram sequence
	// number from both hosts
	SequenceNumber uint32
}

// PacketHandler handles each datagram and
// calls all registered packet handlers to handle the packet,
// once ready, sends it to the right host
type PacketHandler struct {
	// Connection is the Connection between the
	// Server and the Client
	Conn *Connection
	// a map where handler functions are stored
	// parameters:
	// 1. the bytes received from the packet
	// 2. the host to where the packet is headed
	// 3. the Connection between the Server and Client
	// return true if the packet should be cancelled
	handlers map[byte][]func([]byte, Host, *Connection) bool
	// used for the collection of indexes
	// related to datagrams and encapsulated packets.
	splits Indexes
	// Orders is used to re-order
	// the indexes of encapsulated packets
	// this is to avoid message order
	// conflicts with the Client and Server
	orders Orders
	// DatagramBuilds builds datagram
	// either from a packet or from raw bytes
	// datagrams are the communication method
	// for the RakNet protocol
	DatagramBuilder DatagramBuilder
	// the number of the last datagram
	// sent by both hosts
	lastSequenceNumber uint32
	// the last datagram received from both hosts
	datagram *protocol.Datagram
	// a slice of raw packets that will be
	// sent out to the client next tick, packets are sent
	// every tick via the packet handler,
	// 20 ticks = 1 second, 1 tick ~ 0.05 seconds
	OutBoundPacketsClient [][]byte
	// a slice of raw packets that will be
	// sent out to the server next tick, packets are sent
	// every tick via the packet handler,
	// 20 ticks = 1 second, 1 tick ~ 0.05 seconds
	OutBoundPacketsServer [][]byte
	// a bool that returns true if the
	// Client is ready for packets
	Ready bool
}

// returns a new packet handler
// PacketHandler handles each datagram and
// calls all registered packet handlers to handle the packet,
// once ready, sends it to the right host
func NewPacketHandler() *PacketHandler {
	h := PacketHandler{}
	h.handlers = make(map[byte][]func([]byte, Host, *Connection) bool)
	h.splits = Indexes{splits: make(map[int16][]*protocol.EncapsulatedPacket), splitCounts: make(map[int16]uint)}
	h.DatagramBuilder = NewDatagramBuilder()
	h.DatagramBuilder.pkHandler = &h
	h.OutBoundPacketsClient = [][]byte{}
	h.OutBoundPacketsServer = [][]byte{}
	h.orders = Orders{}
	return &h
}

// registers a packet handler with a certain packet id
// every function should have these parameters:
// 1. the bytes received from the packet
// 2. the host to where the packet is headed
// 3. the Connection between the Server and Client
// return true if the packet should be cancelled
func (pkHandler *PacketHandler) RegisterHandler(pkId byte, f func([]byte, Host, *Connection) bool) {
	if _, ok := pkHandler.handlers[pkId]; !ok {
		pkHandler.handlers[pkId] = []func([]byte, Host, *Connection) bool {}
	}
	pkHandler.handlers[pkId] = append(pkHandler.handlers[pkId], f)
}

// calls packet handlers for a certain packet id
// returns true if some handler cancelled the packet
func (pkHandler *PacketHandler) CallPacketHandlers(pkId byte, host Host, packet []byte) bool {
	cancelled := false
	if _, ok := pkHandler.handlers[pkId]; ok {
		for _, v := range pkHandler.handlers[pkId] {
			if v(packet, host, pkHandler.Conn) {
				cancelled = true
			}
		}
	}
	return cancelled //return true if packet is cancelled
}

// Sends the last datagram received from one host
// to the other host
func (pkHandler *PacketHandler) FlowDatagram(host Host) {
	datagram := pkHandler.datagram
	datagram.Encode()
	host.WritePacket(datagram.Buffer)
}

// adds a packet to the slice of
// packets that will be sent out to the client next tick,
// packets are sent every tick via the packet handler,
// 20 ticks = 1 second, 1 tick ~ 0.05 seconds
func (pkHandler *PacketHandler) AddOutboundPacketClient(packet packets.IPacket)  {
	packet.EncodeHeader()
	packet.Encode()
	pkHandler.OutBoundPacketsClient = append(pkHandler.OutBoundPacketsClient, packet.GetBuffer())
}

// adds a raw packet to the slice of
// packets that will be sent out to the client next tick,
// packets are sent every tick via the packet handler,
// 20 ticks = 1 second, 1 tick ~ 0.05 seconds
func (pkHandler *PacketHandler) AddOutboundRawPacketClient(packet []byte)  {
	pkHandler.OutBoundPacketsClient = append(pkHandler.OutBoundPacketsClient, packet)
}

// adds a packet to the slice of
// packets that will be sent out to the server next tick,
// packets are sent every tick via the packet handler,
// 20 ticks = 1 second, 1 tick ~ 0.05 seconds
func (pkHandler *PacketHandler) AddOutboundPacketServer(packet packets.IPacket)  {
	packet.EncodeHeader()
	packet.Encode()
	pkHandler.OutBoundPacketsServer = append(pkHandler.OutBoundPacketsServer, packet.GetBuffer())
}

// adds a raw packet to the slice of
// packets that will be sent out to the server next tick,
// packets are sent every tick via the packet handler,
// 20 ticks = 1 second, 1 tick ~ 0.05 seconds
func (pkHandler *PacketHandler) AddOutboundRawPacketServer(packet []byte)  {
	pkHandler.OutBoundPacketsServer = append(pkHandler.OutBoundPacketsServer, packet)
}

// handles encapsulated packet
// if the packet is batch it will call the packet handlers
// if not it will just continue on sending the datagram
func (pkHandler *PacketHandler) HandleEncapsulated(packet *protocol.EncapsulatedPacket, host Host, addr net.UDPAddr) bool {
	handled := false
	PkId := packet.Buffer[0]
	if PkId == BatchId {
		if pkHandler.Ready {
			batch := NewMinecraftPacketBatch()
			batch.SetBuffer(packet.Buffer)
			batch.Decode()
			batch2 := NewMinecraftPacketBatch()

			for _, pk := range batch.GetRawPackets() {
				pkId := pk[0]
				if !pkHandler.CallPacketHandlers(pkId, host, pk) {
					batch2.AddRawPacket(pk)
				}
			}

			if host.IsServer() {
				if toServer := pkHandler.OutBoundPacketsServer; len(toServer) > 0 {
					for _, pk := range toServer {
						batch2.AddRawPacket(pk)
					}
					pkHandler.OutBoundPacketsServer = nil
				}
			}else{
				if toClient := pkHandler.OutBoundPacketsClient; len(toClient) > 0 {
					for _, pk := range toClient {
						batch2.AddRawPacket(pk)
					}
					pkHandler.OutBoundPacketsClient = nil
				}
			}

			batch2.Encode()
			encap := protocol.NewEncapsulatedPacket()
			encap.Buffer = batch2.Buffer
			encap.Reliability = packet.Reliability
			encap.HasSplit = false
			encap.Length = packet.Length
			encap.MessageIndex = packet.MessageIndex
			encap.OrderIndex = packet.OrderIndex
			dgram := protocol.NewDatagram()
			dgram.AddPacket(encap)
			dgram.SequenceNumber = pkHandler.lastSequenceNumber
			dgram.Encode()

			host.WritePacket(dgram.Buffer)

			handled = true
		}

	}else if PkId == ClientHandshakeId {
		Info(AnsiGreen + "Client has connected to the Server.")
		pkHandler.FlowDatagram(host)
		pkHandler.Conn.Client.SendJoinMessage()
		pkHandler.Ready = true
		handled = true
	} else if PkId == ClientCancelConnectId {
		Info(AnsiBrightRed + "Client has disconnected from the Server.")
		pkHandler.Conn.Client.SetConnected(false)
		pkHandler.FlowDatagram(host)
		pkHandler.Ready = false
		handled = true
	} else {

		switch PkId {
		case protocol.IdConnectionRequest:
			Info(AnsiBrightCyan + "Connection request from client")
		case protocol.IdConnectionAccept:
			Info(AnsiBrightCyan + "Connection request accepted from server")
		}

		pkHandler.FlowDatagram(host)
		handled = true
	}
	return handled
}

// Decodes a split encapsulated packet
// and turns it into one
func (pkHandler *PacketHandler) DecodeSplit(packet *protocol.EncapsulatedPacket) *protocol.EncapsulatedPacket {
	id := packet.SplitId

	if pkHandler.splits.splits[id] == nil {
		pkHandler.splits.splits[id] = make([]*protocol.EncapsulatedPacket, packet.SplitCount)
		pkHandler.splits.splitCounts[id] = 0
	}

	if pk := pkHandler.splits.splits[id][packet.SplitIndex]; pk == nil {
		pkHandler.splits.splitCounts[id]++
	}

	pkHandler.splits.splits[id][packet.SplitIndex] = packet

	if pkHandler.splits.splitCounts[id] == packet.SplitCount {
		newPacket := protocol.NewEncapsulatedPacket()
		for _, pk := range pkHandler.splits.splits[id] {
			newPacket.PutBytes(pk.Buffer)
		}

		delete(pkHandler.splits.splits, id)
		return newPacket
	}

	return nil
}

// this handles encapsulated packets that
// come in more than one part, construct a single one
// and calls the HandleEncapsulated function
func (pkHandler *PacketHandler) HandleSplitEncapsulated(packet *protocol.EncapsulatedPacket, host Host, addr net.UDPAddr) bool {
	d := pkHandler.DecodeSplit(packet)

	if d != nil {
		return pkHandler.HandleEncapsulated(d, host, addr)
	}

	return false
}

// this handles all the incoming packets
// after the Client has clicked the Server
// of the Proxy
func (pkHandler *PacketHandler) HandleIncomingPacket(buffer []byte, host Host, addr net.UDPAddr) {
	MessageId := buffer[0]
	if (MessageId & protocol.BitFlagValid) != 0 {
		if (MessageId & protocol.BitFlagIsNak) != 0 {
			nak := protocol.NewNAK()
			nak.SetBuffer(buffer)
			nak.Decode()
			nak.Encode()
			host.WritePacket(nak.Buffer)
		} else if (MessageId & protocol.BitFlagIsAck) != 0 {
			ack := protocol.NewACK()
			ack.SetBuffer(buffer)
			ack.Decode()
			ack.Encode()
			host.WritePacket(ack.Buffer)
		} else {
			datagram := protocol.NewDatagram()
			datagram.SetBuffer(buffer)
			datagram.Decode()

			pkHandler.datagram = datagram
			pkHandler.lastSequenceNumber = datagram.SequenceNumber
			packets2 := datagram.GetPackets()

			if len(*packets2) != 0 {
				for _, pk := range *packets2 {
					if pk.HasSplit {
						if !pkHandler.HandleSplitEncapsulated(pk, host, addr) {
							pkHandler.FlowDatagram(host)
						}
					} else {
						if !pkHandler.HandleEncapsulated(pk, host, addr) {
							pkHandler.FlowDatagram(host)
						}
					}
				}
			}else{
				pkHandler.FlowDatagram(host)
			}
		}
	} else {
		host.WritePacket(buffer)
	}
}