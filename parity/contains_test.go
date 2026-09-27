package parity

import "testing"

// containsOperator applied a substring match to every array element, so
// ["about-us"] contains "about" was true. Ruby's Array#include? is exact
// equality for strings; only String#include? is a substring test.
//
// Substring matching over non-string elements (structs, drops) is preserved —
// see TestConditionBlogPostContains in the integration suite.
func TestContainsIsExactForStringArrays(t *testing.T) {
	data := map[string]interface{}{
		"pages": []interface{}{"about-us", "contact"},
		"title": "about-us",
	}
	assertRender(t, `{% if pages contains "about" %}YES{% else %}NO{% endif %}`, "NO", data)
	assertRender(t, `{% if pages contains "about-us" %}YES{% else %}NO{% endif %}`, "YES", data)
	// strings keep substring semantics
	assertRender(t, `{% if title contains "about" %}YES{% else %}NO{% endif %}`, "YES", data)
}
