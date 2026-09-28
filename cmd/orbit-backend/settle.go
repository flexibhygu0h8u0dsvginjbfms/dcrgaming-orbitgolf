package main

import (
	"context"
	"time"

	"github.com/karamble/dcrgaming-orbitgolf/internal/session"
)

// settler is the part of a game the payout watch needs.
type settler interface {
	Tables() []string
	View(sid string) (session.View, bool)
	Settle(ctx context.Context, sid string) error
}

type settleResult struct {
	sid string
	err error
}

// approvePayout is what a player reads once their match is over.
const approvePayout = "Match over. Approve the payout in dcrpulse."

// settleFinished proposes the payout of every finished match by itself, once
// per table and in the background: the bridge builds it and each operator
// approves it in dcrpulse, which takes as long as a person takes. A proposal
// that failed is tried again on a later poll.
func (a *app) settleFinished(ctx context.Context, g settler) {
	if a.settling == nil {
		a.settling = map[string]bool{}
		a.settled = make(chan settleResult, 16)
	}
	for drained := false; !drained; {
		select {
		case r := <-a.settled:
			if r.err != nil {
				delete(a.settling, r.sid)
				a.status = "Payout not proposed yet: " + r.err.Error()
			}
		default:
			drained = true
		}
	}
	for _, sid := range g.Tables() {
		v, ok := g.View(sid)
		if !ok || !v.Done || a.settling[sid] {
			continue
		}
		a.settling[sid] = true
		a.status = approvePayout
		a.settleWG.Add(1)
		go func() {
			defer a.settleWG.Done()
			call, cancel := context.WithTimeout(ctx, 2*time.Minute)
			defer cancel()
			a.settled <- settleResult{sid: sid, err: g.Settle(call, sid)}
		}()
	}
}
