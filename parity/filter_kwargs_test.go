package parity

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

type kwFilters struct{}

// Join reports the keyword args it received, so a dropped hash is visible.
func (kwFilters) Join(input interface{}, args ...interface{}) string {
	parts := []string{}
	for _, a := range args {
		if m, ok := a.(map[string]interface{}); ok {
			for k, v := range m {
				parts = append(parts, fmt.Sprintf("%s=%v", k, v))
			}
		}
	}
	sort.Strings(parts)
	return fmt.Sprintf("%v[%s]", input, strings.Join(parts, ","))
}

// Variable parsed `| f: a: 1, b: 2` into a keyword hash but Render never read
// filter[2], so custom filters were invoked without their keyword args. Ruby
// appends the evaluated hash as a final positional argument.
func TestFilterKeywordArgumentsReachFilters(t *testing.T) {
	assertRender(t, `{{ "k" | join: first: 1, last: 3 }}`, "k[first=1,last=3]", nil, &kwFilters{})
	assertRender(t, `{{ "k" | join: total: count }}`, "k[total=7]",
		map[string]interface{}{"count": 7}, &kwFilters{})
}
