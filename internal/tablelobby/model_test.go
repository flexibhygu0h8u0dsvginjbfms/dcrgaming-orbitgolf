package tablelobby

import "testing"

func TestDCRRendersEveryAmountExactly(t *testing.T) {
	for atoms, want := range map[int64]string{
		0:               "0.00000000",
		100_000:         "0.00100000",
		123_456_789_012: "1234.56789012",
		-150_000_001:    "-1.50000001",
		21e14 - 1:       "20999999.99999999",
		21e14:           "21000000.00000000",
	} {
		if got := DCR(atoms); got != want {
			t.Errorf("DCR(%d) renders %q, want %q", atoms, got, want)
		}
	}
}
