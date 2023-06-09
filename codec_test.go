package noise

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

type test2 struct {
	data []byte
}

func (t *test2) Marshal() []byte {
	return t.data
}

func unmarshalTest2(data []byte) (*test2, error) {
	return &test2{data: data}, nil
}

type test struct {
	data []byte
}

func (t *test) Marshal() []byte {
	return t.data
}

func unmarshalTest(data []byte) (*test, error) {
	return &test{data: data}, nil
}

func TestCodecRegisterEncodeDecode(t *testing.T) {
	t.Parallel()

	codec := newCodec()

	opcode := codec.register(&test{}, unmarshalTest)

	msg := &test{data: []byte("hello world")}

	expected := make([]byte, OpcodeSize+len(msg.data))
	binary.BigEndian.PutUint32(expected[:OpcodeSize], opcode)
	copy(expected[OpcodeSize:], msg.data)

	data, err := codec.encode(msg)
	assert.NoError(t, err)

	assert.EqualValues(t, opcode, binary.BigEndian.Uint32(data[:OpcodeSize]))
	assert.EqualValues(t, expected, data)

	obj, err := codec.decode(data)
	assert.NoError(t, err)
	assert.IsType(t, obj, &test{})

	// Failure cases.

	data[0] = 99
	_, err = codec.decode(data)
	assert.Error(t, err)

	_, err = codec.encode(&test2{data: []byte("should not be encodable")})
	assert.Error(t, err)

}

func TestPanicIfDuplicateMessagesRegistered(t *testing.T) {
	t.Parallel()

	codec := newCodec()

	assert.Panics(t, func() {
		codec.register(&test{}, unmarshalTest)
		codec.register(&test2{}, unmarshalTest2)
		codec.register(&test{}, unmarshalTest)
	})
}
