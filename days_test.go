package interval

import (
	"testing"
	"time"
	"fmt"
)

func TestHolidayWeekday(t *testing.T) {
	tests := []struct {
		Year int
		Month time.Month
		Weekday time.Weekday
		Count int

		Expect string
	}{
		{2023, time.February, time.Monday, 0, "2023-02-06"},
		{2023, time.February, time.Wednesday, 0, "2023-02-01"},
		{2023, time.February, time.Friday, 0, "2023-02-03"},
		{2023, time.February, time.Monday, 2 , "2023-02-20"},
		{2023, time.February, time.Wednesday, 2, "2023-02-15"},
		{2023, time.February, time.Friday, 2, "2023-02-17"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("subtest%d", i), func(t *testing.T) {

			h := holidayWeekday(test.Year, test.Month, test.Weekday, test.Count)
			expect, err := time.Parse("2006-01-02", test.Expect)
			if err != nil {
				t.Error(err)
				return
			}

			if !h.Equal(expect) {
				t.Errorf("Expected %s %d of %s %d to be %s: Got %s",
					test.Weekday,
					test.Count + 1,
					test.Month,
					test.Year,
					test.Expect,
					h.Format("2006-01-02"),
				)
			}
		})
	}
}

func TestHolidayWeekdayReverse(t *testing.T) {
	tests := []struct {
		Year int
		Month time.Month
		Weekday time.Weekday
		Count int

		Expect string
	}{
		{2023, time.February, time.Monday, 0, "2023-02-27"},
		{2023, time.February, time.Tuesday, 0, "2023-02-28"},
		{2023, time.February, time.Wednesday, 0, "2023-02-22"},
		{2023, time.February, time.Monday, 2, "2023-02-13"},
		{2023, time.February, time.Tuesday, 2, "2023-02-14"},
		{2023, time.February, time.Wednesday, 2, "2023-02-08"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("subtest%d", i), func(t *testing.T) {

			h := holidayWeekdayReverse(test.Year, test.Month, test.Weekday, test.Count)
			expect, err := time.Parse("2006-01-02", test.Expect)
			if err != nil {
				t.Error(err)
				return
			}

			if !h.Equal(expect) {
				t.Errorf("Expected reverse %s %d of %s %d to be %s: Got %s",
					test.Weekday,
					test.Count + 1,
					test.Month,
					test.Year,
					test.Expect,
					h.Format("2006-01-02"),
				)
			}
		})
	}
}
