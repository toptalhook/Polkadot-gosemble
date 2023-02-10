package types

import (
	"bytes"
	"errors"
	"fmt"
	sc "github.com/LimeChain/goscale"
)

type InherentData struct {
	Data map[[8]byte]sc.Sequence[sc.U8]
}

func NewInherentData() *InherentData {
	return &InherentData{
		Data: make(map[[8]byte]sc.Sequence[sc.U8]),
	}
}

func (id *InherentData) Bytes() []byte {
	buf := &bytes.Buffer{}

	sc.Compact(sc.NewU128FromUint64(uint64(len(id.Data)))).Encode(buf)

	for k, v := range id.Data {
		buf.Write(k[:])
		buf.Write(v.Bytes())
	}

	return buf.Bytes()
}

func (id *InherentData) Put(key [8]byte, value sc.Encodable) error {
	if id.Data[key] != nil {
		return errors.New(fmt.Sprintf("InherentDataExists - [%v]", key))
	}

	id.Data[key] = sc.BytesToSequenceU8(value.Bytes())

	return nil
}

func DecodeInherentData(buffer *bytes.Buffer) (*InherentData, error) {
	result := NewInherentData()
	length := sc.DecodeCompact(buffer).ToBigInt().Int64()

	for i := 0; i < int(length); i++ {
		key := [8]byte{}
		len, err := buffer.Read(key[:])
		if err != nil {
			return nil, err
		}
		if len != 8 {
			return nil, errors.New("invalid length")
		}
		value := sc.DecodeSequence[sc.U8](buffer)

		result.Data[key] = value
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}