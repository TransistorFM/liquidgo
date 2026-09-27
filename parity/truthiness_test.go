package parity

import "testing"

// Ruby treats only nil and false as falsy; "" , 0 and [] are truthy. Four
// call sites also excluded "", so `{% if title %}` was false for an empty
// string where Ruby renders the branch.
func TestEmptyStringIsTruthy(t *testing.T) {
	data := map[string]interface{}{
		"blank": "", "zero": 0, "nothing": nil, "list": []interface{}{}, "no": false,
	}
	assertRender(t, `{% if blank %}T{% else %}F{% endif %}`, "T", data)
	assertRender(t, `{% if zero %}T{% else %}F{% endif %}`, "T", data)
	assertRender(t, `{% if list %}T{% else %}F{% endif %}`, "T", data)
	assertRender(t, `{% if nothing %}T{% else %}F{% endif %}`, "F", data)
	assertRender(t, `{% if no %}T{% else %}F{% endif %}`, "F", data)
	assertRender(t, `{% unless blank %}T{% else %}F{% endunless %}`, "F", data)
	assertRender(t, `{% case blank %}{% when "" %}MATCH{% else %}NO{% endcase %}`, "MATCH", data)
}

// The bare-variable and `or` paths in Condition had their own copies of the
// truthiness test, so `{% if obj.blank %}` and `{% if a or b %}` still treated
// "" as falsy after the tag-level fix.
func TestEmptyStringIsTruthyInConditionPaths(t *testing.T) {
	data := map[string]interface{}{
		"obj":     map[string]interface{}{"blank": "", "nothing": nil},
		"blank":   "",
		"nothing": nil,
	}
	assertRender(t, `{% if obj.blank %}T{% else %}F{% endif %}`, "T", data)
	assertRender(t, `{% if nothing or blank %}T{% else %}F{% endif %}`, "T", data)
	assertRender(t, `{% if blank and obj.blank %}T{% else %}F{% endif %}`, "T", data)
	assertRender(t, `{% if obj.nothing %}T{% else %}F{% endif %}`, "F", data)
}
