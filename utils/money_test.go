package utils

import "testing"

func TestParseAmount(t *testing.T) {
	cases := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{in: "1000", want: "1000.00"},
		{in: "1000.5", want: "1000.50"},
		{in: " 99.99 ", want: "99.99"},
		{in: "0.01", want: "0.01"},
		{in: "", wantErr: true},
		{in: "0", wantErr: true},
		{in: "-5", wantErr: true},
		{in: "1.005", wantErr: true},
		{in: "abc", wantErr: true},
		{in: "99999999999999", wantErr: true},
	}

	for _, tc := range cases {
		got, err := ParseAmount(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("ParseAmount(%q) = %q, want error", tc.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseAmount(%q) returned error: %v", tc.in, err)
			continue
		}
		if got != tc.want {
			t.Errorf("ParseAmount(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestParseSignedAmountAllowsZeroAndNegative(t *testing.T) {
	cases := map[string]string{
		"":        "0.00",
		"0":       "0.00",
		"-1500.5": "-1500.50",
		"2500":    "2500.00",
	}

	for in, want := range cases {
		got, err := ParseSignedAmount(in)
		if err != nil {
			t.Errorf("ParseSignedAmount(%q) returned error: %v", in, err)
			continue
		}
		if got != want {
			t.Errorf("ParseSignedAmount(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNegateAmount(t *testing.T) {
	cases := map[string]string{
		"100.00":  "-100.00",
		"-100.00": "100.00",
		"0.00":    "0.00",
		"0":       "0.00",
	}

	for in, want := range cases {
		if got := NegateAmount(in); got != want {
			t.Errorf("NegateAmount(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSubAmount(t *testing.T) {
	cases := []struct {
		a, b, want string
	}{
		{"1000.00", "250.50", "749.50"},
		{"100.00", "100.00", "0.00"},
		{"0.00", "25.25", "-25.25"},
		// A NUMERIC column never returns a blank, but a zero keeps the
		// helper total rather than letting a bad value poison the result.
		{"", "10.00", "-10.00"},
	}

	for _, tc := range cases {
		if got := SubAmount(tc.a, tc.b); got != tc.want {
			t.Errorf("SubAmount(%q, %q) = %q, want %q", tc.a, tc.b, got, tc.want)
		}
	}
}
