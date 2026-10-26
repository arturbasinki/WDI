package z02

import (
	"context"
	"errors"
	"testing"
	"time"
)

// Regression: feeder wychodził przez ctx.Done() BEZ close(jobs) — robotnicy
// czekali na `range jobs` w nieskończoność, wg.Wait() nie wracał, a serwer
// po DELETE zostawiał job „running” na zawsze (reprodukcja 30/30).
func TestSearchCancelReturns(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	type outcome struct {
		p   Progress
		err error
	}
	done := make(chan outcome, 1)
	go func() {
		p, err := Search(ctx, 1e8, nil)
		done <- outcome{p, err}
	}()
	select {
	case o := <-done:
		if !errors.Is(o.err, context.Canceled) {
			t.Fatalf("oczekiwano context.Canceled, jest: %v (progress: %+v)", o.err, o.p)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Search zawiesił się po anulowaniu — feeder musi close(jobs) przez defer")
	}
}
