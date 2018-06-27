package mcpeproxy

import (
	"github.com/Irmine/GoRakLib/protocol"
	"net"
)

const (
	ClientHandshakeId = 0x13
	ClientCancelConnectId = 0x15
)

const (
	// Maximum MTU size is the maximum packet size.
	// Any MTU size above this will get limited to the maximum.
	MaximumMTUSize = 1492
	// MinimumMTUSize is the minimum packet size.
	// Any MTU size below this will get set to the minimum.
	MinimumMTUSize = 400
)


// Indexes is used for the collection of indexes related to datagrams and encapsulated packets.
// It uses several maps and is therefore protected by a mutex.
type Indexes struct {
	splits       map[int16][]*protocol.EncapsulatedPacket
	splitCounts  map[int16]uint
	splitId 	 int16
}

// Orders is used to re-order the indexes of encapsulated packets
// this is to avoid message order conflicts with the client and server
type Orders struct { // not used currently
	SplitIndex uint
	OrderIndex uint32
	MessageIndex uint32
	// the last datagram sequence number
	sequenceNumber uint32
}

// PacketHandler handles each datagram and
// calls all registered packet handlers to handle the packet,
// once ready, sends it to the right host
type PacketHandler struct {
	// connection is the connection between the
	// server and the client
	conn *Connection
	// a map where handler functions are stored
	// parameters:
	// 1. the bytes received from the packet
	// 2. the host to where the packet is headed
	// 3. the connection between the server and client
	// return true if the packet should be cancelled
	handlers map[byte][]func([]byte, Host, *Connection) bool
	// used for the collection of indexes
	// related to datagrams and encapsulated packets.
	splits Indexes
	// used for order of encapsulated packets
	// related to datagrams and encapsulated packets.
	orders Orders
	// DatagramBuilds builds datagram
	// either from a packet or from raw bytes
	// datagrams are the communication method
	// for the RakNet protocol
	datagramBuilder DatagramBuilder
	// the number of the last datagram
	// sent by both the client and server
	lastSequenceNumber uint32
	// the last datagram received
	datagram *protocol.Datagram
	// the last encapsulated received
	encapsulated *protocol.EncapsulatedPacket
	// a bool that returns true if the
	// client is ready for packets
	ready bool
}

// returns a new packet handler
// PacketHandler handles each datagram to make sure
// it goes to the right place, once ready, it calls all registered
// packet handlers to handle the packet
func NewPacketHandler() *PacketHandler {
	h := PacketHandler{}
	h.handlers = make(map[byte][]func([]byte, Host, *Connection) bool)
	h.splits = Indexes{splits: make(map[int16][]*protocol.EncapsulatedPacket), splitCounts: make(map[int16]uint)}
	h.orders = Orders{0, 0, 0, 0}
	h.datagramBuilder = NewDatagramBuilder()
	h.datagramBuilder.pkHandler = &h
	return &h
}

// registers a packet handler with a certain packet id
// every function should have these parameters:
// 1. the bytes received from the packet
// 2. the host to where the packet is headed
// 3. the connection between the server and client
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
			if v(packet, host, pkHandler.conn) {
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

// handles encapsulated packet
// if the packet is batch it will call the packet handlers
// if not it will just continue on sending the datagram
func (pkHandler *PacketHandler) HandleEncapsulated(packet *protocol.EncapsulatedPacket, host Host, addr net.UDPAddr) bool {
	handled := false
	PkId := packet.Buffer[0]
	if PkId == BatchId {
		if pkHandler.ready {
			batch := NewMinecraftPacketBatch()
			batch.SetBuffer(packet.Buffer)
			batch.Decode()

			datagram := protocol.NewDatagram()
			batch2 := NewMinecraftPacketBatch()

			for _, pk := range batch.GetRawPackets() {
				pkId := pk[0]
				if !pkHandler.CallPacketHandlers(pkId, host, pk) {
					batch2.AddRawPacket(pk)
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
			datagram.SequenceNumber = pkHandler.lastSequenceNumber
			pkHandler.lastSequenceNumber++
			datagram.AddPacket(encap)
			datagram.Encode()
			host.WritePacket(datagram.Buffer)

			handled = true
		}
	}else if PkId == ClientHandshakeId {
		Notice(AnsiGreen + "Client has connected to the server.")
		pkHandler.FlowDatagram(host)
		pkHandler.conn.client.SendJoinMessage()
		pkHandler.ready = true
		handled = true
	} else if PkId == ClientCancelConnectId {
		Notice(AnsiBrightRed + "Client has disconnected from the server.")
		pkHandler.conn.client.SetConnected(false)
		pkHandler.FlowDatagram(host)
		handled = true
	} else {
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
// after the client has clicked the server
// of the proxy
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