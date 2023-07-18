package interval

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

var usDateRe = regexp.MustCompile(`(\d{1,2})\s*/\s*(\d{1,2})\s*/\s*(\d{4})`)
var pgDateRe = regexp.MustCompile(`(\d{4})\s*-\s*(\d{1,2})\s*-\s*(\d{1,2})`)
var txDateRe = regexp.MustCompile(`([A-Za-z]+)\s*(\d{1,2})[,\s]\s*(\d{4})`)

var months = []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
}

func (d Date) GoTime() time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
}

func FromGoTime(t time.Time) Date {
	y, m, d := t.Date()
	return Date{Year: y, Month: int(m), Day: d}
}

// Returns number of days between two dates
func (d Date) Sub(v Date) int {
	dd := time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, time.UTC)
	vd := time.Date(v.Year, time.Month(v.Month), v.Day, 0, 0, 0, 0, time.UTC)
	return int(dd.Sub(vd).Hours() / 24)
}

func ParseDate(s string) (Date, error) {
	var (
		y, m, d int64
		err     error
	)

	s = strings.TrimSpace(s)

	if match := usDateRe.FindStringSubmatch(s); match != nil {
		// m/d/y
		m, err = strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return Date{}, err
		}
		d, err = strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return Date{}, err
		}
		y, err = strconv.ParseInt(match[3], 10, 64)
		if err != nil {
			return Date{}, err
		}
	} else if match = pgDateRe.FindStringSubmatch(s); match != nil {
		// y-m-d
		y, err = strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return Date{}, err
		}
		m, err = strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return Date{}, err
		}
		d, err = strconv.ParseInt(match[3], 10, 64)
		if err != nil {
			return Date{}, err
		}
	} else if match = txDateRe.FindStringSubmatch(s); match != nil && len(match[1]) > 2 {
		ms := strings.ToLower(match[1])
		ms = ms[:3]
		for i, mm := range months {
			if ms == mm {
				m = int64(i) + 1
			}
		}
		d, err = strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return Date{}, err
		}
		y, err = strconv.ParseInt(match[3], 10, 64)
		if err != nil {
			return Date{}, err
		}
		if m == 0 {
			return Date{}, fmt.Errorf("Interval: Cannot parse month %#v", match[1])
		}
	} else {
		return Date{}, fmt.Errorf("Interval: Cannot parse date %#v", s)
	}

	// Normalize, so October 32 is Nov 1
	t := time.Date(int(y), time.Month(m), int(d), 0, 0, 0, 0, time.UTC)
	ty, tm, td := t.Date()

	return Date{
		Year:  ty,
		Month: int(tm),
		Day:   td,
	}, nil
}
