package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// The source of the transaction.
//
// Depending on the source we might apply different validation schemes.
// For instance we can disallow specific kinds of transactions if they were not produced
// by our local node (for instance off-chain workers).
const (
	// Transaction is already included in block.
	//
	// This means that we can't really tell where the transaction is coming from,
	// since it's already in the received block. Note that the custom validation logic
	// using either `Local` or `External` should most likely just allow `InBlock`
	// transactions as well.
	InBlock sc.U8 = iota

	// Transaction is coming from a local source.
	//
	// This means that the transaction was produced internally by the node
	// (for instance an Off-Chain Worker, or an Off-Chain Call), as opposed
	// to being received over the network.
	Local

	// Transaction has been received externally.
	//
	// This means the transaction has been received from (usually) "untrusted" source,
	// for instance received over the network or RPC.
	External
)

type TransactionSource sc.VaryingData

func NewTransactionSource(value sc.U8) TransactionSource {
	return TransactionSource(sc.NewVaryingData(value))
}

func (ts TransactionSource) Encode(buffer *bytes.Buffer) {
	// TODO:
}

func DecodeTransactionSource(buffer *bytes.Buffer) TransactionSource {
	// TODO:
	return TransactionSource{}
}

func (ts TransactionSource) Bytes() []byte {
	return sc.EncodedBytes(ts)
}
