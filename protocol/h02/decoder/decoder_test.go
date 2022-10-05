package decoder

import (
	"context"
	"math"
	"reciever-ms/config"
	"testing"
)

const float64EqualityThreshold = 1e-7

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestDecode(t *testing.T) {
	dec := New(&config.Config{})

	t.Run("fails if packet does not start with *HQ, or ends with #", func(t *testing.T) {
		_, err := dec.Decode(context.TODO(), []byte("*HQ,abc"))
		if err == nil {
			t.Fatal()
		}

		_, err = dec.Decode(context.TODO(), []byte("abc#"))
		if err == nil {
			t.Fatal()
		}
	})

	t.Run("fails if packet csv have lenght less than 3", func(t *testing.T) {
		_, err := dec.Decode(context.TODO(), []byte("*HQ,1,2#"))
		if err == nil {
			t.Fatal()
		}
	})

	t.Run("fails message type is unknown", func(t *testing.T) {
		_, err := dec.Decode(context.TODO(), []byte("*HQ,867232051148352,invalid,044639,A,2027.93290,S,05434.94389,W,0.00,0,110722,FFFFFBFF#"))
		if err == nil {
			t.Fatal()
		}
	})

	t.Run("fails if location marks itself as invalid", func(t *testing.T) {
		_, err := dec.Decode(context.TODO(), []byte("*HQ,867232051148352,V1,044639,B,2027.93290,S,05434.94389,W,0.00,20,110722,FFFFFBFF#"))
		if err == nil {
			t.Fatal()
		}
	})

	t.Run("successfully decodes a location message", func(t *testing.T) {
		res, err := dec.Decode(context.TODO(), []byte("*HQ,867232051148352,V1,044639,A,2027.93290,S,05434.94389,W,10.00,20,110722,FFFFFBFF#"))
		if err != nil {
			t.Fatalf("unexpected err %v", err)
		}

		if res.Evt.Type != "h02:LocationMsg" {
			t.Fatal("invalid msg type")
		}

		msg := res.Evt.Data.(LocationMsg)

		if msg.Imei != "867232051148352" {
			t.Fatal("invalid imei")
		}

		if !almostEqual(msg.Lat, -20.4655483333) {
			t.Fatal("invalid lat")
		}

		if !almostEqual(msg.Lng, -54.5823981666) {
			t.Fatal("invalid lng")
		}

		if msg.Direction != 20 {
			t.Fatal("invalid lng")
		}

		if msg.Speed != 18 {
			t.Fatal("invalid speed")
		}

		s := msg.Timestamp

		if s.Year() != 2022 || s.Month() != 7 || s.Day() != 11 || s.Hour() != 4 || s.Minute() != 46 || s.Second() != 39 {
			t.Fatal("invalid date")
		}
	})
}
