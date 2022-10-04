package strings

import "testing"

func TestBytesAsLiteralString(t *testing.T) {
	type inOut struct {
		In  []byte
		Out string
	}

	tests := []inOut{
		{
			In:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			Out: "0102030405060708090A",
		},
		{
			In:  []byte{},
			Out: "",
		},
		{
			In:  []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa},
			Out: "0102030405060708090A",
		},
		{
			In:  []byte{0xf, 0xa},
			Out: "0F0A",
		},
	}

	for _, test := range tests {
		result := BytesAsLiteralString(test.In)
		if test.Out != result {
			t.Fatalf("expected: %s - got: %s", test.Out, result)
		}
	}
}

func TestGetStringsBetween(t *testing.T) {
	s := "abc"
	e := "xyz"

	t.Run("finds strings in between", func(t *testing.T) {
		x := "something before" + s + "find me !" + e + "something after"

		if GetStringInBetween(x, s, e) != "find me !" {
			t.Fail()
		}
	})

	t.Run("returns empty string if start is not found", func(t *testing.T) {
		x := "i dont contain start but have end" + e

		if GetStringInBetween(x, s, e) != "" {
			t.Fail()
		}
	})

	t.Run("returns empty string if end is not found", func(t *testing.T) {
		x := s + "i dont contain end but have start"

		if GetStringInBetween(x, s, e) != "" {
			t.Fail()
		}
	})

	t.Run("returns empty string if neither is found", func(t *testing.T) {
		x := "i dont contain end nor start"

		if GetStringInBetween(x, s, e) != "" {
			t.Fail()
		}
	})

	t.Run("returns the first string in between on multiple occourances", func(t *testing.T) {
		x := s + "occourance_1" + e + "meh" + s + "occourance_2" + e

		if GetStringInBetween(x, s, e) != "occourance_1" {
			t.Fail()
		}
	})

	t.Run("succeds even if start and end are equal", func(t *testing.T) {
		x := s + "meh" + s

		if GetStringInBetween(x, s, s) != "meh" {
			t.Fail()
		}
	})
}
