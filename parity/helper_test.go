// Package parity holds regression tests for behaviours where this engine
// diverged from the Ruby Shopify Liquid gem (5.13.0). Each test names the Ruby
// semantics it pins down, and each is self-contained: no external fixtures.
package parity

import (
	"strconv"
	"testing"

	"github.com/Notifuse/liquidgo/liquid"
	"github.com/Notifuse/liquidgo/liquid/tags"
)

// files backs {% include %} in tests that need a partial.
var files = map[string]string{}

type testFS struct{}

func (testFS) ReadTemplateFile(path string) (string, error) {
	if s, ok := files[path]; ok {
		return s, nil
	}
	return "", nil
}

// render parses and renders src with the standard tags registered.
func render(t *testing.T, src string, data map[string]interface{}, filters ...interface{}) string {
	t.Helper()

	env := liquid.NewEnvironment()
	tags.RegisterStandardTags(env)
	env.SetFileSystem(testFS{})
	for _, f := range filters {
		if err := env.RegisterFilter(f); err != nil {
			t.Fatalf("RegisterFilter: %v", err)
		}
	}

	tmpl, err := liquid.ParseTemplate(src, &liquid.TemplateOptions{Environment: env})
	if err != nil {
		t.Fatalf("parse %q: %v", src, err)
	}
	return tmpl.Render(data, nil)
}

// mustFloat parses a Go float literal for table-driven expectations.
func mustFloat(t *testing.T, lit string) float64 {
	t.Helper()
	f, err := strconv.ParseFloat(lit, 64)
	if err != nil {
		t.Fatalf("bad float literal %q: %v", lit, err)
	}
	return f
}

func assertRender(t *testing.T, src, want string, data map[string]interface{}, filters ...interface{}) {
	t.Helper()
	if got := render(t, src, data, filters...); got != want {
		t.Errorf("template %q\n  got:  %q\n  want: %q (Ruby liquid 5.13.0)", src, got, want)
	}
}
