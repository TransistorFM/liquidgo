package parity

import "testing"

// strftime substituted Go layout tokens into the format string, so flagged
// directives (%-d), %l and %s were left verbatim in the output, and literal
// text could collide with a Go reference-time token. Expectations come from
// Ruby 3.x: Time.parse("2024-03-05 09:07:03 -0500").strftime(...)
func TestStrftimeMatchesRuby(t *testing.T) {
	data := map[string]interface{}{"t": "2024-03-05 09:07:03 -0500"}

	cases := map[string]string{
		"%a": "Tue", "%A": "Tuesday", "%b": "Mar", "%B": "March",
		"%d": "05", "%-d": "5", "%e": " 5",
		"%H": "09", "%-H": "9", "%I": "09", "%-I": "9", "%l": " 9",
		"%j": "065", "%-j": "65", "%m": "03", "%-m": "3",
		"%M": "07", "%S": "03", "%p": "AM", "%P": "am",
		"%s": "1709647623", "%w": "2", "%U": "09", "%W": "10",
		"%y": "24", "%Y": "2024", "%z": "-0500", "%:z": "-05:00", "%Z": "EST",
		"%F": "2024-03-05", "%T": "09:07:03", "%D": "03/05/24",
		"%R": "09:07", "%r": "09:07:03 AM", "%L": "000", "%N": "000000000",
		"%x": "03/05/24", "%X": "09:07:03", "%%": "%",
	}
	for format, want := range cases {
		assertRender(t, `{{ t | date: "`+format+`" }}`, want, data)
	}

	// common theme formats, and literal text that must survive
	assertRender(t, `{{ t | date: "%b %d, %Y" }}`, "Mar 05, 2024", data)
	assertRender(t, `{{ t | date: "Published %-d %B %Y" }}`, "Published 5 March 2024", data)
	assertRender(t, `{{ t | date: "Jan 2 2006 15:04" }}`, "Jan 2 2006 15:04", data)
}
