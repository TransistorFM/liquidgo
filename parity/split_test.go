package parity

import "testing"

// Ruby's String#split returns [] for an empty string and drops trailing empty
// strings; Go's strings.Split keeps them. Themes rely on this in the common
// `{% assign only = only | default: "" | split: "," %}{% if only.size == 0 %}`
// idiom, which silently selects nothing when size is 1 instead of 0.
func TestSplitMatchesRubyStringSplit(t *testing.T) {
	assertRender(t, `{{ "" | split: "," | size }}`, "0", nil)
	assertRender(t, `{{ "a,b," | split: "," | size }}`, "2", nil)
	assertRender(t, `{{ "a,,b" | split: "," | size }}`, "3", nil)
	assertRender(t, `{{ "a,b" | split: "," | size }}`, "2", nil)
	assertRender(t, `{% assign only = "" | split: "," %}{% if only.size == 0 %}EMPTY{% else %}NOT{% endif %}`, "EMPTY", nil)
}
