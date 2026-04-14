package storage

type FieldDescriptor struct {
	DataType Type
	Size     int
}

type TupleDescriptor struct {
	Fields []FieldDescriptor
}

type Tuple struct {
	Fields []Field
	Desc   TupleDescriptor
}

func (tupleDesc *TupleDescriptor) Size() int {
	sum := 0

	for _, field := range tupleDesc.Fields {
		sum += field.Size
	}

	return sum
}

func (tuple Tuple) Serialize() []byte {
	var buf []byte

	for _, field := range tuple.Fields {
		serializedField := field.Serialize()

		buf = append(buf, serializedField...)
	}

	return buf
}

func Deserialize(buf []byte, tupleDesc TupleDescriptor) (Tuple, error) {
	pointer, newTuple := 0, Tuple{Desc: tupleDesc}

	for _, field := range tupleDesc.Fields {
		currentFieldByte := buf[pointer : pointer+field.Size]

		deserializedField, err := DeserializeField(currentFieldByte, field)

		if err != nil {
			return Tuple{}, err
		}

		newTuple.Fields = append(newTuple.Fields, deserializedField)

		pointer += field.Size
	}

	return newTuple, nil
}
