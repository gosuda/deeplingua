package jsonl

import (
	"sync"

	"github.com/valyala/fastjson"
)

var valuePool = &sync.Pool{
	New: func() interface{} {
		return &Value{}
	},
}

type Value struct {
	r *Reader
	p *fastjson.Parser
	*fastjson.Value
}

func (v *Value) Close() error {
	if v == nil {
		return nil
	}

	if v.r != nil && v.p != nil {
		v.r.pool.Put(v.p)
		v.r = nil
		v.p = nil
		v.Value = nil
		valuePool.Put(v)
	}
	return nil
}
