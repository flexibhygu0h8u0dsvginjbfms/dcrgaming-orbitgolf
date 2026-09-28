package tablelobby

import "testing"

func formingView(stage Stage) View {
	return View{
		Match: "abcdef01",
		Stage: stage,
		Terms: Terms{BuyInAtoms: 100_000, BondAtoms: 1_000_000, Seats: 2, RefundBlocks: 288, BondBlocks: 288, Until: 1000},
		Seats: []Seat{
			{Number: 0, You: true, Name: "You", Deposits: []Deposit{
				{Purpose: PurposeSeatBond, Check: CheckConfirming, Local: true, Confirmations: 1, RequiredConfirmations: 2},
			}},
			{Number: 1, Name: "Seat 1", Deposits: []Deposit{
				{Purpose: PurposeSeatBond, Check: CheckVerified, Local: true, Confirmations: 2, RequiredConfirmations: 2},
			}},
		},
		NextStep: "ADMISSION BOND", NextDetail: "Table formation is managed by dcrpulse.",
	}
}

func TestScreenRailFollowsTheStage(t *testing.T) {
	s := formingView(StageRoster).Screen()
	want := []string{"done", "done", "current", "todo", "todo", "todo"}
	if len(s.Steps) != len(want) {
		t.Fatalf("rail has %d steps, want %d", len(s.Steps), len(want))
	}
	for i, step := range s.Steps {
		if step.State != want[i] {
			t.Errorf("step %d (%s) is %s, want %s", i, step.Label, step.State, want[i])
		}
	}
	if s.Steps[0].Label != "INVITATION" || s.Steps[5].Label != "READY" {
		t.Fatalf("rail labels: %+v", s.Steps)
	}
	closed := formingView(StageClosed)
	closed.ClosedReason = "admission expired"
	cs := closed.Screen()
	if !cs.Closed || cs.ClosedReason != "admission expired" {
		t.Fatalf("a closed table does not say so: %+v", cs)
	}
	for _, step := range cs.Steps {
		if step.State != "todo" {
			t.Fatalf("a closed table still shows progress on its rail: %+v", cs.Steps)
		}
	}
}

func TestScreenSeatsShowOwnSeatAndProgress(t *testing.T) {
	s := formingView(StageAdmission).Screen()
	if !s.Seats[0].You || s.Seats[1].You {
		t.Fatalf("own seat not marked: %+v", s.Seats)
	}
	mine, theirs := s.Seats[0].Deposits[0], s.Seats[1].Deposits[0]
	if mine.Name != "ENTRY" || mine.Short != "1 of 2" || mine.Progress != 0.5 || mine.Complete {
		t.Fatalf("own bond shown as %+v", mine)
	}
	if !theirs.Complete || theirs.Progress != 1 {
		t.Fatalf("other bond shown as %+v", theirs)
	}
	if stake := s.Seats[0].Deposits[1]; stake.Name != "STAKE" || stake.Short != "not posted" || stake.Complete {
		t.Fatalf("a stake nobody posted shown as %+v", stake)
	}
	if s.Terms.BuyIn != "0.00100000" || s.Terms.PerSeat != "0.01100000" || s.Terms.Pot != "0.00200000" {
		t.Fatalf("terms shown as %+v", s.Terms)
	}
}

func TestScreenOnStaleEvidenceShowsNothingDone(t *testing.T) {
	v := formingView(StageStake)
	v.CanFund = true
	v.Stale = true
	s := v.Screen()
	if s.CanFund {
		t.Fatal("stale evidence offers funding")
	}
	for _, seat := range s.Seats {
		if seat.Status != "recheck required" {
			t.Fatalf("stale seat says %q", seat.Status)
		}
		for _, d := range seat.Deposits {
			if d.Complete || d.Progress != 0 {
				t.Fatalf("stale deposit shown as %+v", d)
			}
		}
	}
}
