package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	testProject = "demo-project"
	testIssuer  = "https://securetoken.google.com/" + testProject
)

type fakeJWKS struct {
	mu     sync.Mutex
	keys   map[string]*rsa.PublicKey
	hits   int
	maxAge int
	srv    *httptest.Server
}

func newFakeJWKS(t *testing.T, maxAge int) *fakeJWKS {
	t.Helper()
	f := &fakeJWKS{keys: map[string]*rsa.PublicKey{}, maxAge: maxAge}
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.hits++
		var list []map[string]string
		for kid, k := range f.keys {
			list = append(list, map[string]string{
				"kid": kid, "kty": "RSA", "alg": "RS256", "use": "sig",
				"n": base64.RawURLEncoding.EncodeToString(k.N.Bytes()),
				"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(k.E)).Bytes()),
			})
		}
		w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(f.maxAge)+", must-revalidate")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": list})
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeJWKS) add(t *testing.T, kid string) *rsa.PrivateKey {
	t.Helper()
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	f.mu.Lock()
	f.keys[kid] = &k.PublicKey
	f.mu.Unlock()
	return k
}

func (f *fakeJWKS) hitCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.hits
}

type clock struct{ t time.Time }

func (c *clock) now() time.Time { return c.t }

func newVerifier(f *fakeJWKS, c *clock) *Verifier {
	return New(Config{Issuer: testIssuer, Audience: testProject, JWKSURL: f.srv.URL, Now: c.now})
}

func baseClaims(c *clock) jwt.MapClaims {
	return jwt.MapClaims{
		"iss": testIssuer, "aud": testProject, "sub": "user-1",
		"iat": c.t.Add(-time.Minute).Unix(), "auth_time": c.t.Add(-time.Minute).Unix(),
		"exp": c.t.Add(time.Hour).Unix(),
	}
}

