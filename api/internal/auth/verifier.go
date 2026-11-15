// Package auth verifies Firebase ID tokens (RS256 JWTs) against a JWKS. The
// issuer, audience and key URL are configuration, so another OIDC provider can
// replace Firebase later. Only the token's sub is exposed; email is never kept.
package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	firebaseJWKSURL = "https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com"
	leeway          = time.Minute
)

var (
	// ErrInvalidToken covers every rejected token: bad signature, claims or kid.
	ErrInvalidToken = errors.New("invalid token")
	// ErrEmailNotVerified means the token carries an unverified email.
	ErrEmailNotVerified = errors.New("email not verified")
	// ErrKeysUnavailable means the signing keys could not be fetched.
	ErrKeysUnavailable = errors.New("signing keys unavailable")
)

// Claims is the part of a verified token the API uses.
type Claims struct {
	Sub string
}

// Config configures a Verifier. Zero Client and Now get defaults.
type Config struct {
	Issuer   string
	Audience string
	JWKSURL  string
	Client   *http.Client
	Now      func() time.Time
}

// Verifier checks ID tokens. It is safe for concurrent use.
type Verifier struct {
	issuer, audience string
	keys             *keySet
	now              func() time.Time
}

// NewFirebase returns a Verifier for a Firebase project.
func NewFirebase(projectID string) *Verifier {
	return New(Config{
		Issuer:   "https://securetoken.google.com/" + projectID,
		Audience: projectID,
		JWKSURL:  firebaseJWKSURL,
	})
}

func New(cfg Config) *Verifier {
	if cfg.Client == nil {
		cfg.Client = &http.Client{Timeout: 5 * time.Second}
	}
	if cfg.Now == nil {
		cfg.Now = time.Now
	}
	return &Verifier{
		issuer:   cfg.Issuer,
		audience: cfg.Audience,
		keys:     &keySet{url: cfg.JWKSURL, client: cfg.Client, now: cfg.Now},
		now:      cfg.Now,
	}
}

type tokenClaims struct {
	jwt.RegisteredClaims
	AuthTime      *jwt.NumericDate `json:"auth_time"`
	Email         string           `json:"email"`
	EmailVerified bool             `json:"email_verified"`
}

// Verify checks the signature (RS256 only) and the claims. Errors wrap
// ErrInvalidToken, ErrEmailNotVerified or ErrKeysUnavailable.
func (v *Verifier) Verify(ctx context.Context, raw string) (Claims, error) {
	var c tokenClaims
	_, err := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(v.issuer),
		jwt.WithAudience(v.audience),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithLeeway(leeway),
		jwt.WithTimeFunc(v.now),
	).ParseWithClaims(raw, &c, func(t *jwt.Token) (any, error) {
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			return nil, ErrInvalidToken
		}
		return v.keys.get(ctx, kid)
	})
	if err != nil {
		if errors.Is(err, ErrKeysUnavailable) {
			return Claims{}, ErrKeysUnavailable
		}
		return Claims{}, ErrInvalidToken
	}
	if c.Subject == "" {
		return Claims{}, ErrInvalidToken
	}
	if c.AuthTime != nil && c.AuthTime.Time.After(v.now().Add(leeway)) {
		return Claims{}, ErrInvalidToken
	}
	if c.Email != "" && !c.EmailVerified {
		return Claims{}, ErrEmailNotVerified
	}
	return Claims{Sub: c.Subject}, nil
}
