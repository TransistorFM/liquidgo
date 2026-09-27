package liquid

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// strftime formats a time using Ruby's Time#strftime directives.
//
// It scans the format rather than substituting Go layout tokens into it, so
// flags (%-d, %_d, %0e) work and literal text can never collide with a Go
// reference-time token. Expectations are taken from Ruby 3.x.
func strftime(t *time.Time, format string) string {
	if t == nil {
		return ""
	}

	var b strings.Builder
	for i := 0; i < len(format); i++ {
		if format[i] != '%' || i == len(format)-1 {
			b.WriteByte(format[i])
			continue
		}

		i++
		pad := byte(0) // 0 = directive default; '-' none, '_' space, '0' zero
		colons := 0

		for i < len(format) {
			c := format[i]
			if c == '-' || c == '_' || c == '0' {
				pad = c
				i++
				continue
			}
			if c == ':' {
				colons++
				i++
				continue
			}
			break
		}
		if i >= len(format) {
			b.WriteByte('%')
			break
		}

		b.WriteString(strftimeDirective(t, format[i], pad, colons))
	}

	return b.String()
}

func strftimeDirective(t *time.Time, verb byte, pad byte, colons int) string {
	switch verb {
	case 'a':
		return t.Format("Mon")
	case 'A':
		return t.Format("Monday")
	case 'b', 'h':
		return t.Format("Jan")
	case 'B':
		return t.Format("January")
	case 'c':
		return fmt.Sprintf("%s %s %s %s", t.Format("Mon"), t.Format("Jan"),
			padNum(t.Day(), 2, '_'), t.Format("15:04:05 2006"))
	case 'd':
		return padNum(t.Day(), 2, defaultPad(pad, '0'))
	case 'e':
		return padNum(t.Day(), 2, defaultPad(pad, '_'))
	case 'H':
		return padNum(t.Hour(), 2, defaultPad(pad, '0'))
	case 'I':
		return padNum(hour12(t), 2, defaultPad(pad, '0'))
	case 'l':
		return padNum(hour12(t), 2, defaultPad(pad, '_'))
	case 'j':
		return padNum(t.YearDay(), 3, defaultPad(pad, '0'))
	case 'm':
		return padNum(int(t.Month()), 2, defaultPad(pad, '0'))
	case 'M':
		return padNum(t.Minute(), 2, defaultPad(pad, '0'))
	case 'S':
		return padNum(t.Second(), 2, defaultPad(pad, '0'))
	case 'p':
		return t.Format("PM")
	case 'P':
		return strings.ToLower(t.Format("PM"))
	case 's':
		return strconv.FormatInt(t.Unix(), 10)
	case 'w':
		return strconv.Itoa(int(t.Weekday()))
	case 'u':
		d := int(t.Weekday())
		if d == 0 {
			d = 7
		}
		return strconv.Itoa(d)
	case 'U':
		return padNum(weekOfYear(t, time.Sunday), 2, defaultPad(pad, '0'))
	case 'W':
		return padNum(weekOfYear(t, time.Monday), 2, defaultPad(pad, '0'))
	case 'y':
		return padNum(t.Year()%100, 2, defaultPad(pad, '0'))
	case 'Y':
		return strconv.Itoa(t.Year())
	case 'C':
		return padNum(t.Year()/100, 2, defaultPad(pad, '0'))
	case 'z':
		if colons == 1 {
			return t.Format("-07:00")
		}
		return t.Format("-0700")
	case 'Z':
		return t.Format("MST")
	case 'L':
		return fmt.Sprintf("%03d", t.Nanosecond()/1e6)
	case 'N':
		return fmt.Sprintf("%09d", t.Nanosecond())
	case 'F':
		return t.Format("2006-01-02")
	case 'T', 'X':
		return t.Format("15:04:05")
	case 'D', 'x':
		return t.Format("01/02/06")
	case 'R':
		return t.Format("15:04")
	case 'r':
		return t.Format("03:04:05 PM")
	case 'n':
		return "\n"
	case 't':
		return "\t"
	case '%':
		return "%"
	default:
		// Unknown directive: Ruby emits it verbatim.
		return "%" + string(verb)
	}
}

func defaultPad(flag, fallback byte) byte {
	if flag == 0 {
		return fallback
	}
	return flag
}

func padNum(n, width int, pad byte) string {
	s := strconv.Itoa(n)
	if pad == '-' || len(s) >= width {
		return s
	}
	fill := byte('0')
	if pad == '_' {
		fill = ' '
	}
	return strings.Repeat(string(fill), width-len(s)) + s
}

func hour12(t *time.Time) int {
	h := t.Hour() % 12
	if h == 0 {
		h = 12
	}
	return h
}

// weekOfYear counts weeks starting on the given weekday, as %U and %W do.
func weekOfYear(t *time.Time, start time.Weekday) int {
	yday := t.YearDay() - 1
	wday := (int(t.Weekday()) - int(start) + 7) % 7
	return (yday - wday + 7) / 7
}
