package launchpad

import "testing"

func TestGridNote(t *testing.T) {
	cases := []struct{ row, col, want int }{{0, 0, 81}, {0, 7, 88}, {7, 0, 11}, {7, 7, 18}}
	for _, tc := range cases {
		if got := GridNote(tc.row, tc.col); int(got) != tc.want {
			t.Fatalf("(%d,%d)=%d; want %d", tc.row, tc.col, got, tc.want)
		}
	}
}

func TestColor(t *testing.T) {
	if got := Color(3, 0); got != Red {
		t.Fatalf("red=%d", got)
	}
	if got := Color(0, 3); got != Green {
		t.Fatalf("green=%d", got)
	}
	if got := Color(3, 3); got != Amber {
		t.Fatalf("amber=%d", got)
	}
}
