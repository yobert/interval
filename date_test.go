package interval

import (
	"fmt"
	"testing"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		Input  string
		Output Date
	}{
		{"10/21/1984", Date{Year: 1984, Month: 10, Day: 21}},
		{"1984-10-21", Date{Year: 1984, Month: 10, Day: 21}},
		{"October 21, 1984", Date{Year: 1984, Month: 10, Day: 21}},
		{"October 21 1984", Date{Year: 1984, Month: 10, Day: 21}},
		{"OCT 21, 1984", Date{Year: 1984, Month: 10, Day: 21}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("subtest%d", i), func(t *testing.T) {
			i, err := ParseDate(tt.Input)
			if err != nil {
				t.Error(err)
				return
			}
			if i != tt.Output {
				t.Error(fmt.Errorf("Expected date %#v: Got %#v", tt.Output, i))
				return
			}
		})
	}
}
