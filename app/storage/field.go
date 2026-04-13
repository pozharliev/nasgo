package storage

import (
	"encoding/binary"
)

type Type int

const (
	IntType Type = iota
	StringType
)

type Field interface {
	Serialize() []byte
	Size() int
	Type() Type
}

type IntField struct {
	data int
}

type StringField struct {
	data      string
	maxLength int
}

func (field IntField) Serialize() []byte {
	buf := make([]byte, binary.Size(field.data))

	binary.BigEndian.PutUint32(buf, uint32(field.data))
	return buf
}

func (field IntField) Size() int {
	return 4
}

func (field IntField) Type() Type {
	return IntType
}

func (field StringField) Serialize() []byte {
	buf := make([]byte, 4+field.maxLength)

	binary.BigEndian.PutUint32(buf[0:4], uint32(len(field.data)))
	copy(buf[4:], field.data)

	return buf
}

func (field StringField) Size() int {
	return field.maxLength + 4
}

func (field StringField) Type() Type {
	return StringType
}

func DeserializeField(data []byte, fieldDesc FieldDescriptor) Field {
	var result Field

	switch fieldDesc.DataType {
	case IntType:
		result = deserializeIntField(data)
	case StringType:
		result = deserializeStringField(data, fieldDesc)
	}

	return result
}

func deserializeIntField(data []byte) IntField {
	decoded := int(binary.BigEndian.Uint32(data))

	return IntField{data: decoded}
}

func deserializeStringField(data []byte, fieldDesc FieldDescriptor) StringField {
	size := int(binary.BigEndian.Uint32(data[0:4]))
	decoded := string(data[4 : size+4])

	return StringField{data: decoded, maxLength: fieldDesc.Size}
}
