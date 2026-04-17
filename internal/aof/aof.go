package aof

import (
	"bufio"
	"io"
	"os"
	"sync"

	"github.com/AdeshDeshmukh/crimson/internal/resp"
)

type AOF struct {
	file   *os.File
	writer *bufio.Writer
	mu     sync.Mutex
}

func New(path string) (*AOF, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &AOF{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (a *AOF) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.writer.Flush()
	return a.file.Close()
}

func (a *AOF) Write(value resp.Value) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, err := a.writer.Write(value.Marshal())
	if err != nil {
		return err
	}

	return a.writer.Flush()
}

func (a *AOF) Load(fn func(value resp.Value)) error {
	a.file.Seek(0, io.SeekStart)

	reader := resp.NewParser(a.file)

	for {
		value, err := reader.Parse()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		fn(value)
	}

	return nil
}