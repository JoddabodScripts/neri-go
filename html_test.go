package nerimity

import "testing"

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{`<script>`, `&lt;script&gt;`},
		{`a & b`, `a &amp; b`},
		{`"quoted"`, `&quot;quoted&quot;`},
		{`it's`, `it&#39;s`},
		{`<img src="x" onerror='alert(1)'>`, `&lt;img src=&quot;x&quot; onerror=&#39;alert(1)&#39;&gt;`},
		{`plain text`, `plain text`},
		{``, ``},
	}
	for _, tt := range tests {
		got := EscapeHTML(tt.in)
		if got != tt.want {
			t.Errorf("EscapeHTML(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEscapeHTMLAmpersandOrdering(t *testing.T) {
	// & must be escaped first, otherwise escaping "<" to "&lt;" and then "&"
	// to "&amp;" would double-escape it into "&amp;lt;".
	got := EscapeHTML("<")
	if got != "&lt;" {
		t.Errorf("EscapeHTML(%q) = %q, want %q (double-escaping bug)", "<", got, "&lt;")
	}
}
