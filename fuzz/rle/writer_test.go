package rle

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestWriter(t *testing.T) {
	t.Skip("not implemented")

	cases := []struct {
		in, want []byte
	}{
		{[]byte{}, []byte{}},
		{[]byte{0, 0, 0}, []byte{0, 3}},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("case %x", tc.in), func(t *testing.T) {
			buf := new(bytes.Buffer)
			w, err := NewWriter(buf)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			_, err = w.Write(tc.in)
			if err != nil {
				t.Errorf("unexpected write error: %s", err)
			}

			got := buf.Bytes()
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %x; want %x", got, tc.want)
			}
		})
	}
}
