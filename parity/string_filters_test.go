package parity

import "testing"

// Ruby's truncate counts characters; this counted bytes, so a multi-byte
// ellipsis or any non-ASCII text truncated far too aggressively.
func TestTruncateCountsCharacters(t *testing.T) {
	data := map[string]interface{}{
		"html":    "<p>Hello <b>world</b></p>",
		"unicode": "日本語のポッドキャスト番組",
		"emoji":   "Podcast 🎙 Show",
	}
	assertRender(t, `{{ html | truncate: 12, "…" }}`, "<p>Hello <b…", data)
	assertRender(t, `{{ unicode | truncate: 5 }}`, "日本...", data) // 5 chars incl. the "..."
	assertRender(t, `{{ emoji | truncate: 11, "" }}`, "Podcast 🎙 S", data)
	assertRender(t, `{{ "short" | truncate: 50 }}`, "short", data)
}

// Ruby's escape_once leaves existing entities alone; this re-escaped them, so
// a title already holding &amp; rendered as &amp;amp;.
func TestEscapeOnceDoesNotDoubleEscape(t *testing.T) {
	assertRender(t, `{{ "a&amp;b <c>" | escape_once }}`, "a&amp;b &lt;c&gt;", nil)
	assertRender(t, `{{ "Q & A" | escape_once }}`, "Q &amp; A", nil)
	assertRender(t, `{{ "&#39;quoted&#39; <b>" | escape_once }}`, "&#39;quoted&#39; &lt;b&gt;", nil)
	assertRender(t, `{{ "plain" | escape_once }}`, "plain", nil)
}

// Ruby's slice accepts a negative offset counting back from the end, and
// slices by character rather than byte.
func TestSliceHandlesNegativeOffsets(t *testing.T) {
	assertRender(t, `{{ "hello" | slice: -3, 2 }}`, "ll", nil)
	assertRender(t, `{{ "hello" | slice: 1, 3 }}`, "ell", nil)
	assertRender(t, `{{ "hello" | slice: -1 }}`, "o", nil)
	assertRender(t, `{{ "日本語" | slice: 1, 2 }}`, "本語", nil)
	assertRender(t, `{{ "hello" | slice: -99, 2 }}`, "", nil)
}
