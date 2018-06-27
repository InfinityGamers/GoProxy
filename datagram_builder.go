package mcpeproxy

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
// to the proxy's own sequence number order
func (db *DatagramBuilder) BuildFromPackets(packets []packets.IPacket) protocol.Datagram {
	datagram := db.BuildRawFromPackets(packets, 0, db.pkHandler.orders.sequenceNumber)
	db.pkHandler.orders.sequenceNumber++
	return datagram
}

// this builds a datagram from raw bytes
// and sets the order of the datagram
// to the proxy's own sequence number order
func (db *DatagramBuilder) BuildFromBuffer(packets []byte) protocol.Datagram {
	datagram := db.BuildRaw(append([][]byte{}, packets), 0, db.pkHandler.orders.sequenceNumber)
	db.pkHandler.orders.sequenceNumber++
	return datagram
}

// this builds a datagram from packets
// with the encode/decode function
// and sets the order of the datagram
// to the latest sequence number received
func (db *DatagramBuilder) BuildFromPacketsWithLatestSequence(packets []packets.IPacket) protocol.Datagram {
	datagram := db.BuildRawFromPackets(packets, 0, db.pkHandler.lastSequenceNumber)
	db.pkHandler.lastSequenceNumber++
	return datagram
}

// this builds a datagram from raw bytes
// and sets the order of the datagram
// to the latest sequence number received
func (db *DatagramBuilder) BuildFromBufferWithLatestSequence(packets []byte) protocol.Datagram {
	datagram := db.BuildRaw(append([][]byte{}, packets), 0, db.pkHandler.lastSequenceNumber)
	db.pkHandler.lastSequenceNumber++
	return datagram
}

// this builds a datagram from raw bytes
// and sets the reliability of the
// encapsulated packet to a custom one
// and the order of the datagram
// to a custom sequence number
func (db *DatagramBuilder) BuildRaw(packets [][]byte, reliability byte, sequenceNumber uint32) protocol.Datagram {
	batch := NewMinecraftPacketBatch()
	batch.rawPackets = packets
	batch.Encode()
	encapsulated := protocol.NewEncapsulatedPacket()
	encapsulated.Reliability = reliability
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
func (db *DatagramBuilder) BuildRawFromPackets(packets []packets.IPacket, reliability byte, sequenceNumber uint32) protocol.Datagram {
	batch := NewMinecraftPacketBatch()
	for _, p := range packets {
		batch.AddPacket(p)
	}
	batch.Encode()
	encapsulated := protocol.NewEncapsulatedPacket()
	encapsulated.Reliability = reliability
	encapsulated.Buffer = batch.Buffer
	datagram := protocol.NewDatagram()
	datagram.SequenceNumber = sequenceNumber
	datagram.AddPacket(encapsulated)
	datagram.Encode()
	return *datagram
}