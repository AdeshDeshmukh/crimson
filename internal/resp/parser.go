package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = "string"
	ERROR   = "error"
	INTEGER = "integer"
	BULK    = "bulk"
	ARRAY   = "array"
	NULL    = "null"
)

type Value struct {
	Type  string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

type Parser struct {
	reader *bufio.Reader
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(reader),
	}
}

func (p *Parser) Parse() (Value, error) {
	prefix, err := p.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch prefix {
	case '+':
		return p.parseSimpleString()
	case '-':
		return p.parseError()
	case ':':
		return p.parseInteger()
	case '$':
		return p.parseBulkString()
	case '*':
		return p.parseArray()
	default:
		return Value{}, fmt.Errorf("unknown prefix: %s", string(prefix))
	}
}

func (p *Parser) readLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return line[:len(line)-2], nil
}

func (p *Parser) parseSimpleString() (Value, error) {
	line, err := p.readLine()
	if err != nil {
		return Value{}, err
	}
	return Value{Type: STRING, Str: line}, nil
}

func (p *Parser) parseError() (Value, error) {
	line, err := p.readLine()
	if err != nil {
		return Value{}, err
	}
	return Value{Type: ERROR, Str: line}, nil
}

func (p *Parser) parseInteger() (Value, error) {
	line, err := p.readLine()
	if err != nil {
		return Value{}, err
	}
	num, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, fmt.Errorf("invalid integer: %s", line)
	}
	return Value{Type: INTEGER, Num: num}, nil
}

func (p *Parser) parseBulkString() (Value, error) {
	line, err := p.readLine()
	if err != nil {
		return Value{}, err
	}

	length, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, fmt.Errorf("invalid bulk length: %s", line)
	}

	if length == -1 {
		return Value{Type: NULL}, nil
	}

	data := make([]byte, length)
	_, err = io.ReadFull(p.reader, data)
	if err != nil {
		return Value{}, err
	}

	p.reader.ReadByte()
	p.reader.ReadByte()

	return Value{Type: BULK, Bulk: string(data)}, nil
}

func (p *Parser) parseArray() (Value, error) {
	line, err := p.readLine()
	if err != nil {
		return Value{}, err
	}

	count, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, fmt.Errorf("invalid array length: %s", line)
	}

	if count == -1 {
		return Value{Type: NULL}, nil
	}

	array := make([]Value, count)
	for i := 0; i < count; i++ {
		val, err := p.Parse()
		if err != nil {
			return Value{}, err
		}
		array[i] = val
	}

	return Value{Type: ARRAY, Array: array}, nil
}