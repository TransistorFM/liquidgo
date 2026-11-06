package parity

import (
	"os"
	"os/exec"
	"runtime/debug"
	"testing"

	"github.com/Notifuse/liquidgo/liquid"
)

// Ruby compares hashes and arrays by value. Comparing episode hashes with
// Go's == panicked, rendering Beacon and Cardboard as "Liquid error: internal".
func TestCollectionEquality(t *testing.T) {
	drop := liquid.NewDrop()
	cases := []struct {
		name        string
		left, right interface{}
		want        string
	}{
		{"equal episodes", map[string]interface{}{"path": "/recommended", "type": "full"}, map[string]interface{}{"type": "full", "path": "/recommended"}, "EQ|EQ"},
		{"different episodes", map[string]interface{}{"path": "/recommended"}, map[string]interface{}{"path": "/latest"}, "NE|NE"},
		{"different keys", map[string]interface{}{"path": "/episode"}, map[string]interface{}{"title": "/episode"}, "NE|NE"},
		{"extra hash key", map[string]interface{}{"path": "/episode"}, map[string]interface{}{"path": "/episode", "title": nil}, "NE|NE"},
		{"nested collections", map[string]interface{}{"people": []interface{}{map[string]interface{}{"name": "Justin"}}}, map[string]interface{}{"people": []interface{}{map[string]interface{}{"name": "Justin"}}}, "EQ|EQ"},
		{"equal arrays", []interface{}{1, "two"}, []interface{}{1, "two"}, "EQ|EQ"},
		{"array order matters", []interface{}{1, 2}, []interface{}{2, 1}, "NE|NE"},
		{"array length matters", []interface{}{1}, []interface{}{1, 2}, "NE|NE"},
		{"missing episode", map[string]interface{}{"path": "/recommended"}, nil, "NE|NE"},
		{"numeric array elements", []interface{}{1, 2.0}, []interface{}{1.0, 2}, "EQ|EQ"},
		{"large integers in arrays", [1]int64{1 << 53}, [1]int64{1<<53 + 1}, "NE|NE"},
		{"large integers in slices", []int64{1 << 53}, []int64{1<<53 + 1}, "NE|NE"},
		{"numeric hash values", map[string]interface{}{"id": 1}, map[string]interface{}{"id": 1.0}, "EQ|EQ"},
		{"nested numeric values", map[string]interface{}{"ids": []interface{}{1}}, map[string]interface{}{"ids": []interface{}{1.0}}, "EQ|EQ"},
		{"numeric hash keys stay distinct", map[interface{}]interface{}{1: "one"}, map[interface{}]interface{}{1.0: "one"}, "NE|NE"},
		{"typed array", []string{"one"}, []interface{}{"one"}, "EQ|EQ"},
		{"empty arrays", []interface{}(nil), []interface{}{}, "EQ|EQ"},
		{"distinct drops", drop, liquid.NewDrop(), "NE|NE"},
		{"same drop", drop, drop, "EQ|EQ"},
		{"nested distinct drops", map[string]interface{}{"people": []interface{}{drop}}, map[string]interface{}{"people": []interface{}{liquid.NewDrop()}}, "NE|NE"},
		{"nested same drop", map[string]interface{}{"people": []interface{}{drop}}, map[string]interface{}{"people": []interface{}{drop}}, "EQ|EQ"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assertRender(t, `{% if a == b %}EQ{% else %}NE{% endif %}|{% if a != b %}NE{% else %}EQ{% endif %}`,
				c.want, map[string]interface{}{"a": c.left, "b": c.right})
		})
	}
}

func TestCyclicCollectionEquality(t *testing.T) {
	// A stack overflow cannot be recovered. Isolate the render so a regression
	// fails this test instead of killing the entire suite, with a bounded stack.
	if os.Getenv("LIQUIDGO_CYCLE_TEST") != "1" {
		cmd := exec.Command(os.Args[0], "-test.run=^TestCyclicCollectionEquality$")
		cmd.Env = append(os.Environ(), "LIQUIDGO_CYCLE_TEST=1")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("cyclic collection rendering failed: %v\n%s", err, out)
		}
		return
	}
	debug.SetMaxStack(128 * 1024)

	a, b := map[string]interface{}{}, map[string]interface{}{}
	a["self"], b["self"] = a, b
	x, y := make([]interface{}, 2), make([]interface{}, 2)
	x[0], y[0] = x, y
	x[1], y[1] = 1, 1
	src := `{% if a == b %}EQ{% else %}NE{% endif %}|{% if a != b %}NE{% else %}EQ{% endif %}`

	assertRender(t, src, "EQ|EQ", map[string]interface{}{"a": a, "b": a})
	assertRender(t, src, "EQ|EQ", map[string]interface{}{"a": a, "b": b})
	assertRender(t, src, "EQ|EQ", map[string]interface{}{"a": x, "b": x})
	assertRender(t, src, "EQ|EQ", map[string]interface{}{"a": x, "b": y})

	// A matching cycle must not hide a later unequal value.
	a["value"], b["value"] = 1, 2
	x[1], y[1] = 1, 2
	assertRender(t, src, "NE|NE", map[string]interface{}{"a": a, "b": b})
	assertRender(t, src, "NE|NE", map[string]interface{}{"a": x, "b": y})
}
