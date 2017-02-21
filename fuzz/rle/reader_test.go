package rle

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	cases := []struct {
		in, want []byte
	}{
		{[]byte{}, []byte{}},
		{[]byte{0, 3}, []byte{0, 0, 0}},
		{[]byte{0, 12}, bytes.Repeat([]byte{0}, 10)},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("case %x", tc.in), func(t *testing.T) {
			r, err := NewReader(bytes.NewReader(tc.in))
			if err != nil && err != io.EOF {
				t.Fatalf("unexpected error: %s", err)
			}

			buf := make([]byte, 10)
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				t.Fatalf("unexpected read error: %s", err)
			}

			if n > len(buf) {
				t.Errorf("n = %d; len(buf) = %d", n, len(buf))
			}

			if n != len(tc.want) {
				t.Errorf("n = %d; want %d", n, len(tc.want))
			}

			if !arrayEquals(buf[:n], tc.want) {
				t.Errorf("got % x; want % x", buf, tc.want)
			}
		})
	}
}

func TestReadFromBuf(t *testing.T) {
	r := &Reader{
		buf: []byte{5, 4, 4, 1},
	}

	p := make([]byte, 3)
	n := r.readFromBuf(p)
	if n != 3 {
		t.Errorf("n = %d; want %d", n, 3)
	}

	want := []byte{5, 5, 5}
	if !arrayEquals(p, want) {
		t.Errorf("got %v; want %v", p, want)
	}

	n = r.readFromBuf(p)
	if n != 2 {
		t.Errorf("n = %d; want %d", n, 2)
	}

	want = []byte{5, 4}
	if !arrayEquals(p[:n], want) {
		t.Errorf("got %v; want %v", p[:n], want)
	}
}

func arrayEquals(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
