package main

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/karamble/dcrgaming-orbitgolf/internal/session"
)

type fakeSettler struct {
	mu    sync.Mutex
	done  bool
	calls int
	fail  error
}

func (f *fakeSettler) Tables() []string { return []string{"t1"} }
func (f *fakeSettler) View(string) (session.View, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return session.View{Done: f.done}, true
}
func (f *fakeSettler) Settle(context.Context, string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	return f.fail
}
func (f *fakeSettler) count() int { f.mu.Lock(); defer f.mu.Unlock(); return f.calls }

// A finished match proposes its payout by itself, once, and asks the player to
// approve it; a proposal that failed is tried again.
func TestAFinishedMatchProposesItsPayout(t *testing.T) {
	a := &app{dir: t.TempDir(), practice: -1, jobs: make(chan error, 1)}
	f := &fakeSettler{}
	ctx := context.Background()
	a.settleFinished(ctx, f)
	a.settleWG.Wait()
	if f.count() != 0 {
		t.Fatal("a match still being played was settled")
	}
	f.mu.Lock()
	f.done, f.fail = true, errors.New("bridge busy")
	f.mu.Unlock()
	a.settleFinished(ctx, f)
	a.settleWG.Wait()
	if f.count() != 1 || a.status != approvePayout {
		t.Fatalf("finished match: %d proposals, status %q", f.count(), a.status)
	}
	f.mu.Lock()
	f.fail = nil
	f.mu.Unlock()
	a.settleFinished(ctx, f) // reads the failure, then tries again
	a.settleWG.Wait()
	if f.count() != 2 {
		t.Fatalf("a failed proposal was not retried: %d", f.count())
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		a.settleFinished(ctx, f)
		a.settleWG.Wait()
	}
	if f.count() != 2 {
		t.Fatalf("a proposed payout was proposed again: %d", f.count())
	}
	if a.status != approvePayout {
		t.Fatalf("status after a proposal: %q", a.status)
	}
}
