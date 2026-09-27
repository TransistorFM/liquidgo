package parity

import "testing"

// ToDate could not parse Ruby's Time#to_s / ActiveSupport TimeWithZone#to_s,
// so `| date:` returned its input unchanged for any timestamp produced by a
// Ruby backend.
func TestDateParsesRubyTimeFormats(t *testing.T) {
	data := map[string]interface{}{
		"rubyTime": "2024-03-01 09:07:00 -0500",
		"isoTime":  "2024-03-01T09:07:00-05:00",
		"plain":    "2024-03-01 09:07:00",
	}
	assertRender(t, `{{ rubyTime | date: "%b %d, %Y" }}`, "Mar 01, 2024", data)
	assertRender(t, `{{ isoTime | date: "%b %d, %Y" }}`, "Mar 01, 2024", data)
	assertRender(t, `{{ plain | date: "%b %d, %Y" }}`, "Mar 01, 2024", data)
	assertRender(t, `{{ "2024-12-25" | date: "%m/%d/%Y" }}`, "12/25/2024", data)
}
