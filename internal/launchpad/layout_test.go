package launchpad

import "testing"

func TestGridNote(t *testing.T) {
	cases := []struct{ row, col, want int }{{0, 0, 0}, {0, 7, 7}, {7, 0, 112}, {7, 7, 119}}
	for _, tc := range cases {
		if got := GridNote(tc.row, tc.col); int(got) != tc.want {
			t.Fatalf("(%d,%d)=%d; want %d", tc.row, tc.col, got, tc.want)
		}
	}
}

func TestThreadBlock(t *testing.T) {
	cases := []struct{ thread, col, row int }{
		{0, 0, 0}, {1, 1, 0}, {2, 0, 4}, {3, 1, 4}, {4, 2, 0}, {7, 3, 4},
	}
	for _, tc := range cases {
		col, row := ThreadBlock(tc.thread)
		if col != tc.col || row != tc.row {
			t.Fatalf("thread %d: got (%d,%d), want (%d,%d)", tc.thread, col, row, tc.col, tc.row)
		}
	}
}

func TestSideButtonNote(t *testing.T) {
	if got := SideButtonNote(0); got != 8 {
		t.Fatalf("top=%d", got)
	}
	if got := SideButtonNote(7); got != 120 {
		t.Fatalf("bottom=%d", got)
	}
}
