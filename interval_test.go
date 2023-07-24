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
		{"666 businessdays", Interval{WorkDays: 666}},
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
		{"10/22/1984", "-1 day", "10/21/1984"},
		{"10/31/1984", "1 day", "11/1/1984"},
		{"11/1/1984", "-1 day", "10/31/1984"},
		{"10/31/1984", "1 month", "11/30/1984"},
		{"11/30/1984", "-1 month", "10/30/1984"},
		{"12/31/1984", "-1 month", "11/30/1984"},
		{"1/31/1984", "1 month", "2/29/1984"},
		{"2/29/1984", "-1 month", "1/29/1984"},
		{"3/31/1984", "-1 month", "2/29/1984"},
		{"1/31/1985", "1 month", "2/28/1985"},
		{"2/28/1985", "-1 month", "1/28/1985"},
		{"3/31/1985", "-1 month", "2/28/1985"},
		{"11/1/1984", "-1 day", "10/31/1984"},
		{"3/31/1984", "-1 month", "2/29/1984"},

		// business day/holiday logic
		{"12/30/2021", "1 businessday", "1/3/2022"}, // new years on saturday
		{"1/3/2022", "-1 businessday", "12/30/2021"},
		{"12/30/2022", "10 businessday", "1/17/2023"}, // weekends, new years on sunday, mlk day
		{"2/17/2023", "1 businessday", "2/21/2023"},   // presidents day
		{"2/21/2023", "-1 businessday", "2/17/2023"},  // presidents day
		{"5/1/2023", "4 workdays", "5/5/2023"},        // regular weekdays
		{"5/5/2023", "-4 workdays", "5/1/2023"},
		{"5/1/2023", "5 workdays", "5/8/2023"}, // weekend
		{"5/8/2023", "-5 workdays", "5/1/2023"},
		{"5/1/2023", "20 workdays", "5/30/2023"},  // memorial day
		{"6/30/2023", "-9 workdays", "6/16/2023"}, // juneteenth
		{"7/4/2023", "1 workday", "7/5/2023"},
		{"7/3/2023", "1 workday", "7/5/2023"}, // independence day
		{"7/4/2023", "-1 workday", "7/3/2023"},
		{"7/5/2023", "-1 workday", "7/3/2023"},
		{"9/1/2023", "1 workday", "9/5/2023"}, // labor day
		{"9/5/2023", "-1 workday", "9/1/2023"},
		{"10/8/2023", "1 workday", "10/10/2023"}, // columbus day
		{"10/10/2023", "-1 workday", "10/6/2023"},
		{"11/9/2023", "1 workday", "11/13/2023"}, // veterans day
		{"11/13/2023", "-1 workday", "11/9/2023"},
		{"11/22/2023", "1 workday", "11/24/2023"}, // thanksgiving
		{"11/24/2023", "-1 workday", "11/22/2023"},
		{"12/24/2023", "1 workday", "12/26/2023"}, // christmas
		{"12/26/2023", "-1 workday", "12/22/2023"},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("subtest%d", i), func(t *testing.T) {
			input, err := ParseDate(tt.Input)
			if err != nil {
				t.Error(err)
				return
			}
			inter, err := Parse(tt.Interval)
			if err != nil {
				t.Error(err)
				return
			}
			output, err := ParseDate(tt.Output)
			if err != nil {
				t.Error(err)
				return
			}
			got := inter.AddToDate(input)

			if got != output {
				t.Error(fmt.Errorf("Expected %s + %s = %s: Got %s", input, inter, output, got))
				return
			}
		})
	}
}
