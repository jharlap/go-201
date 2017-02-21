package tests

import (
	"strconv"
	"testing"
)

func TestSquareOld(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{1, 1},
		{2, 4},
		{8, 64},
		{0, 0},
		{-1, 1},
	}

	for _, tc := range cases {
		got := Square(tc.in)
		if got != tc.want {
			t.Errorf("Square(%d) = %d; want %d", tc.in, got, tc.want)
		}
	}
}

func TestSquareSubtests(t *testing.T) {
	cases := []struct {
		in, want int
	}{
		{1, 1},
		{2, 4},
		{8, 64},
		{0, 0},
		{-1, 1},
	}

	for _, tc := range cases {
		t.Run(strconv.Itoa(tc.in), func(t *testing.T) {
			got := Square(tc.in)
			if got != tc.want {
				t.Errorf("got %d; want %d", got, tc.want)
			}
		})
	}
}
