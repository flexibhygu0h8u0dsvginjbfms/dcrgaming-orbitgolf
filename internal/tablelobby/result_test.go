package tablelobby

import "testing"

func TestResultSaysWhoWonFromThisSeat(t *testing.T) {
	score := [2]int{2, 5}
	won := ResultFor(1, 1, true, false, score, 100_000, 2, CheckVerified, false)
	if won.Outcome != "won" || won.Headline != "YOU WON" || won.Score != "5 : 2" {
		t.Fatalf("winner sees %+v", won)
	}
	if won.Money != "You take the pot: 0.00200000 DCR before fees." {
		t.Fatalf("winner's money line: %q", won.Money)
	}
	lost := ResultFor(0, 1, true, false, score, 100_000, 2, CheckVerified, false)
	if lost.Outcome != "lost" || lost.Headline != "YOU LOST" || lost.Score != "2 : 5" {
		t.Fatalf("loser sees %+v", lost)
	}
	if lost.Money != "Your 0.00100000 DCR stake goes to the winner." {
		t.Fatalf("loser's money line: %q", lost.Money)
	}
	draw := ResultFor(0, 0, false, false, [2]int{3, 3}, 100_000, 2, CheckVerified, false)
	if draw.Outcome != "draw" || draw.Money != "Your 0.00100000 DCR stake comes back to you." {
		t.Fatalf("draw shown as %+v", draw)
	}
	void := ResultFor(0, 0, true, true, score, 100_000, 2, CheckVerified, false)
	if void.Outcome != "void" || void.Headline != "MATCH VOID" {
		t.Fatalf("an abandoned match shown as %+v", void)
	}
}

func TestResultFollowsThePayout(t *testing.T) {
	states := func(r Result) string {
		out := ""
		for _, s := range r.Payout {
			out += s.State[:1]
		}
		return out
	}
	for _, c := range []struct {
		check string
		paid  bool
		want  string
	}{
		{CheckVerified, false, "ctt"},
		{CheckSpending, false, "ddc"},
		{CheckVerified, true, "ddc"},
		{CheckSpent, true, "ddd"},
	} {
		r := ResultFor(0, 0, true, false, [2]int{1, 0}, 100_000, 2, c.check, c.paid)
		if got := states(r); got != c.want {
			t.Errorf("stake %s, paid %v: payout %s, want %s (%s)", c.check, c.paid, got, c.want, r.PayoutNote)
		}
	}
}
