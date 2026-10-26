package z06

import (
	"strings"
	"testing"
)

// Regression: parser był rekurencyjny bez limitu — 10^6 znaków „(”
// zabijało proces (fatal error: stack overflow, poza odzyskiwaniem net/http).
func TestParserDepthLimit(t *testing.T) {
	deep := strings.Repeat("(", 5000) + "x" + strings.Repeat(")", 5000)
	if _, err := PlotFun(deep, -1, 1, 601); err == nil {
		t.Fatal("spodziewano się błędu limitu głębokości dla 5000 zagnieżdżeń")
	}
	ok := strings.Repeat("(", 40) + "x+1" + strings.Repeat(")", 40)
	if _, err := PlotFun(ok, -1, 1, 601); err != nil {
		t.Fatalf("40 zagnieżdżeń powinno się parsować: %v", err)
	}
}
