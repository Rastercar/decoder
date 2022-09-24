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
