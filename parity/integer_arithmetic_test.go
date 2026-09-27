package parity

import "testing"

// ToNumber normalises every numeric to float64, so applyOperation's
// `inputNum.(int)` check never held and Integer op Integer produced a Float.
// That is visible in output because Ruby prints a Float as "6.0" and an
// Integer as "6", and it breaks `{% if n == 1 %}` after `| plus: 1`.
//
// Ruby's Integer#/ is also integer division, flooring toward negative infinity.
func TestIntegerArithmeticStaysInteger(t *testing.T) {
	assertRender(t, `{{ 1 | plus: 2 }}`, "3", nil)
	assertRender(t, `{{ 10 | minus: 3 }}`, "7", nil)
	assertRender(t, `{{ 3 | times: 2 }}`, "6", nil)
	assertRender(t, `{{ 7 | modulo: 3 }}`, "1", nil)
	assertRender(t, `{{ "3" | plus: "4" }}`, "7", nil)

	// any Float makes the result a Float
	assertRender(t, `{{ 4.5 | times: 2 }}`, "9.0", nil)
	assertRender(t, `{{ 10 | minus: 2.0 }}`, "8.0", nil)
	assertRender(t, `{{ 1.0 | times: 3 }}`, "3.0", nil)
}

func TestIntegerDivisionMatchesRuby(t *testing.T) {
	assertRender(t, `{{ 12 | divided_by: 5 }}`, "2", nil)
	assertRender(t, `{{ 7 | divided_by: 2 }}`, "3", nil)
	assertRender(t, `{{ 16 | divided_by: 4 }}`, "4", nil)
	assertRender(t, `{{ 12 | divided_by: 5.0 }}`, "2.4", nil)
}

// sum and abs had the same root cause: ToNumber normalises to float64, so
// their int checks never held and `{{ posts | sum: 'count' }}` rendered "43.0"
// where Ruby gives "43". Upstream's comprehensive_test.go covers the sum case,
// but that suite encodes non-Ruby behaviour elsewhere, so it is not relied on.
func TestSumAndAbsPreserveIntegers(t *testing.T) {
	data := map[string]interface{}{
		"ints":   []interface{}{40, 3},
		"floats": []interface{}{1.5, 2.5},
		"mixed":  []interface{}{1, 2.5},
		"rows": []interface{}{
			map[string]interface{}{"count": 40},
			map[string]interface{}{"count": 3},
		},
	}

	assertRender(t, `{{ ints | sum }}`, "43", data)
	assertRender(t, `{{ floats | sum }}`, "4.0", data)
	assertRender(t, `{{ mixed | sum }}`, "3.5", data)
	assertRender(t, `{{ rows | sum: "count" }}`, "43", data)

	assertRender(t, `{{ -5 | abs }}`, "5", data)
	assertRender(t, `{{ -5.5 | abs }}`, "5.5", data)
	assertRender(t, `{{ 17 | abs }}`, "17", data)
}
