package parity

import "testing"

// Tokenizer.Shift returns "" both at end-of-input and for a blank line, so
// parseForLiquidTag treated the first blank line as EOF and silently discarded
// every statement after it. Ruby's `while (token = tokenizer.shift)` stops only
// on nil; an empty line is truthy and simply skipped.
//
// This failed silently: no error, just assigns that never happened.
func TestLiquidTagBlankLineIsNotEndOfInput(t *testing.T) {
	assertRender(t, "{% liquid\n  assign a = \"1\"\n\n  assign b = \"2\"\n%}{{ a }}{{ b }}", "12", nil)

	// blank line after a nested block, the shape that surfaced this
	src := "{% liquid\n  assign x = \"\"\n  if x == \"\"\n    assign x = \"set\"\n  endif\n\n  assign y = \"after\"\n%}{{ x }}|{{ y }}"
	assertRender(t, src, "set|after", nil)

	// a variable assigned "" after a blank line must be "", not nil
	src2 := "{% liquid\n  assign a = \"1\"\n\n  assign fb = \"\"\n%}{% if fb == '' %}EMPTY{% else %}NIL{% endif %}"
	assertRender(t, src2, "EMPTY", nil)
}
