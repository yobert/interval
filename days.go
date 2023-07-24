package interval

import (
	"time"
)

// Magic sauce for "business days" which means holiday calculation, according to:
// 5 U.S.C. 6103(a)
// 5 U.S.C. 6103(b)
// Section 3(a) of Executive Order 11582, February 11, 1971
// https://www.opm.gov/policy-data-oversight/pay-leave/work-schedules/fact-sheets/Federal-Holidays-In-Lieu-Of-Determination

func isBusinessDay(v time.Time) bool {
	// US rules.

	wd := v.Weekday()

	if wd == time.Sunday || wd == time.Saturday {
		return false
	}

	y, m, _ := v.Date()

	switch m {
	case time.January:
		// New years day
		if v.Equal(holidayDate(y, 1, 1)) {
			return false
		}
		// MLK day: Third monday in Jan
		if v.Equal(holidayWeekday(y, time.January, time.Monday, 2)) {
			return false
		}
	case time.February:
		// Presidents day: Third monday in Feb
		if v.Equal(holidayWeekday(y, time.February, time.Monday, 2)) {
			return false
		}
	case time.May:
		// Memorial day: Last monday in May
		if v.Equal(holidayWeekdayReverse(y, time.May, time.Monday, 0)) {
			return false
		}
	case time.June:
		// Juneteenth
		if v.Equal(holidayDate(y, time.June, 19)) {
			return false
		}
	case time.July:
		// Independence day
		if v.Equal(holidayDate(y, time.July, 4)) {
			return false
		}
	case time.September:
		// Labor day: First monday in Sept
		if v.Equal(holidayWeekday(y, time.September, time.Monday, 0)) {
			return false
		}
	case time.October:
		// Columbus day
		if v.Equal(holidayWeekday(y, time.October, time.Monday, 1)) {
			return false
		}
	case time.November:
		// Veterans day
		if v.Equal(holidayDate(y, time.November, 11)) {
			return false
		}
		// Thanksgiving: 4th Thursday in Nov
		if v.Equal(holidayWeekday(y, time.November, time.Thursday, 3)) {
			return false
		}
	case time.December:
		// Christmas
		if v.Equal(holidayDate(y, time.December, 25)) {
			return false
		}
		// New years day next year
		if v.Equal(holidayDate(y+1, 1, 1)) {
			return false
		}
	}
	return true
}

func holidayDate(year int, month time.Month, day int) time.Time {
	hd := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	wd := hd.Weekday()
	if wd == time.Saturday {
		hd = hd.AddDate(0, 0, -1)
	}
	if wd == time.Sunday {
		hd = hd.AddDate(0, 0, 1)
	}
	return hd
}

func holidayWeekday(year int, month time.Month, day time.Weekday, count int) time.Time {
	v := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	vwd := v.Weekday()
	offset := int(day - vwd)
	if offset < 0 {
		offset += 7
	}
	v = v.AddDate(0, 0, offset)

	for {
		if count == 0 {
			return v
		}
		v = v.AddDate(0, 0, 7)
		count--
	}
}

func holidayWeekdayReverse(year int, month time.Month, day time.Weekday, count int) time.Time {
	v := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	v = v.AddDate(0, 0, -1)

	vwd := v.Weekday()
	offset := int(vwd - day)
	if offset < 0 {
		offset += 7
	}
	v = v.AddDate(0, 0, -offset)

	for {
		if count == 0 {
			return v
		}
		v = v.AddDate(0, 0, -7)
		count--
	}
}
