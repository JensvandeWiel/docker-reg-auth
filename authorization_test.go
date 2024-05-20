package registry

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockAuthorizer struct {
	actions ActionSet
}

func (m *mockAuthorizer) Authorize(req *AuthorizationRequest) (ActionSet, error) {
	found := m.actions.ContainsAll(req.Actions)
	if !found {
		return nil, fmt.Errorf("one or more actions are not allowed")
	}

	return req.Actions, nil
}

func TestCanParseAuthorizationRequest(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.QueryParams().Set("account", "test")
	c.QueryParams().Set("service", "test")
	c.QueryParams().Set("scope", "repository:foo/bar:pull,push")
	want := &AuthorizationRequest{
		Account: "test",
		Service: "test",
		Type:    "repository",
		Name:    "foo/bar",
		Actions: ActionSet{ActionPull, ActionPush},
	}

	got, err := AuthorizationRequestFromContext(c)
	if err != nil {
		t.Errorf("AuthorizationRequestFromContext() error = %v", err)
		return
	}

	if got.Account != want.Account {
		t.Errorf("AuthorizationRequestFromContext() got.Account = %v, want %v", got.Account, want.Account)
	}

	if got.Service != want.Service {
		t.Errorf("AuthorizationRequestFromContext() got.Service = %v, want %v", got.Service, want.Service)
	}

	if got.Type != want.Type {
		t.Errorf("AuthorizationRequestFromContext() got.Type = %v, want %v", got.Type, want.Type)
	}

	if got.Name != want.Name {
		t.Errorf("AuthorizationRequestFromContext() got.Name = %v, want %v", got.Name, want.Name)
	}

	if len(got.Actions) != len(want.Actions) {
		t.Errorf("AuthorizationRequestFromContext() got.Actions = %v, want %v", got.Actions, want.Actions)
	}

	if !got.Actions.ContainsAll(want.Actions) {
		t.Errorf("AuthorizationRequestFromContext() got.Actions = %v, want %v", got.Actions, want.Actions)
	}
}

func TestAuth_Authorize(t *testing.T) {
	tests := []struct {
		name    string
		has     ActionSet
		req     *AuthorizationRequest
		wantErr bool
	}{
		{
			name: "HasAllActions",
			has:  ActionSet{ActionPull, ActionPush, ActionCatalog, ActionAll, ActionAdmin},
			req: &AuthorizationRequest{
				Actions: ActionSet{ActionPull, ActionPush},
			},
			wantErr: false,
		},
		{
			name: "HasSomeActions",
			has:  ActionSet{ActionPull, ActionPush},
			req: &AuthorizationRequest{
				Actions: ActionSet{ActionPull, ActionPush, ActionCatalog},
			},
			wantErr: true,
		},
		{
			name: "HasNoActions",
			has:  ActionSet{},
			req: &AuthorizationRequest{
				Actions: ActionSet{ActionPull, ActionPush, ActionCatalog, ActionAll, ActionAdmin},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			mock := &mockAuthorizer{
				actions: tt.has,
			}

			got, err := mock.Authorize(tt.req)

			if err != nil {
				if !tt.wantErr {
					t.Errorf("Authorizer.Authorize() error = %v, wantErr %v, hasActions: %v, gotActions: %v", err, tt.wantErr, tt.has, got)
				}

				return
			}
			assert.Equal(t, got, tt.req.Actions)

		})

	}
}
