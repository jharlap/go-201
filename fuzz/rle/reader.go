package rle

import (
	"errors"
	"io"
)

// Reader uncompresses RLE data from a Reader stream
type Reader struct {
	r   io.Reader
	buf []byte
}

func NewReader(r io.Reader) (*Reader, error) {
	return &Reader{r: r}, nil
}

var errTruncatedInput = errors.New("truncated input")

func (r *Reader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	n := 0
	for n < len(p) {
		if len(r.buf) == 0 {
			err := r.fillBuf()
			if err != nil {
				if err == io.EOF {
					return n, err
				}
				return 0, err
			}
		}

		if len(r.buf) > 0 {
			n += r.readFromBuf(p[n:])
		} else {
			// no more data available
			return n, nil
		}
	}

	return n, nil
}

func (r *Reader) fillBuf() error {
	buf := make([]byte, 1024)
	n, err := r.r.Read(buf)
	if err != nil {
		return err
	}

	r.buf = buf[:n]
	return nil
}

func (r *Reader) readFromBuf(p []byte) int {
	n := 0
	for n < len(p) && len(r.buf) > 1 {
		i := byte(0)
		for i < r.buf[1] {
			p[n] = r.buf[0]
			n++
			i++

			if n == len(p) {
				r.buf[1] -= i
				return n
			}
		}
		r.buf = r.buf[2:]
	}
	return n
}
