package parity

import "testing"

// Ruby iterates a Hash as [key, value] pairs. This returned an empty segment
// for any map, so `{% for pair in hash %}` rendered nothing at all — silent
// content loss with no error.
//
// Keys are sorted rather than insertion-ordered: Go maps have no order to
// recover. A value implementing liquid.OrderedKeys can supply Ruby's order.
func TestForIteratesMapsAsKeyValuePairs(t *testing.T) {
	data := map[string]interface{}{
		"links": map[string]interface{}{"apple": "a", "spotify": "s", "rss": "r"},
		"empty": map[string]interface{}{},
	}

	assertRender(t, `{% for pair in links %}{{ pair[0] }}={{ pair[1] }};{% endfor %}`,
		"apple=a;rss=r;spotify=s;", data)

	// the first/last aliases Ruby exposes on the pair
	assertRender(t, `{% for pair in links %}{{ pair.first }}{% endfor %}`, "applerssspotify", data)

	assertRender(t, `{% for pair in links limit: 2 %}{{ pair[0] }},{% endfor %}`, "apple,rss,", data)
	assertRender(t, `{% for pair in empty %}x{% else %}none{% endfor %}`, "none", data)

	// forloop metadata must reflect the real size
	assertRender(t, `{% for pair in links %}{{ forloop.length }}{% break %}{% endfor %}`, "3", data)
}
