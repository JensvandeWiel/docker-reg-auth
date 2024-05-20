package registry

import "errors"

type mockAuthenticator struct {
	mockUser, mockPass string
}

func (a *mockAuthenticator) Authenticate(user string, pass string) error {
	if user != a.mockUser || pass != a.mockPass {
		return errors.New("invalid credentials")
	}
	return nil
}

func newMockAuthenticator(user, pass string) *mockAuthenticator {
	return &mockAuthenticator{mockUser: user, mockPass: pass}
}

//TODO(JensvandeWiel) Add tests
