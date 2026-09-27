package parity

import "testing"

// equalVariables used Go's ==, which is false across dynamic types, so an
// Integer never equalled a Float. Ruby's 1 == 1.0 is true. Combined with
// arithmetic returning float64, `{% assign n = n | plus: 1 %}{% if n == 1 %}`
// never fired.
func TestNumericEqualityAcrossIntAndFloat(t *testing.T) {
	assertRender(t, `{% if a == b %}EQ{% else %}NE{% endif %}`, "EQ",
		map[string]interface{}{"a": 1, "b": 1.0})
	assertRender(t, `{% if a == 1 %}EQ{% else %}NE{% endif %}`, "EQ",
		map[string]interface{}{"a": 1.0})
	assertRender(t, `{% if a != b %}NE{% else %}EQ{% endif %}`, "NE",
		map[string]interface{}{"a": 1, "b": 2.0})

	// the loop-counter shape that surfaced it
	src := `{% assign n = 0 %}{% assign n = n | plus: 1 %}{% if n == 1 %}ONE{% else %}NOT{% endif %}`
	assertRender(t, src, "ONE", nil)
}
