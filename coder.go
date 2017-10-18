// TODO:
// bench test with std.binary and gob
// function test
// field tag parse
// read buffer not enough, need return errr, not panic

package binary

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrNotEnoughSpace = errors.New("not enough space")
)

type coder struct {
	buff   []byte
	pos    int
	endian Endian
}

// Buffer returns the byte slice that has been encoding/decoding.
func (this *coder) Buffer() []byte {
	return this.buff[:this.pos]
}

// Len returns unmber of bytes that has been encoding/decoding.
func (this *coder) Len() int {
	return this.pos
}

// Cap returns number total bytes of this coder buffer.
func (this *coder) Cap() int {
	return len(this.buff)
}

// Skip ignore size bytes for encoding/decoding.
// If with errors, it will return -1
func (this *coder) Skip(size int) int {
	newPos := this.pos + size
	if size >= 0 && newPos <= this.Cap() {
		for i, b := int(size-1), this.buff[this.pos:newPos]; i >= 0; i-- { //zero skiped bytes
			b[i] = 0
		}
		this.pos = newPos
		return size
	}
	return -1
}

// Reset move the read/write pointer to the beginning of buffer.
func (this *coder) Reset() {
	for i := int(this.pos - 1); i >= 0; i-- { //zero encoded bytes
		this.buff[i] = 0
	}
	this.pos = 0
}

// reserve returns next size bytes for encoding/decoding.
func (this *coder) reserve(size int) []byte {
	newPos := this.pos + size
	if newPos > this.Cap() {
		panic(fmt.Errorf("binary.Coder:buffer overflow pos=%d cap=%d require=%d, not enough space!", this.pos, this.Cap(), size))
	}
	if size > 0 && newPos <= this.Cap() {
		b := this.buff[this.pos:newPos]
		this.pos = newPos
		return b
	}
	return nil
}

type BytesReader []byte

func (p *BytesReader) Read(data []byte) (n int, err error) {
	n = copy(data, *p)
	if n == len(*p) {
		err = io.EOF
	}
	*p = (*p)[n:]
	return
}

type BytesWriter []byte

func (p *BytesWriter) Write(data []byte) (n int, err error) {
	n = copy(*p, data)
	if n < len(data) {
		err = ErrNotEnoughSpace
	}
	*p = (*p)[n:]
	return
}
