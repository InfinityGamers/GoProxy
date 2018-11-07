package goproxy

const (
	Unreliable = iota + 0x00
	UnreliableSequenced
	Reliable
	ReliableOrdered
	ReliableSequenced
	UnreliableWithAckReceipt
	ReliableWithAckReceipt
	ReliableOrderedWithAckReceipt
)

func isReliable(reliability byte) bool {
	return reliability == Reliable || reliability == ReliableOrdered || reliability == ReliableSequenced || reliability == ReliableWithAckReceipt || reliability == ReliableOrderedWithAckReceipt
}

func isSequenced(reliability byte) bool  {
	return reliability == UnreliableSequenced || reliability == ReliableSequenced
}

func isOrdered(reliability byte) bool {
	return reliability == ReliableOrdered || reliability == ReliableOrderedWithAckReceipt
}

func isSequencedOrdered(reliability byte) bool {
	return reliability == UnreliableSequenced || reliability == ReliableOrdered || reliability == ReliableSequenced || reliability == ReliableOrderedWithAckReceipt
}