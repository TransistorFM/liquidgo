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

// Integer equality must not round through float64, including when one side is
// a Float: Ruby distinguishes 9007199254740993 from 9007199254740992.0.
func TestNumericEqualityPreservesIntegerPrecision(t *testing.T) {
	cases := []struct {
		name        string
		left, right interface{}
		want        string
	}{
		{"adjacent signed integers", int64(1 << 53), int64(1<<53 + 1), "NE|NE"},
		{"adjacent unsigned integers", uint64(1 << 63), uint64(1<<63 + 1), "NE|NE"},
		{"equal signed and unsigned", int64(1<<53 + 1), uint64(1<<53 + 1), "EQ|EQ"},
		{"negative and unsigned", int64(-1), ^uint64(0), "NE|NE"},
		{"rounded float", int64(1<<53 + 1), float64(1 << 53), "NE|NE"},
		{"exact float", int64(1 << 53), float64(1 << 53), "EQ|EQ"},
		{"fractional float", int64(3), 3.5, "NE|NE"},
		{"unsigned fractional float", uint64(3), 3.5, "NE|NE"},
		{"minimum signed integer", int64(-1 << 63), -0x1p63, "EQ|EQ"},
		{"signed float boundary", int64(1<<63 - 1), 0x1p63, "NE|NE"},
		{"unsigned float boundary", ^uint64(0), 0x1p64, "NE|NE"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			src := `{% if a == b %}EQ{% else %}NE{% endif %}|{% if a != b %}NE{% else %}EQ{% endif %}`
			assertRender(t, src, c.want, map[string]interface{}{"a": c.left, "b": c.right})
			assertRender(t, src, c.want, map[string]interface{}{"a": c.right, "b": c.left})
		})
	}
}
