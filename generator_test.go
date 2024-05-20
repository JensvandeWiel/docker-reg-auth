package registry

import (
	"github.com/distribution/distribution/registry/auth/token"
	"testing"
)

func TestDefaultTokenGenerator_GenerateToken(t *testing.T) {
	g := NewDefaultTokenGenerator(".devcerts/RootCA.crt", ".devcerts/RootCA.key")
	tokOpt := &TokenOptions{
		Issuer:    "paca-node",
		Audience:  "registry",
		ExpiresIn: 3600,
	}
	rawToken, err := g.GenerateToken(&AuthorizationRequest{
		Account: "jens",
		Service: "registry",
		Type:    ScopeTypeRegistry,
		Name:    "catalog",
		Actions: ActionSet{ActionAll},
	}, ActionSet{ActionAll}, tokOpt)
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
		return
	}

	if rawToken == nil {
		t.Errorf("GenerateToken() token is nil")
	}

	tok, err := token.NewToken(rawToken.Token)
	if err != nil {
		t.Errorf("NewToken() error = %v", err)
		return
	}

	if tok.Header.SigningAlg != "RS256" {
		t.Errorf("GenerateToken() SigningAlgorithm = %v, want RS256", tok.Header.SigningAlg)
	}

	if tok.Claims.Issuer != tokOpt.Issuer {
		t.Errorf("GenerateToken() Issuer = %v, want paca-node", tok.Claims.Issuer)
	}

	if tok.Claims.Subject != "jens" {
		t.Errorf("GenerateToken() Subject = %v, want jens", tok.Claims.Subject)
	}

	if tok.Claims.Audience != tokOpt.Audience {
		t.Errorf("GenerateToken() Audience = %v, want Authentication", tok.Claims.Audience)
	}

	if tok.Claims.Expiration-tok.Claims.IssuedAt != 3600 {
		t.Errorf("GenerateToken() Expiration = %v, want 3600", tok.Claims.Expiration)
	}

	if tok.Claims.Access[0].Type != "registry" {
		t.Errorf("GenerateToken() Type = %v, want registry", tok.Claims.Access[0].Type)
	}

	if tok.Claims.Access[0].Name != "catalog" {
		t.Errorf("GenerateToken() Name = %v, want catalog", tok.Claims.Access[0].Name)
	}

	if len(tok.Claims.Access[0].Actions) != 1 {
		t.Errorf("GenerateToken() Actions = %v, want 1", len(tok.Claims.Access[0].Actions))
	}

	if tok.Claims.Access[0].Actions[0] != "*" {
		t.Errorf("GenerateToken() Actions = %v, want *", tok.Claims.Access[0].Actions[0])
	}
}
