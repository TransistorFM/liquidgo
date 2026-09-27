package parity

import "testing"

// The {% liquid %} tag is inherently multi-line, but blockBodyFullToken was
// compiled without (?s) so `.` never matched a newline and the whole tag fell
// through to the unknown-tag handler.
func TestLiquidTagSpansMultipleLines(t *testing.T) {
	assertRender(t, "{% liquid\n  assign x = 1\n%}{{ x }}", "1", nil)
	assertRender(t, "{% liquid\n  assign a = \"A\"\n  assign b = \"B\"\n  echo a\n  echo b\n%}", "AB", nil)
}
