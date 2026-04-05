package storage

type FieldDescriptor struct {
	dataType Type
	size     *int
}

type TupleDescriptor struct {
	fields []FieldDescriptor
}

type Tuple struct {
	fields []Field
}
