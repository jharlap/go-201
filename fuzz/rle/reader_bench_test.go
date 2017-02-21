package rle

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

func BenchmarkReader(b *testing.B) {
	br := bytes.NewReader([]byte{87, 1, 65, 255, 32, 1, 87, 1, 65, 255, 32, 1, 87, 1, 65, 255})
	for n := 0; n < b.N; n++ {
		br.Seek(0, io.SeekStart)
		r, _ := NewReader(br)
		_, err := ioutil.ReadAll(r)
		if err != nil {
			b.Fatalf("unexpected error: %s", err)
		}
	}
}
