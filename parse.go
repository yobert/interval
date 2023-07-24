package interval

import (
	"errors"
	"fmt"
)

var unitMap = map[string]Interval{
	"s":       Interval{Seconds: 1},
	"sec":     Interval{Seconds: 1},
	"secs":    Interval{Seconds: 1},
	"second":  Interval{Seconds: 1},
	"seconds": Interval{Seconds: 1},

	"m":       Interval{Seconds: 60},
	"min":     Interval{Seconds: 60},
	"mins":    Interval{Seconds: 60},
	"minute":  Interval{Seconds: 60},
	"minutes": Interval{Seconds: 60},

	"h":     Interval{Seconds: 60 * 60},
	"hr":    Interval{Seconds: 60 * 60},
	"hrs":   Interval{Seconds: 60 * 60},
	"hour":  Interval{Seconds: 60 * 60},
	"hours": Interval{Seconds: 60 * 60},

	"d":    Interval{Days: 1},
	"day":  Interval{Days: 1},
	"days": Interval{Days: 1},

	"businessday":  Interval{WorkDays: 1},
	"businessdays": Interval{WorkDays: 1},
	"workday":      Interval{WorkDays: 1},
	"workdays":     Interval{WorkDays: 1},

	"w":     Interval{Days: 7},
	"wk":    Interval{Days: 7},
	"wks":   Interval{Days: 7},
	"week":  Interval{Days: 7},
	"weeks": Interval{Days: 7},

	"mon":    Interval{Months: 1},
	"mons":   Interval{Months: 1},
	"month":  Interval{Months: 1},
	"months": Interval{Months: 1},

	"year":  Interval{Months: 12},
	"years": Interval{Months: 12},
}

func Parse(s string) (Interval, error) {
	orig := s

	var out Interval

	for len(s) > 0 {
		// Ignore whitespace
		if s[0] == ' ' {
			s = s[1:]
			continue
		}

		// Consume [-+]?
		neg := 1
		if s[0] == '-' || s[0] == '+' {
			if s[0] == '-' {
				neg = -1
			}
			s = s[1:]
		}

		if s == "" {
			return Interval{}, fmt.Errorf("Interval: Cannot parse %#v", orig)
		}

		for s != "" {
			var (
				v   uint64
				err error
			)

			// Consume [0-9]+
			pl := len(s)
			v, s, err = leadingInt(s)
			if err != nil || pl == len(s) {
				return Interval{}, fmt.Errorf("Interval: Cannot parse %#v", orig)
			}

			// Consume whitespace
			for s != "" && s[0] == ' ' {
				s = s[1:]
			}

			// Consume unit
			i := 0
			for ; i < len(s); i++ {
				c := s[i]
				if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') {
					break
				}
			}
			if i == 0 {
				return Interval{}, fmt.Errorf("Interval: missing unit in %#v", orig)
			}
			u := s[:i]
			s = s[i:]
			unit, ok := unitMap[u]
			if !ok {
				return Interval{}, fmt.Errorf("Interval: Unknown unit %#v in %#v", u, orig)
			}
			vi := int(v) * neg
			out.Seconds += unit.Seconds * vi
			out.WorkDays += unit.WorkDays * vi
			out.Days += unit.Days * vi
			out.Months += unit.Months * vi
		}
	}

	return out, nil
}

// Some code lifted from Go's stdlib time package :)
var errLeadingInt = errors.New("Interval: bad [0-9]*") // never printed

// leadingInt consumes the leading [0-9]* from s.
func leadingInt[bytes []byte | string](s bytes) (x uint64, rem bytes, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			// overflow
			return 0, rem, errLeadingInt
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			// overflow
			return 0, rem, errLeadingInt
		}
	}
	return x, s[i:], nil
}
