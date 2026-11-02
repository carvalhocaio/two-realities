package level

import "testing"

func TestCampaignActs(t *testing.T) {
	levels, err := LoadAll()
	if err != nil {
		t.Fatal(err)
	}
	want := []int{1, 1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3}
	if len(levels) != len(want) {
		t.Fatalf("%d levels, want %d", len(levels), len(want))
	}
	for i, l := range levels {
		if l.Act != want[i] {
			t.Errorf("%s act = %d, want %d", l.Name, l.Act, want[i])
		}
	}
}
