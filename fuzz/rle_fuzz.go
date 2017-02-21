package fuzz

import (
	"bytes"
	"io"

	"github.com/jharlap/go-201/fuzz/rle"
)

func Fuzz(data []byte) int {
	r, err := rle.NewReader(bytes.NewReader(data))
	if err != nil {
		// error handling worked
		return 0
	}

	buf := make([]byte, 64<<10)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		return 0
	}

	// empty in, empty out
	if len(data) == 0 && n == 0 {
		return 0
	}

	// output size matches expectation
	if len(data) > 0 && n <= len(buf) && n == expectedOutputSize(data) {
		return 0
	}

	return 1
}

func expectedOutputSize(data []byte) int {
	var r int
	for i := 1; i < len(data); i += 2 {
		r += int(data[i])
	}
	return r
}
