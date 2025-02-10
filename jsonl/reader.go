package jsonl

import (
	"bytes"
	"io"
	"os"

	"github.com/valyala/fastjson"
	"gosuda.org/deeplingua/internal/mmap"
)

type Reader struct {
	file     *os.File
	fileView []byte
	offset   int64
	size     int64
}

func NewReader(f *os.File) (*Reader, error) {
	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	view, err := mmap.Map(uintptr(f.Fd()), 0, int(stat.Size()), mmap.PROT_READ, mmap.MAP_SHARED)
	if err != nil {
		return nil, err
	}

	g := &Reader{
		file:     f,
		fileView: view,
		offset:   0,
		size:     stat.Size(),
	}

	return g, nil
}

func (g *Reader) Scan() (*Value, error) {
	if g.offset >= g.size || g.offset < 0 {
		return nil, io.EOF
	}

	// scan for newline
	idx := bytes.IndexByte(g.fileView[g.offset:], '\n')
	if idx == -1 {
		// try parse the last line
		v := &Value{}
		fv, err := fastjson.ParseBytes(g.fileView[g.offset:])
		if err != nil {
			if g.size-g.offset < 2 {
				return nil, io.EOF
			}
			return nil, err
		}
		v.Value = fv
		g.offset = g.size

		g.Close() // close the reader

		return v, nil
	}

	v := &Value{}
	fv, err := fastjson.ParseBytes(g.fileView[g.offset : g.offset+int64(idx)])
	if err != nil {
		return nil, err
	}
	v.Value = fv
	g.offset += int64(idx) + 1

	return v, nil
}

func (g *Reader) Close() error {
	if g.fileView == nil {
		return nil
	}

	err := mmap.UnMap(g.fileView)
	if err != nil {
		return err
	}
	g.fileView = nil
	g.offset = -1

	return nil
}
