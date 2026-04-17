package resp

import (
	"fmt"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (w *Writer) Write(v Value) error {
	bytes := v.Marshal()
	_, err := w.writer.Write(bytes)
	return err
}

func (v Value) Marshal() []byte {
	switch v.Type {
	case STRING:
		return marshalSimpleString(v.Str)
	case ERROR:
		return marshalError(v.Str)
	case INTEGER:
		return marshalInteger(v.Num)
	case BULK:
		return marshalBulkString(v.Bulk)
	case ARRAY:
		return marshalArray(v.Array)
	case NULL:
		return marshalNull()
	default:
		return []byte{}
	}
}

func marshalSimpleString(s string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", s))
}

func marshalError(s string) []byte {
	return []byte(fmt.Sprintf("-%s\r\n", s))
}

func marshalInteger(n int) []byte {
	return []byte(fmt.Sprintf(":%d\r\n", n))
}

func marshalBulkString(s string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
}

func marshalNull() []byte {
	return []byte("$-1\r\n")
}

func marshalArray(array []Value) []byte {
	result := []byte(fmt.Sprintf("*%d\r\n", len(array)))
	for _, v := range array {
		result = append(result, v.Marshal()...)
	}
	return result
}