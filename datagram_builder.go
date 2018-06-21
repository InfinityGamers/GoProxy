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

func (db *DatagramBuilder) BuildFromPackets(packets []packets.IPacket) protocol.Datagram {
	batch := NewMinecraftPacketBatch()
	for _, p := range packets {
		batch.AddPacket(p)
	}
	batch.Encode()
	encapsulated := protocol.NewEncapsulatedPacket()
	encapsulated.Reliability = 0
	encapsulated.Buffer = batch.Buffer
	datagram := protocol.NewDatagram()
	datagram.SequenceNumber = db.pkHandler.orders.sequenceNumber
	db.pkHandler.orders.sequenceNumber++
	datagram.AddPacket(encapsulated)
	datagram.Encode()
	return *datagram
}

func (db *DatagramBuilder) BuildFromBuffer(packets []byte) protocol.Datagram {
	batch := NewMinecraftPacketBatch()
	batch.rawPackets = [][]byte{packets}
	batch.Encode()
	encapsulated := protocol.NewEncapsulatedPacket()
	encapsulated.Reliability = 0
	encapsulated.Buffer = batch.Buffer
	datagram := protocol.NewDatagram()
	datagram.SequenceNumber = db.pkHandler.orders.sequenceNumber
	db.pkHandler.orders.sequenceNumber++
	datagram.AddPacket(encapsulated)
	datagram.Encode()
	return *datagram
}

func (db *DatagramBuilder) BuildRaw(packets []byte, reliability byte, sequenceNumber uint32) protocol.Datagram {
	batch := NewMinecraftPacketBatch()
	batch.rawPackets = [][]byte{packets}
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