func sign(t *testing.T, method jwt.SigningMethod, kid string, key any, claims jwt.MapClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(method, claims)
	if kid != "" {
		tok.Header["kid"] = kid
	}
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestVerify(t *testing.T) {
	c := &clock{t: time.Now().UTC().Truncate(time.Second)}
	f := newFakeJWKS(t, 3600)
	priv := f.add(t, "k1")
	v := newVerifier(f, c)
	ctx := context.Background()

	good, err := v.Verify(ctx, sign(t, jwt.SigningMethodRS256, "k1", priv, baseClaims(c)))
	if err != nil || good.Sub != "user-1" {
		t.Fatalf("valid token: %+v %v", good, err)
	}

	mod := func(f func(jwt.MapClaims)) jwt.MapClaims {
		m := baseClaims(c)
		f(m)
		return m
	}
	hmacTok := sign(t, jwt.SigningMethodHS256, "k1", f.keys["k1"].N.Bytes(), baseClaims(c))
	noneTok := sign(t, jwt.SigningMethodNone, "k1", jwt.UnsafeAllowNoneSignatureType, baseClaims(c))
	otherKey, _ := rsa.GenerateKey(rand.Reader, 2048)

	cases := []struct {
		name string
		tok  string
	}{
		{"wrong audience", sign(t, jwt.SigningMethodRS256, "k1", priv, mod(func(m jwt.MapClaims) { m["aud"] = "other" }))},
		{"wrong issuer", sign(t, jwt.SigningMethodRS256, "k1", priv, mod(func(m jwt.MapClaims) { m["iss"] = "https://evil.example" }))},
		{"expired", sign(t, jwt.SigningMethodRS256, "k1", priv, mod(func(m jwt.MapClaims) { m["exp"] = c.t.Add(-2 * time.Minute).Unix() }))},
		{"missing exp", sign(t, jwt.SigningMethodRS256, "k1", priv, mod(func(m jwt.MapClaims) { delete(m, "exp") }))},
		{"iat in the future", sign(t, jwt.SigningMethodRS256, "k1", priv, mod(func(m jwt.MapClaims) { m["iat"] = c.t.Add(10 * time.Minute).Unix() }))},
		{"alg none", noneTok},
		{"HS256 with the public key as secret", hmacTok},
		{"unknown kid", sign(t, jwt.SigningMethodRS256, "nope", priv, baseClaims(c))},
		{"missing kid", sign(t, jwt.SigningMethodRS256, "", priv, baseClaims(c))},
		{"signed by another key", sign(t, jwt.SigningMethodRS256, "k1", otherKey, baseClaims(c))},
		{"empty sub", sign(t, jwt.SigningMethodRS256, "k1", priv, mod(func(m jwt.MapClaims) { m["sub"] = "" }))},
		{"missing sub", sign(t, jwt.SigningMethodRS256, "k1", priv, mod(func(m jwt.MapClaims) { delete(m, "sub") }))},
		{"auth_time in the future", sign(t, jwt.SigningMethodRS256, "k1", priv, mod(func(m jwt.MapClaims) { m["auth_time"] = c.t.Add(10 * time.Minute).Unix() }))},
		{"garbage", "not.a.jwt"},
		{"empty", ""},
	}
	for _, tc := range cases {
		if _, err := v.Verify(ctx, tc.tok); !errors.Is(err, ErrInvalidToken) {
			t.Errorf("%s: err = %v, want ErrInvalidToken", tc.name, err)
		}
	}

	// Within the one-minute leeway an iat slightly ahead of the clock passes.
	skew := sign(t, jwt.SigningMethodRS256, "k1", priv, mod(func(m jwt.MapClaims) { m["iat"] = c.t.Add(30 * time.Second).Unix() }))
	if _, err := v.Verify(ctx, skew); err != nil {
		t.Errorf("iat within leeway rejected: %v", err)
	}
}

func TestVerifyEmailVerification(t *testing.T) {
	c := &clock{t: time.Now().UTC().Truncate(time.Second)}
	f := newFakeJWKS(t, 3600)
	priv := f.add(t, "k1")
	v := newVerifier(f, c)
	ctx := context.Background()

	m := baseClaims(c)
	m["email"], m["email_verified"] = "a@example.com", false
	if _, err := v.Verify(ctx, sign(t, jwt.SigningMethodRS256, "k1", priv, m)); !errors.Is(err, ErrEmailNotVerified) {
		t.Errorf("unverified email: err = %v, want ErrEmailNotVerified", err)
	}
	m["email_verified"] = true
	if cl, err := v.Verify(ctx, sign(t, jwt.SigningMethodRS256, "k1", priv, m)); err != nil || cl.Sub != "user-1" {
		t.Errorf("verified email: %+v %v", cl, err)
	}
}

func TestKeyCacheFollowsMaxAge(t *testing.T) {
	c := &clock{t: time.Now().UTC().Truncate(time.Second)}
	f := newFakeJWKS(t, 600)
	priv := f.add(t, "k1")
	v := newVerifier(f, c)
	ctx := context.Background()
	tok := func() string { return sign(t, jwt.SigningMethodRS256, "k1", priv, baseClaims(c)) }

	for i := 0; i < 3; i++ {
		if _, err := v.Verify(ctx, tok()); err != nil {
			t.Fatal(err)
		}
	}
	if h := f.hitCount(); h != 1 {
		t.Fatalf("hits within max-age = %d, want 1", h)
	}
	c.t = c.t.Add(601 * time.Second)
	if _, err := v.Verify(ctx, tok()); err != nil {
		t.Fatal(err)
	}
	if h := f.hitCount(); h != 2 {
		t.Fatalf("hits after max-age = %d, want 2", h)
	}
}

func TestUnknownKidRefetchLimited(t *testing.T) {
	c := &clock{t: time.Now().UTC().Truncate(time.Second)}
	f := newFakeJWKS(t, 3600)
	k1 := f.add(t, "k1")
	v := newVerifier(f, c)
	ctx := context.Background()

	if _, err := v.Verify(ctx, sign(t, jwt.SigningMethodRS256, "k1", k1, baseClaims(c))); err != nil {
		t.Fatal(err)
	}
	base := f.hitCount()

	// A new key appears at the endpoint; the first token with its kid is
	// within the minute since the last fetch, so it is rejected without a fetch.
	k2 := f.add(t, "k2")
	tok2 := sign(t, jwt.SigningMethodRS256, "k2", k2, baseClaims(c))
	if _, err := v.Verify(ctx, tok2); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("k2 inside refetch window: %v", err)
	}
	if h := f.hitCount(); h != base {
		t.Fatalf("refetch inside the window: hits %d -> %d", base, h)
	}
	// After a minute one refetch picks up k2; further unknown kids do not refetch again.
	c.t = c.t.Add(61 * time.Second)
	if _, err := v.Verify(ctx, sign(t, jwt.SigningMethodRS256, "k2", k2, baseClaims(c))); err != nil {
		t.Fatalf("k2 after window: %v", err)
	}
	if h := f.hitCount(); h != base+1 {
		t.Fatalf("hits after refetch = %d, want %d", h, base+1)
	}
	for i := 0; i < 5; i++ {
		_, _ = v.Verify(ctx, sign(t, jwt.SigningMethodRS256, "ghost", k2, baseClaims(c)))
	}
	if h := f.hitCount(); h != base+1 {
		t.Fatalf("unknown kids caused %d extra fetches", h-base-1)
	}
}

func TestKeysUnavailable(t *testing.T) {
	c := &clock{t: time.Now().UTC().Truncate(time.Second)}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusInternalServerError)
	}))
	defer srv.Close()
	v := New(Config{Issuer: testIssuer, Audience: testProject, JWKSURL: srv.URL, Now: c.now})
	k, _ := rsa.GenerateKey(rand.Reader, 2048)
	_, err := v.Verify(context.Background(), sign(t, jwt.SigningMethodRS256, "k1", k, baseClaims(c)))
	if !errors.Is(err, ErrKeysUnavailable) {
		t.Errorf("err = %v, want ErrKeysUnavailable", err)
	}
}
