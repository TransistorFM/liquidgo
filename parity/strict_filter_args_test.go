package parity

import (
	"testing"

	"github.com/Notifuse/liquidgo/liquid"
	"github.com/Notifuse/liquidgo/liquid/tags"
)

// In strict mode a filter colon promises an argument. This accepted
// `| truncate:` and trailing commas, silently applying the filter with no
// arguments, where Ruby raises:
//
//	Liquid syntax error: [:end_of_string] is not a valid expression
func TestStrictModeRejectsMissingFilterArgument(t *testing.T) {
	rejected := []string{
		`{{ "x" | truncate: }}`,
		`{{ "x" | truncate: 5, }}`,
		`{{ "x" | replace: "a", }}`,
	}
	for _, src := range rejected {
		if err := parseStrict(src); err == nil {
			t.Errorf("template %q parsed, want a syntax error (Ruby liquid 5.13.0 raises)", src)
		}
	}

	accepted := []string{
		`{{ "x" | truncate: 5 }}`,
		`{{ "x" | replace: "a", "b" }}`,
		`{{ "x" | upcase }}`,
		`{{ "x" | t: first: 1, last: 2 }}`,
	}
	for _, src := range accepted {
		if err := parseStrict(src); err != nil {
			t.Errorf("template %q failed to parse: %v", src, err)
		}
	}
}

func parseStrict(src string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = liquid.NewSyntaxError("panic")
		}
	}()

	env := liquid.NewEnvironment()
	tags.RegisterStandardTags(env)
	env.SetErrorMode("strict")
	_, err = liquid.ParseTemplate(src, &liquid.TemplateOptions{Environment: env})
	return err
}
