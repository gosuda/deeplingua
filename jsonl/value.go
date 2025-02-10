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
	*fastjson.Value
}
