package registry

import "fmt"

type ActionType string

type ActionSet []ActionType

const (
	ActionPull    ActionType = "pull"
	ActionPush    ActionType = "push"
	ActionCatalog ActionType = "catalog"
	ActionAll     ActionType = "*"
	ActionAdmin   ActionType = "admin"
)

// ParseAction parses a string into an ActionType. If the string is not a valid action, an error is returned.
func ParseAction(action string) (ActionType, error) {
	switch action {
	case "pull":
		return ActionPull, nil
	case "push":
		return ActionPush, nil
	case "*":
		return ActionAll, nil
	case "catalog":
		return ActionCatalog, nil
	case "admin":
		return ActionAdmin, nil
	}

	return "", fmt.Errorf("unknown action: %s", action)
}

// ParseActions parses a slice of strings into an ActionSet. If any of the strings are not valid actions, an error is returned.
func ParseActions(actions []string) (ActionSet, error) {
	actionSet := make(ActionSet, len(actions))
	for i, action := range actions {
		actionType, err := ParseAction(action)
		if err != nil {
			return nil, err
		}
		actionSet[i] = actionType
	}
	return actionSet, nil
}

// String returns the string representation of an ActionType.
func (a ActionType) String() string {
	return string(a)
}

// Contains checks if the ActionSet contains the given action.
func (a ActionSet) Contains(action ActionType) bool {
	for _, a := range a {
		if a == action {
			return true
		}
	}
	return false
}

// ContainsAll checks if the ActionSet contains all the actions in the given slice, and returns a slice of the actions that are in the slice.
func (a ActionSet) ContainsAll(actions []ActionType) bool {
	for _, action := range actions {
		if !a.Contains(action) {
			return false
		}
	}

	return true
}

// ToStrings returns a slice of strings representing the actions in the ActionSet.
func (a ActionSet) ToStrings() []string {
	actions := make([]string, len(a))
	for i, action := range a {
		actions[i] = action.String()
	}
	return actions
}
