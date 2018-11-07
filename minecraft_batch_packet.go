package goproxy

import (
	"github.com/irmine/binutils"
	"github.com/Irmine/GoMine/net/packets"
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

const (
	BatchId = 0xFE
)

type MinecraftPacketBatch struct {
	*binutils.Stream
	rawPackets [][]byte
}

func NewMinecraftPacketBatch() *MinecraftPacketBatch {
	var batch = &MinecraftPacketBatch{}
	batch.Stream = binutils.NewStream()
	return batch
}

func (batch *MinecraftPacketBatch) AddPacket(packet packets.IPacket) {
	packet.EncodeHeader()
	packet.Encode()
	batch.AddRawPacket(packet.GetBuffer())
}

func (batch *MinecraftPacketBatch) AddRawPacket(packet []byte) {
	batch.rawPackets = append(batch.rawPackets, packet)
}

func (batch *MinecraftPacketBatch) GetRawPackets() [][]byte {
	return batch.rawPackets
}

func (batch *MinecraftPacketBatch) Encode() {
	batch.ResetStream()
	batch.PutByte(BatchId)

	stream := binutils.NewStream()
	for _, p := range batch.rawPackets {
		stream.PutLengthPrefixedBytes(p)
	}

	// zlib compression
	buffer := bytes.Buffer{}
	writer := zlib.NewWriter(&buffer)
	writer.Write(stream.Buffer)
	writer.Close()

	batch.PutBytes(buffer.Bytes())
}

func (batch *MinecraftPacketBatch) Decode() {
	defer func() {
		if err := recover(); err != nil {
			Alert(err)
		}
	}()

	batchId := batch.GetByte()

	if batchId != BatchId {
		Notice("Packet received is not batch:", batchId)
		return
	}

	data := batch.Buffer[batch.Offset:]

	reader := bytes.NewReader(data)
	zlibReader, err := zlib.NewReader(reader)

	if err != nil {
		Alert("Zlib Reader Error:", err)
		return
	}

	if zlibReader == nil {
		Alert("Error while reading from zlib")
		return
	}

	zlibReader.Close()

	data, err = ioutil.ReadAll(zlibReader)

	if err != nil {
		Alert("ioutil Error:", err)
		return
	}

	batch.ResetStream()
	batch.SetBuffer(data)

	for !batch.Feof() {
		batch.rawPackets = append(batch.rawPackets, batch.GetLengthPrefixedBytes())
	}
}