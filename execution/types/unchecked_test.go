package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	system "github.com/LimeChain/gosemble/frame/system/dispatchables"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var remarkCall = Call{
	CallIndex: types.CallIndex{
		ModuleIndex:   0,
		FunctionIndex: 0,
	},
	function: system.FnRemark{},
	Args:     sc.NewVaryingData(sc.Sequence[sc.U8]{}),
}

func Test_EncodeUncheckedExtrinsic_Unsigned(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       UncheckedExtrinsic
		expectation []byte
	}{
		{
			label:       "Encode(UnsignedUncheckedExtrinsic)",
			input:       NewUnsignedUncheckedExtrinsic(remarkCall),
			expectation: []byte{0x10, 0x4, 0x0, 0x0, 0x0},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			assert.Equal(t, testExample.expectation, buffer.Bytes())
			assert.Equal(t, sc.Bool(false), testExample.input.IsSigned())
		})
	}
}

func Test_DecodeUncheckedExtrinsic_Unsigned(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation UncheckedExtrinsic
	}{
		{
			label:       "Decode(UnsignedUncheckedExtrinsic)",
			input:       []byte{0x10, 0x4, 0x0, 0x0, 0x0},
			expectation: NewUnsignedUncheckedExtrinsic(remarkCall),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := DecodeUncheckedExtrinsic(buffer)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}

func Test_EncodeUncheckedExtrinsic_Signed(t *testing.T) {
	signer := types.NewMultiAddressId(types.AccountId{types.NewAddress32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1)})
	signature := types.NewMultiSignatureEd25519(types.NewEd25519(sc.FixedSequence[sc.U8]{0x00, 0x62, 0x37, 0x61, 0x33, 0x63, 0x31, 0x32, 0x64, 0x63, 0x30, 0x63, 0x38, 0x63, 0x37, 0x34, 0x38, 0x61, 0x62, 0x30, 0x37, 0x35, 0x32, 0x35, 0x62, 0x37, 0x30, 0x31, 0x31, 0x32, 0x32, 0x62, 0x38, 0x38, 0x62, 0x64, 0x37, 0x38, 0x66, 0x36, 0x30, 0x30, 0x63, 0x37, 0x36, 0x33, 0x34, 0x32, 0x64, 0x32, 0x37, 0x66, 0x32, 0x35, 0x65, 0x35, 0x66, 0x39, 0x32, 0x34, 0x34, 0x34, 0x63, 0x64}...))
	extra := types.SignedExtra{
		Era:   types.NewImmortalEra(),
		Nonce: 0,
		Fee:   0,
	}

	var testExamples = []struct {
		label       string
		input       UncheckedExtrinsic
		expectation []byte
	}{
		{
			label:       "Encode(SignedUncheckedExtrinsic)",
			input:       NewSignedUncheckedExtrinsic(remarkCall, signer, signature, extra),
			expectation: []byte{0xa5, 0x1, 0x84, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x62, 0x37, 0x61, 0x33, 0x63, 0x31, 0x32, 0x64, 0x63, 0x30, 0x63, 0x38, 0x63, 0x37, 0x34, 0x38, 0x61, 0x62, 0x30, 0x37, 0x35, 0x32, 0x35, 0x62, 0x37, 0x30, 0x31, 0x31, 0x32, 0x32, 0x62, 0x38, 0x38, 0x62, 0x64, 0x37, 0x38, 0x66, 0x36, 0x30, 0x30, 0x63, 0x37, 0x36, 0x33, 0x34, 0x32, 0x64, 0x32, 0x37, 0x66, 0x32, 0x35, 0x65, 0x35, 0x66, 0x39, 0x32, 0x34, 0x34, 0x34, 0x63, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			assert.Equal(t, testExample.expectation, buffer.Bytes())
			assert.Equal(t, sc.Bool(true), testExample.input.IsSigned())
			assert.Equal(t, sc.U8(132), testExample.input.Version)
		})
	}
}

func Test_DecodeUncheckedExtrinsic_Signed(t *testing.T) {
	signer := types.NewMultiAddressId(types.AccountId{types.NewAddress32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1)})
	signature := types.NewMultiSignatureEd25519(types.NewEd25519(sc.FixedSequence[sc.U8]{0x0, 0x62, 0x37, 0x61, 0x33, 0x63, 0x31, 0x32, 0x64, 0x63, 0x30, 0x63, 0x38, 0x63, 0x37, 0x34, 0x38, 0x61, 0x62, 0x30, 0x37, 0x35, 0x32, 0x35, 0x62, 0x37, 0x30, 0x31, 0x31, 0x32, 0x32, 0x62, 0x38, 0x38, 0x62, 0x64, 0x37, 0x38, 0x66, 0x36, 0x30, 0x30, 0x63, 0x37, 0x36, 0x33, 0x34, 0x32, 0x64, 0x32, 0x37, 0x66, 0x32, 0x35, 0x65, 0x35, 0x66, 0x39, 0x32, 0x34, 0x34, 0x34, 0x63, 0x64}...))
	extra := types.SignedExtra{
		Era:   types.NewImmortalEra(),
		Nonce: 0,
		Fee:   0,
	}

	var testExamples = []struct {
		label       string
		input       []byte
		expectation UncheckedExtrinsic
	}{
		{
			label:       "Decode(SignedUncheckedExtrinsic)",
			input:       []byte{0xa5, 0x1, 0x84, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x62, 0x37, 0x61, 0x33, 0x63, 0x31, 0x32, 0x64, 0x63, 0x30, 0x63, 0x38, 0x63, 0x37, 0x34, 0x38, 0x61, 0x62, 0x30, 0x37, 0x35, 0x32, 0x35, 0x62, 0x37, 0x30, 0x31, 0x31, 0x32, 0x32, 0x62, 0x38, 0x38, 0x62, 0x64, 0x37, 0x38, 0x66, 0x36, 0x30, 0x30, 0x63, 0x37, 0x36, 0x33, 0x34, 0x32, 0x64, 0x32, 0x37, 0x66, 0x32, 0x35, 0x65, 0x35, 0x66, 0x39, 0x32, 0x34, 0x34, 0x34, 0x63, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			expectation: NewSignedUncheckedExtrinsic(remarkCall, signer, signature, extra),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := DecodeUncheckedExtrinsic(buffer)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}

func Test_DecodeUncheckedExtrinsic_Panics_InvalidLength(t *testing.T) {
	input := []byte{0xa9, 0x1, 0x84, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x62, 0x37, 0x61, 0x33, 0x63, 0x31, 0x32, 0x64, 0x63, 0x30, 0x63, 0x38, 0x63, 0x37, 0x34, 0x38, 0x61, 0x62, 0x30, 0x37, 0x35, 0x32, 0x35, 0x62, 0x37, 0x30, 0x31, 0x31, 0x32, 0x32, 0x62, 0x38, 0x38, 0x62, 0x64, 0x37, 0x38, 0x66, 0x36, 0x30, 0x30, 0x63, 0x37, 0x36, 0x33, 0x34, 0x32, 0x64, 0x32, 0x37, 0x66, 0x32, 0x35, 0x65, 0x35, 0x66, 0x39, 0x32, 0x34, 0x34, 0x34, 0x63, 0x64, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	buffer := &bytes.Buffer{}
	buffer.Write(input)

	assert.PanicsWithValue(
		t,
		"invalid length prefix",
		func() {
			DecodeUncheckedExtrinsic(buffer)
		},
	)
}
