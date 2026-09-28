package tablelobby

// Screen is the table lobby ready to draw: every word and state decided here,
// so the client only lays it out.
type Screen struct {
	Match        string       `json:"match"`
	Stage        string       `json:"stage"`
	Steps        []ScreenStep `json:"steps"`
	Seats        []ScreenSeat `json:"seats"`
	Terms        ScreenTerms  `json:"terms"`
	NextStep     string       `json:"nextStep"`
	NextDetail   string       `json:"nextDetail"`
	CanFund      bool         `json:"canFund"`
	Closed       bool         `json:"closed"`
	ClosedReason string       `json:"closedReason"`
	Error        string       `json:"error"`
	Stale        bool         `json:"stale"`
}

// ScreenStep is one box on the rail: done, current or todo.
type ScreenStep struct {
	Label string `json:"label"`
	State string `json:"state"`
}

// ScreenSeat is one player's card.
type ScreenSeat struct {
	Name     string          `json:"name"`
	You      bool            `json:"you"`
	Status   string          `json:"status"`
	Deposits []ScreenDeposit `json:"deposits"`
}

// ScreenDeposit is one dot on a seat card.
type ScreenDeposit struct {
	Name     string  `json:"name"`
	Short    string  `json:"short"`
	Progress float64 `json:"progress"`
	Complete bool    `json:"complete"`
	Wrong    bool    `json:"wrong"`
}

// ScreenTerms is what the table costs, as it is shown.
type ScreenTerms struct {
	BuyIn        string `json:"buyIn"`
	Bond         string `json:"bond"`
	PerSeat      string `json:"perSeat"`
	Pot          string `json:"pot"`
	RefundBlocks uint32 `json:"refundBlocks"`
	BondBlocks   uint32 `json:"bondBlocks"`
	Until        uint32 `json:"until"`
}

// Screen describes the view for drawing. Stale evidence shows nothing as
// complete, the same rule the view's own checks follow.
func (v View) Screen() Screen {
	s := Screen{
		Match: v.Match, Stage: v.Stage.Label(),
		NextStep: v.NextStep, NextDetail: v.NextDetail, CanFund: v.CanFund && !v.Stale,
		Closed: v.Stage == StageClosed, ClosedReason: v.ClosedReason, Error: v.Error, Stale: v.Stale,
		Terms: ScreenTerms{
			BuyIn: DCR(v.Terms.BuyInAtoms), Bond: DCR(v.Terms.BondAtoms),
			PerSeat: DCR(v.Terms.PerSeat()), Pot: DCR(v.Terms.Pot()),
			RefundBlocks: v.Terms.RefundBlocks, BondBlocks: v.Terms.BondBlocks, Until: v.Terms.Until,
		},
	}
	for _, step := range Steps {
		state := "todo"
		switch {
		case s.Closed:
			// A closed table will not form: nothing on its rail is done or
			// current, and the screen says why it closed.
		case step < v.Stage:
			state = "done"
		case step == v.Stage:
			state = "current"
		}
		s.Steps = append(s.Steps, ScreenStep{Label: step.Label(), State: state})
	}
	for _, seat := range v.Seats {
		card := ScreenSeat{Name: seat.Name, You: seat.You, Status: seat.StatusAt(v.Stale)}
		for _, purpose := range []string{PurposeSeatBond, PurposeStake} {
			d, _ := seat.Deposit(purpose)
			card.Deposits = append(card.Deposits, ScreenDeposit{
				Name: Name(purpose), Short: d.Short(), Progress: v.Progress(d),
				Complete: d.Complete() && !v.Stale, Wrong: d.Wrong(),
			})
		}
		s.Seats = append(s.Seats, card)
	}
	return s
}
