package tablelobby

import "fmt"

// Result is how a finished bridge match ended for one seat, ready to draw.
type Result struct {
	Outcome    string       `json:"outcome"`
	Headline   string       `json:"headline"`
	Score      string       `json:"score"`
	Money      string       `json:"money"`
	Payout     []ScreenStep `json:"payout"`
	PayoutNote string       `json:"payoutNote"`
}

// ResultFor describes a finished match from one seat: whether it won, what
// happens to its stake, and how far the payout has come. stakeCheck is this
// seat's own stake as last checked; paid is the payout having been seen.
func ResultFor(seat, winner uint32, won, void bool, score [2]int, buyIn int64, seats uint32, stakeCheck string, paid bool) Result {
	r := Result{Score: fmt.Sprintf("%d : %d", score[seat%2], score[(seat+1)%2])}
	stake := DCR(buyIn)
	switch {
	case void:
		r.Outcome, r.Headline, r.Money = "void", "MATCH VOID", fmt.Sprintf("Your %s DCR stake comes back to you.", stake)
	case !won:
		r.Outcome, r.Headline, r.Money = "draw", "DRAW", fmt.Sprintf("Your %s DCR stake comes back to you.", stake)
	case winner == seat:
		r.Outcome, r.Headline, r.Money = "won", "YOU WON", fmt.Sprintf("You take the pot: %s DCR before fees.", DCR(buyIn*int64(seats)))
	default:
		r.Outcome, r.Headline, r.Money = "lost", "YOU LOST", fmt.Sprintf("Your %s DCR stake goes to the winner.", stake)
	}
	states := [3]string{"current", "todo", "todo"}
	r.PayoutNote = "Approve the payout in dcrpulse. Both players sign it."
	switch {
	case stakeCheck == CheckSpent:
		states = [3]string{"done", "done", "done"}
		r.PayoutNote = "Paid out."
	case stakeCheck == CheckSpending || paid:
		states = [3]string{"done", "done", "current"}
		r.PayoutNote = "Payout broadcast, waiting for a block."
	}
	for i, label := range []string{"APPROVE", "BROADCAST", "PAID"} {
		r.Payout = append(r.Payout, ScreenStep{Label: label, State: states[i]})
	}
	return r
}
