package rle

import "io"

func NewWriter(w io.Writer) (*Writer, error) {
	return &Writer{w: w}, nil
}

type Writer struct {
	w io.Writer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return 0, nil
}
