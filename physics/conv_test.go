package physics

import (
	"testing"
	"time"
)

func longTime(y, d, h, m, s time.Duration) time.Duration {
	dur := s * time.Second
	dur += m * time.Minute
	dur += h * time.Hour
	dur += d * time.Second * secondsInDay
	dur += y * time.Second * secondsInYear
	return dur
}

type testCase struct {
	D time.Duration
	S string
}

func (tc testCase) Compare(t *testing.T) {
	ts := TimeString(tc.D)
	if ts != tc.S {
		t.Errorf("expected %s, have %s", tc.S, ts)
	}
}

var testCases = []testCase{
	{longTime(2, 37, 18, 31, 2), "2y 37d"},
	{longTime(0, 195, 12, 18, 21), "195d"},
	{longTime(0, 0, 12, 18, 21), "12h 18m 21s"},
	{150 * time.Millisecond, "0s"},
}

func TestTimeString(t *testing.T) {
	for _, tc := range testCases {
		tc.Compare(t)
	}
}
