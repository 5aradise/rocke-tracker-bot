package md

import "strings"

var mdEscaper = strings.NewReplacer(
	"\\", "\\\\",
	"`", "\\`",
	"*", "\\*",
	"_", "\\_",
	"{", "\\{",
	"}", "\\}",
	"[", "\\[",
	"]", "\\]",
	"(", "\\(",
	")", "\\)",
	"#", "\\#",
	"+", "\\+",
	"-", "\\-",
	".", "\\.",
	"!", "\\!",
)

func Escape(s string) string {
	return mdEscaper.Replace(s)
}
