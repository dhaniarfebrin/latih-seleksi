package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

const (
	// minRefetch is the shortest time between two fetches triggered by an
	// unknown kid, so a stream of forged tokens cannot hammer the key endpoint.
	minRefetch = time.Minute
	// defaultTTL applies when the response carries no usable max-age.
	defaultTTL   = time.Hour
	maxJWKSBytes = 1 << 20
)

var maxAgeRE = regexp.MustCompile(`(?i)max-age=(\d+)`)

// keySet caches the RSA public keys of a JWKS endpoint.
type keySet struct {
	url    string
	client *http.Client
	now    func() time.Time

	mu        sync.Mutex
	keys      map[string]*rsa.PublicKey
	expires   time.Time
	lastFetch time.Time
}

// get returns the key for kid. A fresh cache answers without a request; an
// unknown kid triggers one refetch, at most once per minRefetch.
func (k *keySet) get(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	k.mu.Lock()
	defer k.mu.Unlock()
	now := k.now()
	fresh := k.keys != nil && now.Before(k.expires)
	if fresh {
		if key, ok := k.keys[kid]; ok {
			return key, nil
		}
		if now.Sub(k.lastFetch) < minRefetch {
			return nil, ErrInvalidToken
		}
	}
	if err := k.fetch(ctx, now); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrKeysUnavailable, err)
	}
	if key, ok := k.keys[kid]; ok {
		return key, nil
	}
	return nil, ErrInvalidToken
}

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// fetch replaces the cache with the endpoint's current keys. Callers hold mu.
func (k *keySet) fetch(ctx context.Context, now time.Time) error {
	k.lastFetch = now
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, k.url, nil)
	if err != nil {
		return err
	}
	resp, err := k.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("jwks endpoint returned status %d", resp.StatusCode)
	}
	var doc struct {
		Keys []jwk `json:"keys"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxJWKSBytes)).Decode(&doc); err != nil {
		return fmt.Errorf("decode jwks: %w", err)
	}
	keys := make(map[string]*rsa.PublicKey, len(doc.Keys))
	for _, j := range doc.Keys {
		if j.Kty != "RSA" || j.Kid == "" {
			continue
		}
		pub, err := j.public()
		if err != nil {
			continue
		}
		keys[j.Kid] = pub
	}
	if len(keys) == 0 {
		return errors.New("jwks has no usable RSA key")
	}
	ttl := defaultTTL
	if m := maxAgeRE.FindStringSubmatch(resp.Header.Get("Cache-Control")); m != nil {
		if s, err := strconv.Atoi(m[1]); err == nil && s > 0 {
			ttl = time.Duration(s) * time.Second
		}
	}
	k.keys, k.expires = keys, now.Add(ttl)
	return nil
}

func (j jwk) public() (*rsa.PublicKey, error) {
	nb, err := base64.RawURLEncoding.DecodeString(j.N)
	if err != nil {
		return nil, err
	}
	eb, err := base64.RawURLEncoding.DecodeString(j.E)
	if err != nil {
		return nil, err
	}
	e := new(big.Int).SetBytes(eb)
	if !e.IsInt64() || e.Int64() < 3 || e.Int64() > 1<<31-1 {
		return nil, errors.New("bad exponent")
	}
	return &rsa.PublicKey{N: new(big.Int).SetBytes(nb), E: int(e.Int64())}, nil
}
