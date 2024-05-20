package registry

import (
	"fmt"
	"github.com/labstack/echo/v4"
)

//TODO(JensvandeWiel) Implement real authorizer

// Authorizer is an interface for authorizing requests. It is used to check if a token request is allowed to perform certain actions.
type Authorizer interface {
	// Authorize authorizes the request and returns the allowed actions. If one of the actions is not allowed, an error is returned.
	Authorize(req *AuthorizationRequest) (ActionSet, error)
}

// AuthorizationRequest is a request for authorization.
type AuthorizationRequest struct {
	Account string
	Service string
	Type    ScopeType
	Name    string
	IP      string
	Actions ActionSet
}

func AuthorizationRequestFromContext(ctx echo.Context) (*AuthorizationRequest, error) {
	q := ctx.QueryParams()
	req := &AuthorizationRequest{}
	if account := q.Get("account"); account != "" {
		req.Account = account
	} else {
		return nil, fmt.Errorf("account is required")
	}

	if service := q.Get("service"); service != "" {
		req.Service = service
	} else {
		return nil, fmt.Errorf("service is required")
	}

	if scope := q.Get("scope"); scope != "" {
		scope, err := ParseScope(scope)
		if err != nil {
			return nil, err
		}
		req.Type = scope.Type
		req.Name = scope.Name
		req.Actions = scope.Actions
		if req.Account == "" {
			req.Account = scope.Name
		}
	} else {
		return nil, fmt.Errorf("scope is required")
	}

	return req, nil
}

// DummyAuthorizer is an authorizer that always authorizes the request.
type DummyAuthorizer struct {
}

func NewDummyAuthorizer() *DummyAuthorizer {
	return &DummyAuthorizer{}
}

// Authorize always returns the request actions.
func (a *DummyAuthorizer) Authorize(req *AuthorizationRequest) (ActionSet, error) {
	return req.Actions, nil
}
