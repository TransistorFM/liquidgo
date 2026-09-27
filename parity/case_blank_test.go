package parity

import "testing"

// CaseTag.Blank() reported the tag's own nodelist (the text between {% case %}
// and the first {% when %}, always whitespace) instead of the when/else bodies.
// An enclosing block then considered itself blank and stripped its whitespace
// text nodes. Ruby: @blocks.none? { |block| !block.attachment.blank? }
func TestCaseTagBlankReflectsWhenBodies(t *testing.T) {
	// The case renders content, so surrounding whitespace must survive.
	src := "{% if true %}\n  {% case kind %}\n  {% when \"a\" %}A{% endcase %}\n{% endif %}"
	assertRender(t, src, "\n  A\n", map[string]interface{}{"kind": "a"})

	// A genuinely blank case stays blank.
	blank := "{% if true %}X{% case kind %}{% when \"a\" %}{% endcase %}Y{% endif %}"
	assertRender(t, blank, "XY", map[string]interface{}{"kind": "a"})
}
