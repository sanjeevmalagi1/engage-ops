package core

import "testing"

func TestAddressRoundTrip(t *testing.T) {
	addr := Address{Provider: "customerio", Type: "segment", Name: "active_users"}

	str := addr.String()
	if want := "customerio_segment.active_users"; str != want {
		t.Fatalf("String() = %q, want %q", str, want)
	}

	got, err := ParseAddress(str)
	if err != nil {
		t.Fatalf("ParseAddress() error = %v", err)
	}
	if got != addr {
		t.Fatalf("ParseAddress() = %+v, want %+v", got, addr)
	}
}

func TestParseAddressErrors(t *testing.T) {
	cases := []string{"no-dot-here", "noseparator.name"}
	for _, c := range cases {
		if _, err := ParseAddress(c); err == nil {
			t.Errorf("ParseAddress(%q) expected error, got nil", c)
		}
	}
}
