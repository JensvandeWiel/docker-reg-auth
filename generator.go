package registry

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/distribution/distribution/registry/auth/token"
	"github.com/docker/libtrust"
	"math/rand"
	"strings"
	"time"
)

const (
	SignAuth = "AUTH"
)

// TokenOptions contains the options for generating a token.
type TokenOptions struct {
	ExpiresIn int64
	Issuer    string
	Audience  string
}

// Token represents a token and an access token.
type Token struct {
	Token        string `json:"token"`
	AccessToken  string `json:"access_token"`
	IssuedAt     int64  `json:"issued_at"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// TokenGenerator is an interface for generating tokens.
type TokenGenerator interface {
	GenerateToken(req *AuthorizationRequest, actions ActionSet, options *TokenOptions) (*Token, error)
}

// DefaultTokenGenerator is a default implementation of the TokenGenerator interface.
type DefaultTokenGenerator struct {
	privKey libtrust.PrivateKey
	pubKey  libtrust.PublicKey
}

// NewDefaultTokenGenerator creates a new DefaultTokenGenerator.
func NewDefaultTokenGenerator(certPath, keyPath string) (*DefaultTokenGenerator, error) {
	pubKey, privKey, err := LoadCertificateAndKey(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	return &DefaultTokenGenerator{
		privKey: privKey,
		pubKey:  pubKey,
	}, nil
}

// GenerateToken generates a token for the given request and actions with the given options.
func (g *DefaultTokenGenerator) GenerateToken(req *AuthorizationRequest, actions ActionSet, tokenOptions *TokenOptions) (*Token, error) {
	if g.privKey == nil || g.pubKey == nil {
		return nil, fmt.Errorf("private or public key is nil")
	}
	_, algo, err := g.privKey.Sign(strings.NewReader(SignAuth), 0)
	if err != nil {
		return nil, err
	}

	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	if tokenOptions == nil {
		return nil, fmt.Errorf("token options are nil")
	}

	header := token.Header{
		Type:       "JWT",
		SigningAlg: algo,
		KeyID:      g.pubKey.KeyID(),
	}

	headerJson, err := json.Marshal(header)

	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	if req.Service != tokenOptions.Audience {
		return nil, fmt.Errorf("service does not match audience")
	}

	// Check if the all requested actions are allowed
	if !actions.ContainsAll(req.Actions) {
		return nil, fmt.Errorf("request actions do not match allowed actions")
	}

	claim := token.ClaimSet{
		Issuer:     tokenOptions.Issuer,
		Subject:    req.Account,
		Audience:   tokenOptions.Audience,
		Expiration: now + tokenOptions.ExpiresIn,
		NotBefore:  now - 10,
		IssuedAt:   now,
		JWTID:      fmt.Sprintf("%d", rand.Int63()),
		Access:     make([]*token.ResourceActions, 0), //[]*token.ResourceActions{}
	}

	claim.Access = append(claim.Access, &token.ResourceActions{
		Type:    req.Type.String(),
		Name:    req.Name,
		Actions: actions.ToStrings(),
	})

	claimJson, err := json.Marshal(claim)

	payload := fmt.Sprintf("%s%s%s", encodeBase64(headerJson), token.TokenSeparator, encodeBase64(claimJson))

	sig, sigAlgo, err := g.privKey.Sign(strings.NewReader(payload), 0)
	if err != nil && sigAlgo != algo {
		return nil, err
	}

	tok := fmt.Sprintf("%s%s%s", payload, token.TokenSeparator, encodeBase64(sig))

	return &Token{
		Token:       tok,
		AccessToken: tok,
		IssuedAt:    now,
		ExpiresIn:   tokenOptions.ExpiresIn,
	}, nil
}

func encodeBase64(data []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(data), "=")
}
