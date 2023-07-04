package interval

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		Input  string
		Output Interval
	}{
		{"1 day", Interval{Days: 1}},
		{"1 month", Interval{Months: 1}},
		{"1 second", Interval{Seconds: 1}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("subtest%d", i), func(t *testing.T) {
			i, err := Parse(tt.Input)
			if err != nil {
				t.Error(err)
				return
			}
			if i != tt.Output {
				t.Error(fmt.Errorf("Expected interval %#v: Got %#v", tt.Output, i))
				return
			}
		})
	}
}

func TestAddToDate(t *testing.T) {
	tests := []struct {
		Input    string
		Interval string
		Output   string
	}{
		{"10/21/1984", "1 day", "10/22/1984"},
		{"10/31/1984", "1 day", "11/1/1984"},
		{"10/31/1984", "1 month", "11/30/1984"},
		{"11/1/1984", "-1 day", "10/31/1984"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("subtest%d", i), func(t *testing.T) {
			input, err := ParseDate(tt.Input)
			if err != nil {
				t.Error(err)
				return
			}
			interval, err := Parse(tt.Interval)
			if err != nil {
				t.Error(err)
				return
			}
			output, err := ParseDate(tt.Output)
			if err != nil {
				t.Error(err)
				return
			}
			got := interval.AddToDate(input)

			if got != output {
				t.Error(fmt.Errorf("Expected %#v + %#v = %#v: Got %#v", tt.Input, tt.Interval, tt.Output, got))
				return
			}
		})
	}
}
