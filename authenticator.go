package registry

// Authenticator is an interface for authenticating users for the registry server.
type Authenticator interface {
	// Authenticate authenticates the user with the given username and password. If the authentication fails, an error is returned. Depending on the implementation, the way you need to provide the password may differ.
	Authenticate(user string, pass string) error
}

// DummyAuthenticator is an authenticator that always succeeds.
type DummyAuthenticator struct {
}

func NewDummyAuthenticator() *DummyAuthenticator {
	return &DummyAuthenticator{}
}

func (a *DummyAuthenticator) Authenticate(user string, pass string) error {
	return nil
}
