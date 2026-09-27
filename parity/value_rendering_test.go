package parity

import "testing"

// Ruby's Float#to_s always shows a decimal point and switches to exponential
// form at exponent >= 15 or <= -5. Go's default formatting drops the ".0" and
// uses a different threshold. Expectations taken from Ruby 3.x + liquid 5.13.0.
func TestFloatRendersLikeRuby(t *testing.T) {
	cases := map[string]string{
		"9.0":                 "9.0",
		"0.5":                 "0.5",
		"4.5":                 "4.5",
		"-9.0":                "-9.0",
		"0.30000000000000004": "0.30000000000000004",
		"100000.0":            "100000.0",
		"123456789.0":         "123456789.0",
		"1.0e14":              "100000000000000.0",
		"1.0e15":              "1.0e+15",
		"1.0e16":              "1.0e+16",
		"0.0001":              "0.0001",
		"1.0e-5":              "1.0e-05",
		"1.5e20":              "1.5e+20",
	}
	for literal, want := range cases {
		assertRender(t, "{{ f }}", want, map[string]interface{}{"f": mustFloat(t, literal)})
	}

	// arithmetic that yields a whole float still shows .0
	assertRender(t, `{{ 4.5 | times: 2 }}`, "9.0", nil)
	assertRender(t, `{{ 10 | minus: 2.0 }}`, "8.0", nil)
}

// Liquid writes an array by joining its elements with no separator.
func TestArrayRendersAsJoinedElements(t *testing.T) {
	data := map[string]interface{}{
		"nums":    []interface{}{1, 2, 3},
		"empty":   []interface{}{},
		"strs":    []interface{}{"a", "b"},
		"withNil": []interface{}{1, nil, 2},
	}
	assertRender(t, `{{ nums }}`, "123", data)
	assertRender(t, `{{ empty }}`, "", data)
	assertRender(t, `{{ strs }}`, "ab", data)
	assertRender(t, `{{ withNil }}`, "12", data)
}

// KNOWN GAP — not fixed.
//
// Ruby Liquid's math filters run through BigDecimal (Utils.to_number converts a
// Float via its decimal string, then apply_operation calls .to_f on the
// result), so 0.1 + 0.2 is exactly 0.3. Go does binary float arithmetic and
// gets 0.30000000000000004, which is what Ruby's *Float* arithmetic gives too —
// the difference is BigDecimal, not formatting.
//
// Matching this needs decimal arithmetic throughout the math filters. No
// Transistor theme does float arithmetic, so it is documented rather than
// fixed. Revisit if a theme ever adds floats.
func TestFloatArithmeticPrecisionIsAKnownGap(t *testing.T) {
	const rubyWouldGive = "0.3"
	got := render(t, `{{ 0.1 | plus: 0.2 }}`, nil)
	if got == rubyWouldGive {
		t.Errorf("decimal arithmetic now matches Ruby (%q) — remove this test and "+
			"the gap note in the README", got)
	}
}
