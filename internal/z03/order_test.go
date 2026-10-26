package z03

import (
	"reflect"
	"testing"
)

// Regression: SieveMap zwracał liczby w kolejności iteracji mapy
// (niedeterministycznie), podczas gdy SieveSlice/SieveBit — rosnąco.
func TestSieveMapOrdered(t *testing.T) {
	if got, want := SieveMap(200), SieveSlice(200); !reflect.DeepEqual(got, want) {
		t.Fatalf("SieveMap nieposortowany:\n%v\n!=\n%v", got, want)
	}
}
