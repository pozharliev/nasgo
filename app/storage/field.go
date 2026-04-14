package storage

import (
	"encoding/binary"
	"fmt"
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

func NewIntField(data int) (IntField, error) {
	return IntField{data: data}, nil
}

func NewStringField(data string, maxLength int) (StringField, error) {
	if len(data) > maxLength {
		return StringField{},
			fmt.Errorf("string data length %d exceeds maxLength %d", len(data), maxLength)
	}

	return StringField{data: data, maxLength: maxLength}, nil
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

func DeserializeField(data []byte, fieldDesc FieldDescriptor) (Field, error) {
	var result Field
	var err error

	switch fieldDesc.DataType {
	case IntType:
		result, err = deserializeIntField(data)
	case StringType:
		result, err = deserializeStringField(data, fieldDesc)
	}

	return result, err
}

func deserializeIntField(data []byte) (IntField, error) {
	decoded := int(binary.BigEndian.Uint32(data))

	field, err := NewIntField(decoded)
	return field, err

}

func deserializeStringField(data []byte, fieldDesc FieldDescriptor) (StringField, error) {
	size := binary.BigEndian.Uint32(data[0:4])
	decoded := string(data[4 : size+4])

	field, err := NewStringField(decoded, fieldDesc.Size)

	if err != nil {
		return StringField{}, err
	}

	return field, nil
}
