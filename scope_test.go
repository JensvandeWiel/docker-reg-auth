package registry

import "testing"

func TestParseScope(t *testing.T) {
	tests := []struct {
		rawScope string
		want     Scope
	}{
		{
			"repository:foo/bar:pull,push",
			Scope{
				Type:    ScopeTypeRepository,
				Name:    "foo/bar",
				Actions: ActionSet{ActionPull, ActionPush},
			},
		},
		{
			"repository:foo/baz/bar:pull",
			Scope{
				Type:    ScopeTypeRepository,
				Name:    "foo/baz/bar",
				Actions: ActionSet{ActionPull},
			},
		},
		{
			"repository:foo/bar:push",
			Scope{
				Type:    ScopeTypeRepository,
				Name:    "foo/bar",
				Actions: ActionSet{ActionPush},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.rawScope, func(t *testing.T) {
			got, err := ParseScope(tt.rawScope)
			if err != nil {
				t.Errorf("ParseScope() error = %v", err)
				return
			}
			if got.Type != tt.want.Type {
				t.Errorf("ParseScope() got.Type = %v, want %v", got.Type, tt.want.Type)
			}
			if got.Name != tt.want.Name {
				t.Errorf("ParseScope() got.Name = %v, want %v", got.Name, tt.want.Name)
			}
			if len(got.Actions) != len(tt.want.Actions) {
				t.Errorf("ParseScope() got.Actions = %v, want %v", got.Actions, tt.want.Actions)
			}
			for i, action := range got.Actions {
				if action != tt.want.Actions[i] {
					t.Errorf("ParseScope() got.Actions[%d] = %v, want %v", i, action, tt.want.Actions[i])
				}
			}
		})
	}
}

func TestFailParseScope(t *testing.T) {
	tests := []struct {
		rawScope string
	}{
		{
			"non:foo/bar:pull,push,",
		},
		{
			"repository:foo/baz:bar:pull",
		},
		{
			"repository:foo/bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.rawScope, func(t *testing.T) {
			_, err := ParseScope(tt.rawScope)
			if err == nil {
				t.Errorf("ParseScope() error = %v, want error", err)
			}
		})
	}
}
