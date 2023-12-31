package interval

import (
	"fmt"
	"time"
)

type Interval struct {
	Seconds int
	Days    int
	Months  int

	WorkDays int
}

var zeroInterval Interval

func fmtnum(v int, s1 string, s2 string) string {
	if v == 1 || v == -1 {
		return fmt.Sprintf("%d %s", v, s1)
	}
	return fmt.Sprintf("%d %s", v, s2)
}
func (i Interval) String() string {
	if i == zeroInterval {
		return "00:00:00"
	}
	out := ""
	m := i.Months
	if m > 11 || m < -11 {
		out += fmtnum(i.Months/12, "year", "years")
		m = m % 12
	}

	if m != 0 {
		if out != "" {
			out += " "
		}
		out += fmtnum(m, "mon", "mons")
	}

	if i.WorkDays != 0 {
		if out != "" {
			out += " "
		}
		out += fmtnum(i.WorkDays, "workday", "workdays")
	}

	if i.Days != 0 {
		if out != "" {
			out += " "
		}
		out += fmtnum(i.Days, "day", "days")
	}

	if i.Seconds != 0 {
		if out != "" {
			out += " "
		}

		hours := i.Seconds / 60 / 60
		minutes := i.Seconds / 60 % 60
		seconds := i.Seconds % 60

		out += fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return out
}

// Add two intervals together
func (i Interval) Add(v Interval) Interval {
	return Interval{
		Seconds:  i.Seconds + v.Seconds,
		WorkDays: i.WorkDays + v.WorkDays,
		Days:     i.Days + v.Days,
		Months:   i.Months + v.Months,
	}
}

// Multiply an interval by an integer value
func (i Interval) MultInt(v int) Interval {
	return Interval{
		Seconds:  i.Seconds * v,
		WorkDays: i.WorkDays * v,
		Days:     i.Days * v,
		Months:   i.Months * v,
	}
}

// Add an interval to a date
func (i Interval) AddToDate(d Date) Date {
	v := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)

	// Move months
	if i.Months != 0 {
		v = v.AddDate(0, i.Months, 0)

		// Check if we month rolled. If we did, rewind by 1 day at a time until we're right.
		_, _, newd := v.Date()
		if newd < d.Day && d.Day > 27 {
			for {
				_, _, newd = v.Date()
				if newd > 27 {
					break
				}
				v = v.AddDate(0, 0, -1)
			}
		}
	}

	// Move business days
	// TODO this will get slow for high numbers.
	if i.WorkDays != 0 {
		var (
			dir, abs int
		)
		if i.WorkDays > 0 {
			dir = 1
			abs = i.WorkDays
		} else {
			dir = -1
			abs = -i.WorkDays
		}

		for abs > 0 {
			v = v.AddDate(0, 0, dir)
			abs--
			for !isBusinessDay(v) {
				v = v.AddDate(0, 0, dir)
			}
		}
	}

	// Move days
	if i.Days != 0 {
		v = v.AddDate(0, 0, i.Days)
	}

	// Move seconds
	if i.Seconds != 0 {
		v = v.Add(time.Duration(i.Seconds) * time.Second)
	}

	ry, rm, rd := v.Date()
	return Date{Year: ry, Month: int(rm), Day: rd}
}
