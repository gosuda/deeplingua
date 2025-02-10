package jsonl

import (
	"os"
)

const (
	bufferFlushSize = 4 * 1024 * 1024 // 4MB
)

type Writer struct {
	file   *os.File
	buffer []byte
}

func NewWriter(f *os.File) (*Writer, error) {
	return &Writer{
		file: f,
	}, nil
}

func (g *Writer) Write(v *Value) error {
	if v == nil || v.Value == nil {
		return nil
	}

	g.buffer = v.Value.MarshalTo(g.buffer)
	g.buffer = append(g.buffer, '\n')
	if len(g.buffer) >= bufferFlushSize {
		return g.Flush()
	}

	err := g.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (g *Writer) Flush() error {
	if len(g.buffer) == 0 {
		return nil
	}

	_, err := g.file.Write(g.buffer)
	g.buffer = g.buffer[:0]
	return err
}

func (g *Writer) Close() error {
	err := g.Flush()
	if err != nil {
		return err
	}
	g.file = nil
	g.buffer = nil
	return nil
}
