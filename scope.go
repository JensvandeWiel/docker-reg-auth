package registry

import (
	"fmt"
	"strings"
)

// ScopeType represents the type of scope.
type ScopeType string

const (
	ScopeTypeRepository ScopeType = "repository"
	ScopeTypeRegistry   ScopeType = "registry"
	//ScopeTypeAll        ScopeType = "*"
)

// String returns the string representation of a ScopeType.
func (s ScopeType) String() string {
	return string(s)
}

// ParseName parses a name and returns a valid name. If the name is empty or ".", "catalog" is returned.
func ParseName(name string) string {
	if name == "" || name == "." {
		return "catalog"
	}
	return name
}

// ParseScopeType parses a string into a ScopeType. If the string is not a valid scope type, an error is returned.
func ParseScopeType(scopeType string) (ScopeType, error) {
	switch scopeType {
	case "repository":
		return ScopeTypeRepository, nil
	case "registry":
		return ScopeTypeRegistry, nil
		/*case "*":
		return ScopeTypeAll, nil
		*/
	}

	return "", fmt.Errorf("unknown scope type: %s", scopeType)
}

// Scope contains ScopeType, Name, and Actions for a scope.
type Scope struct {
	Type    ScopeType
	Name    string
	Actions ActionSet
}

// ParseScope parses a string into a Scope. If the string is not a valid scope, an error is returned.
func ParseScope(scope string) (*Scope, error) {
	parts := strings.Split(scope, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid scope: %s", scope)
	}

	scopeType, err := ParseScopeType(parts[0])
	if err != nil {
		return nil, err
	}

	parsedActions, err := ParseActions(strings.Split(parts[2], ","))
	if err != nil {
		return nil, err
	}

	return &Scope{
		Type:    scopeType,
		Name:    parts[1],
		Actions: parsedActions,
	}, nil
}
