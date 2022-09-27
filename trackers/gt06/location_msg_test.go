package gt06

import "testing"

func TestDecodeDateTime(t *testing.T) {
	t.Run("succeds on valid date as bytes", func(t *testing.T) {
		var failIf = func(cond bool) {
			if cond {
				t.Fail()
			}
		}

		time, err := decodeDateTime([]byte{0x0a, 0x03, 0x17, 0x0f, 0x32, 0x17})

		failIf(err != nil)
		failIf(time.Year() != 2010)
		failIf(time.Month() != 3)
		failIf(time.Day() != 23)
		failIf(time.Hour() != 15)
		failIf(time.Minute() != 50)
		failIf(time.Second() != 23)
	})

	t.Run("fails on insuficient bytes", func(t *testing.T) {
		_, err := decodeDateTime([]byte{0x03, 0x17, 0x0f, 0x32, 0x17})
		if err == nil {
			t.Fatal("should have failed on insuficient bytes")
		}
	})

	t.Run("fails on bad data bytes", func(t *testing.T) {
		var month16 uint8 = 0x10

		_, err := decodeDateTime([]byte{0x0a, month16, 0x17, 0x0f, 0x32, 0x17})
		if err == nil {
			t.Fatal("should have failed on month > 12")
		}
	})
}

func TestDecodeSateliteCnt(t *testing.T) {
	type io struct {
		i byte
		o int
	}

	tests := []io{
		{
			i: 0xCB,
			o: 11,
		},
		{
			i: 0x09,
			o: 9,
		},
		{
			i: 0x00,
			o: 0,
		},
		{
			i: 0xFF,
			o: 15,
		},
	}

	for i := 0; i < 16; i++ {
		tests = append(tests, io{i: byte(i), o: i})
	}

	for i, test := range tests {
		r, err := decodeSateliteCnt(test.i)
		if err != nil {
			t.Fatalf("test failed for input: %d", i)
		}

		if r != test.o {
			t.Fatalf("test failed for input: %d, wanted %d got %d", i, test.o, r)
		}
	}
}

func TestDecodeCoord(t *testing.T) {
	type io struct {
		i []byte
		o float64
	}

	tests := []io{
		{
			i: []byte{0x02, 0x6B, 0x3F, 0x3E},
			o: 22.546096660000,
		},
	}

	for i, test := range tests {
		r := decodeCoord(test.i)
		if !almostEqual(r, test.o) {
			t.Fatalf("test failed for input: %d, wanted near %f got %f", i, test.o, r)
		}
	}
}
