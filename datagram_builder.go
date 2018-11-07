package goproxy

import (
	"github.com/Irmine/GoMine/net/packets"
	"github.com/Irmine/GoRakLib/protocol"
)

// DatagramBuilds builds datagram
// either from a packet or from raw bytes
// datagrams are the communication method
// for the RakNet protocol
type DatagramBuilder struct {
	pkHandler *PacketHandler
}

func NewDatagramBuilder() DatagramBuilder {
	return DatagramBuilder{}
}

// this builds a datagram from packets
// with the encode/decode function
// and sets the order of the datagram
// to the Proxy's own sequence number order
func (db *DatagramBuilder) BuildFromPackets(packets []packets.IPacket) protocol.Datagram {
	db.pkHandler.orders.SequenceNumber++
	datagram := db.BuildRawFromPackets(packets, Unreliable, db.pkHandler.orders.SequenceNumber, 0, 0)
	return datagram
}

// this builds a datagram from raw bytes
// and sets the order of the datagram
// to the Proxy's own sequence number order
func (db *DatagramBuilder) BuildFromBuffer(packets []byte) protocol.Datagram {
	db.pkHandler.orders.SequenceNumber++
	datagram := db.BuildRaw(append([][]byte{}, packets), Unreliable, db.pkHandler.orders.SequenceNumber, 0, 0)
	return datagram
}

// this builds a datagram from packets
// with the encode/decode function
// and sets the order of the datagram
// to the latest sequence number received
func (db *DatagramBuilder) BuildFromPacketsWithLatestSequence(packets []packets.IPacket) protocol.Datagram {
	db.pkHandler.lastSequenceNumber++
	datagram := db.BuildRawFromPackets(packets, Unreliable, db.pkHandler.lastSequenceNumber, 0, 0)
	return datagram
}

// this builds a datagram from raw bytes
// and sets the order of the datagram
// to the latest sequence number received
func (db *DatagramBuilder) BuildFromBufferWithLatestSequence(packets []byte) protocol.Datagram {
	db.pkHandler.lastSequenceNumber++
	datagram := db.BuildRaw(append([][]byte{}, packets), Unreliable, db.pkHandler.lastSequenceNumber, 0, 0)
	return datagram
}


// this builds a datagram from a packet batch
// and sets the reliability of the encapsulated
// packet to a custom one and the order of the
// datagram to a custom sequence number
func (db *DatagramBuilder) BuildFromBatch(batch *MinecraftPacketBatch, reliability byte, sequenceNumber, orderIndex, messageIndex uint32) protocol.Datagram {
	encapsulated := protocol.NewEncapsulatedPacket()
	encapsulated.Reliability = reliability
	encapsulated.OrderIndex = orderIndex
	encapsulated.MessageIndex = messageIndex
	encapsulated.Buffer = batch.Buffer
	datagram := protocol.NewDatagram()
	datagram.SequenceNumber = sequenceNumber
	datagram.AddPacket(encapsulated)
	datagram.Encode()
	return *datagram
}

// this builds a datagram from a slice
// of encapsulated packets and sets the
// reliability of the encapsulated packet
// to a custom one and the order of the
// datagram to a custom sequence number
func (db *DatagramBuilder) BuildFromEncapsulated(packets []*protocol.EncapsulatedPacket, sequenceNumber uint32) protocol.Datagram {
	datagram := protocol.NewDatagram()
	datagram.SequenceNumber = sequenceNumber
	for _, p := range packets {
		datagram.AddPacket(p)
	}
	datagram.Encode()
	return *datagram
}

// this builds a datagram from raw bytes
// and sets the reliability of the
// encapsulated packet to a custom one
// and the order of the datagram
// to a custom sequence number
func (db *DatagramBuilder) BuildRaw(packets [][]byte, reliability byte, sequenceNumber, orderIndex, messageIndex uint32) protocol.Datagram {
	batch := NewMinecraftPacketBatch()
	batch.rawPackets = packets
	batch.Encode()
	encapsulated := protocol.NewEncapsulatedPacket()
	encapsulated.Reliability = reliability
	encapsulated.OrderIndex = orderIndex
	encapsulated.MessageIndex = messageIndex
	encapsulated.Buffer = batch.Buffer
	datagram := protocol.NewDatagram()
	datagram.SequenceNumber = sequenceNumber
	datagram.AddPacket(encapsulated)
	datagram.Encode()
	return *datagram
}

// this builds a datagram from packets
// with the encode/decode function
// and sets the reliability of the
// encapsulated packet to a custom one
// and the order of the datagram
// to a custom sequence number
func (db *DatagramBuilder) BuildRawFromPackets(packets []packets.IPacket, reliability byte, sequenceNumber, orderIndex, messageIndex uint32) protocol.Datagram {
	batch := NewMinecraftPacketBatch()
	for _, p := range packets {
		batch.AddPacket(p)
	}
	batch.Encode()
	encapsulated := protocol.NewEncapsulatedPacket()
	encapsulated.HasSplit = false
	encapsulated.Reliability = reliability
	encapsulated.OrderIndex = orderIndex
	encapsulated.MessageIndex = messageIndex
	encapsulated.Buffer = batch.Buffer
	datagram := protocol.NewDatagram()
	datagram.SequenceNumber = sequenceNumber
	datagram.AddPacket(encapsulated)
	datagram.Encode()
	return *datagram
}