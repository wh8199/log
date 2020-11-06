package log

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	Pool sync.Pool
}

func NewBufferPool() *BufferPool {
	pool := &BufferPool{
		Pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}

	return pool
}

func (p *BufferPool) Get() *bytes.Buffer {
	buf := p.Pool.Get().(*bytes.Buffer)
	buf.Reset()

	return buf
}

func (p *BufferPool) Put(buf *bytes.Buffer) {
	p.Pool.Put(buf)
}
