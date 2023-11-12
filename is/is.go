package is

import "testing"

func Equal[C comparable](t *testing.T, want, got C) {
	t.Helper()
	if want != got {
		t.Fatalf("Values not equal\nwant=%+v\ngot=%+v", want, got)
	}
}
