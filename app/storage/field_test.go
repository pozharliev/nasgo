package storage

import (
	"encoding/binary"
	"testing"
)

func TestStringFieldSerialize(t *testing.T) {
	field, err := NewStringField("testdata123", 12)

	if err != nil {
		t.Error(err)
	}

	serialized := field.Serialize()

	if len(serialized) != 16 {
		t.Errorf("serialized length should be 16, got %d", len(serialized))
	}

	size := int(binary.BigEndian.Uint32(serialized[0:4]))

	if size != 11 {
		t.Errorf("serialized size should be 11, got %d", size)
	}

	decoded := string(serialized[4 : size+4])

	if decoded != "testdata123" {
		t.Errorf("decoded value should be testdata123, got %s", decoded)
	}

	if serialized[len(serialized)-1] != 0 {
		t.Errorf("serialized value should be nil, got %v", serialized[len(serialized)-1])
	}
}

func TestStringField_Deserialize_ReturnsCorrectData(t *testing.T) {
	field, err := NewStringField("testdata123", 12)

	if err != nil {
		t.Error(err)
	}

	fieldDesc := &FieldDescriptor{Size: 12, DataType: StringType}
	tupleDesc := &TupleDescriptor{Fields: []FieldDescriptor{*fieldDesc}}

	serialized := field.Serialize()
	deserialized := Deserialize(serialized, *tupleDesc)

	if deserialized.Fields[0].(StringField).data != field.data {
		t.Error("serialize and deserialize should have the same data")
	}

}
