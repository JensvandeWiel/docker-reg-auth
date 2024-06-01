package registry

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
)

//TODO(JensvandeWiel) Implement real authorizer

// Authorizer is an interface for authorizing requests. It is used to check if a token request is allowed to perform certain actions.
type Authorizer interface {
	// Authorize authorizes the request and returns the allowed actions. If one of the actions is not allowed, an error is returned. Context is used to pass extra information to the authorizer, like the request context.
	Authorize(ctx context.Context, req *AuthorizationRequest) (ActionSet, error)
}

type AccessType string

const (
	AccessTypeOnline  AccessType = "online"
	AccessTypeOffline AccessType = "offline"
)

// ParseAccessType parses a string into an AccessType. If the string is not a valid access type, AccessTypeOnline is returned.
func parseAccessType(accessType string) AccessType {
	switch accessType {
	case "offline":
		return AccessTypeOffline
	default:
		return AccessTypeOnline
	}
}

// AuthorizationRequest is a request for authorization.
type AuthorizationRequest struct {
	Account    string
	Service    string
	Type       ScopeType
	Name       string
	IP         string
	Actions    ActionSet
	ClientId   string
	AccessType AccessType
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

	if clientId := q.Get("client_id"); clientId != "" {
		req.ClientId = clientId
	} else {
		return nil, fmt.Errorf("client_id is required")
	}

	req.AccessType = parseAccessType(q.Get("access_type"))

	if req.AccessType == AccessTypeOffline {
		return nil, fmt.Errorf("offline access type is not yet supported")
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
func (a *DummyAuthorizer) Authorize(ctx context.Context, req *AuthorizationRequest) (ActionSet, error) {
	return req.Actions, nil
}
