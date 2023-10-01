package astTraversal

import "strings"

func FormatDoc(doc string) string {
	return strings.TrimSpace(strings.TrimPrefix(doc, "//"))
}
