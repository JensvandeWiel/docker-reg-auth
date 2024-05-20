package registry

import "testing"

func TestParseAction(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		want    ActionType
		wantErr bool
	}{
		{"TestPull", "pull", ActionPull, false},
		{"TestPush", "push", ActionPush, false},
		{"TestAll", "*", ActionAll, false},
		{"TestCatalog", "catalog", ActionCatalog, false},
		{"TestAdmin", "admin", ActionAdmin, false},
		{"TestEmpty", "", "", true},
		{"TestUnknown", "unknown", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAction(tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseAction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionSet_Contains(t *testing.T) {
	set := ActionSet{ActionPull, ActionPush, ActionAll, ActionCatalog, ActionAdmin}
	tests := []struct {
		name string
		a    ActionType
		want bool
	}{
		{"TestPull", ActionPull, true},
		{"TestPush", ActionPush, true},
		{"TestAll", ActionAll, true},
		{"TestCatalog", ActionCatalog, true},
		{"TestAdmin", ActionAdmin, true},
		{"TestUnknown", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.Contains(tt.a); got != tt.want {
				t.Errorf("ActionSet.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionSet_ContainsAll(t *testing.T) {
	set := ActionSet{ActionPull, ActionPush, ActionCatalog, ActionAdmin}
	tests := []struct {
		name    string
		actions []ActionType
		want    bool
		length  int
	}{
		{"TestAllInSet", []ActionType{ActionPull, ActionPush, ActionCatalog, ActionAdmin}, true, 4},
		{"TestSomeInSet", []ActionType{ActionPull, ActionAll, ActionPush}, false, 2},
		{"TestSomeKnown", []ActionType{ActionPull, ActionPush}, true, 2},
		{"TestNoneInSet", []ActionType{ActionAll}, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.ContainsAll(tt.actions); got != tt.want {
				t.Errorf("ActionSet.ContainsAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestActionSet_ToStrings(t *testing.T) {
	set := ActionSet{ActionPull, ActionPush, ActionAll, ActionCatalog, ActionAdmin}
	want := []string{"pull", "push", "*", "catalog", "admin"}

	got := set.ToStrings()
	if len(got) != len(want) {
		t.Errorf("ActionSet.ToStrings() = %v, want %v", got, want)
	}

	for i, action := range got {
		if action != want[i] {
			t.Errorf("ActionSet.ToStrings() = %v, want %v", got, want)
		}
	}
}
