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

var usDateRe = regexp.MustCompile(`(\d{1,2})/(\d{1,2})/(\d{4})`)
var pgDateRe = regexp.MustCompile(`(\d{4})-(\d{1,2})-(\d{1,2})`)

func (d Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.Day)
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
