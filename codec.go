package noise

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"reflect"
	"sync"
)

const (
	OpcodeSize   = crc32.Size
	maxMsgNumber = 256
)

// Serializable attributes whether or not a type has a byte representation that it may be serialized into.
type Serializable interface {
	// Marshal converts this type into it's byte representation as a slice.
	Marshal() []byte
}

type codec struct {
	sync.RWMutex

	ser map[reflect.Type]uint32
	de  map[uint32]reflect.Value
}

func newCodec() *codec {
	return &codec{
		ser: make(map[reflect.Type]uint32, maxMsgNumber),
		de:  make(map[uint32]reflect.Value, maxMsgNumber),
	}
}

func (c *codec) register(ser Serializable, de interface{}) uint32 {
	c.Lock()
	defer c.Unlock()

	t := reflect.TypeOf(ser)
	d := reflect.ValueOf(de)

	if opcode, registered := c.ser[t]; registered {
		panic(fmt.Errorf("attempted to register type %+v which is already registered under opcode %d", t, opcode))
	}

	in := []reflect.Type{reflect.TypeOf(([]byte)(nil))}
	var out []reflect.Type

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		out = []reflect.Type{reflect.New(t).Type(), reflect.TypeOf((*error)(nil)).Elem()}
	} else {
		out = []reflect.Type{t, reflect.TypeOf((*error)(nil)).Elem()}
	}

	expected := reflect.FuncOf(in, out, false)

	if d.Type() != expected {
		panic(fmt.Errorf("provided decoder for message type %+v is %s, but expected %s", t, d, expected))
	}

	newOpcode := crc32.ChecksumIEEE([]byte(t.Name()))

	if _, registered := c.de[newOpcode]; registered {
		panic(fmt.Errorf("attempted to register type %+v whose opcode %d collides with opcode of already registered type", t, newOpcode))
	}

	c.ser[t] = newOpcode
	c.de[newOpcode] = d

	return newOpcode
}

func (c *codec) encode(msg Serializable) ([]byte, error) {
	c.RLock()
	defer c.RUnlock()

	t := reflect.TypeOf(msg)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	opcode, registered := c.ser[t]
	if !registered {
		return nil, fmt.Errorf("opcode not registered for message type %+v", t)
	}

	buf := make([]byte, OpcodeSize)
	binary.BigEndian.PutUint32(buf[:OpcodeSize], opcode)

	return append(buf, msg.Marshal()...), nil
}

func (c *codec) decode(data []byte) (Serializable, error) {
	if len(data) < OpcodeSize {
		return nil, io.ErrUnexpectedEOF
	}

	opcode := binary.BigEndian.Uint32(data[:OpcodeSize])
	data = data[OpcodeSize:]

	c.RLock()
	defer c.RUnlock()

	decoder, registered := c.de[opcode]
	if !registered {
		return nil, fmt.Errorf("opcode %d is not registered", opcode)
	}

	results := decoder.Call([]reflect.Value{reflect.ValueOf(data)})

	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}

	return results[0].Interface().(Serializable), nil
}
