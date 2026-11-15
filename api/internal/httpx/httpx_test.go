package httpx

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"latihseleksi/api/internal/auth"
	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/exam"
	"latihseleksi/api/internal/gen/catalog"
	"latihseleksi/api/internal/gen/persentase"
	"latihseleksi/api/internal/spec"
	"latihseleksi/api/internal/store"
)

func newServer(t *testing.T, cfg Config) *httptest.Server {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		t.Skip("DATABASE_URL not set")
	}
	st, err := store.Open(context.Background(), url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(st.Close)
	if err := st.Migrate(context.Background()); err != nil {
		t.Fatal(err)
	}
	specs, err := spec.Load()
	if err != nil {
		t.Fatal(err)
	}
	examSvc := exam.New(st, specs, exam.NewGeneratedSource(catalog.All()))
	srv := httptest.NewServer(NewHandler(drill.New(st, catalog.All()), examSvc, st, cfg))
	t.Cleanup(srv.Close)
	return srv
}

func do(t *testing.T, method, url, body string) (int, string, http.Header) {
	t.Helper()
	return doToken(t, method, url, body, "")
}

func doToken(t *testing.T, method, url, body, token string) (int, string, http.Header) {
	t.Helper()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(b), resp.Header
}

type createResp struct {
	ID    string `json:"id"`
	Items []struct {
		Pos     int `json:"pos"`
		Options []struct {
			Key string `json:"key"`
		} `json:"options"`
	} `json:"items"`
}

func create(t *testing.T, srv *httptest.Server, body string) createResp {
	t.Helper()
	code, out, _ := do(t, "POST", srv.URL+"/v1/attempts", body)
	if code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s", code, out)
	}
	var c createResp
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatal(err)
	}
	return c
}

func TestHealthAndPatterns(t *testing.T) {
	srv := newServer(t, Config{})
	if code, _, hdr := do(t, "GET", srv.URL+"/healthz", ""); code != 200 || hdr.Get("X-Request-Id") == "" {
		t.Errorf("healthz code=%d request-id=%q", code, hdr.Get("X-Request-Id"))
	}
	code, out, _ := do(t, "GET", srv.URL+"/v1/patterns", "")
	if code != 200 {
		t.Fatalf("patterns code=%d body=%s", code, out)
	}
	var body struct {
		Patterns []struct {
			ID             string            `json:"id"`
			Misconceptions map[string]string `json:"misconceptions"`
		} `json:"patterns"`
	}
	if err := json.Unmarshal([]byte(out), &body); err != nil {
		t.Fatal(err)
	}
	have := map[string]int{}
	for _, p := range body.Patterns {
		have[p.ID] = len(p.Misconceptions)
	}
	for _, f := range catalog.All() {
		id := f.ID()
		if have[id] == 0 {
			t.Errorf("pattern %s missing or without misconceptions in %s", id, out)
		}
	}
}

func TestFlowHidesAndRevealsKey(t *testing.T) {
	srv := newServer(t, Config{})
	c := create(t, srv, `{"pattern_id":"`+persentase.PatternID+`"}`) // defaults n=8 band=2
	if len(c.Items) != 8 {
		t.Fatalf("default n: items = %d, want 8", len(c.Items))
	}
	_, got, _ := do(t, "GET", srv.URL+"/v1/attempts/"+c.ID, "")
	for _, banned := range []string{"answer_key", "steps", "misconception_id"} {
		if strings.Contains(got, banned) {
			t.Errorf("GET before answer contains %q", banned)
		}
	}

	url := srv.URL + "/v1/attempts/" + c.ID + "/responses/1"
	code, first, _ := do(t, "PUT", url, `{"answer":"A","time_ms":1500}`)
	if code != 200 {
		t.Fatalf("PUT status = %d body=%s", code, first)
	}
	for _, want := range []string{"answer_key", "steps"} {
		if !strings.Contains(first, want) {
			t.Errorf("PUT response missing %q: %s", want, first)
		}
	}
	code, second, _ := do(t, "PUT", url, `{"answer":"B","time_ms":9}`)
	if code != 200 || second != first {
		t.Errorf("second PUT not idempotent:\nfirst  = %s\nsecond = %s", first, second)
	}
	_, after, _ := do(t, "GET", srv.URL+"/v1/attempts/"+c.ID, "")
	if !strings.Contains(after, "answer_key") {
		t.Errorf("GET after answer lacks answer_key")
	}

	code, res, _ := do(t, "POST", srv.URL+"/v1/attempts/"+c.ID+"/submit", "")
	if code != 200 {
		t.Fatalf("submit status = %d body=%s", code, res)
	}
	var r struct {
		N        int `json:"n"`
		Answered int `json:"answered"`
		Score    int `json:"score"`
	}
	if err := json.Unmarshal([]byte(res), &r); err != nil || r.N != 8 || r.Answered != 1 {
		t.Errorf("submit result = %s err=%v", res, err)
	}
}

func TestBadRequestsAndNotFound(t *testing.T) {
	srv := newServer(t, Config{})
	p := persentase.PatternID
	for _, body := range []string{
		`{"pattern_id":"` + p + `","n":0}`,
		`{"pattern_id":"` + p + `","n":21}`,
		`{"pattern_id":"` + p + `","band":0}`,
		`{"pattern_id":"` + p + `","band":4}`,
		`{"pattern_id":"unknown"}`,
		`{"pattern_id":"` + p + `","extra":1}`,
		`not json`,
	} {
		if code, out, _ := do(t, "POST", srv.URL+"/v1/attempts", body); code != 400 {
			t.Errorf("POST %s: status %d, want 400 (%s)", body, code, out)
		}
	}

	c := create(t, srv, `{"pattern_id":"`+p+`","n":2,"band":1}`)
	base := srv.URL + "/v1/attempts/" + c.ID
	for _, body := range []string{`{"answer":"F","time_ms":1}`, `{"answer":"","time_ms":1}`, `{"answer":"A","time_ms":-5}`} {
		if code, out, _ := do(t, "PUT", base+"/responses/1", body); code != 400 {
			t.Errorf("PUT %s: status %d, want 400 (%s)", body, code, out)
		}
	}
	for _, pos := range []string{"0", "3", "99", "x"} {
		if code, out, _ := do(t, "PUT", base+"/responses/"+pos, `{"answer":"A","time_ms":1}`); code != 404 {
			t.Errorf("PUT pos %s: status %d, want 404 (%s)", pos, code, out)
		}
	}
	for _, id := range []string{"00000000-0000-0000-0000-000000000000", "bogus"} {
		if code, _, _ := do(t, "GET", srv.URL+"/v1/attempts/"+id, ""); code != 404 {
			t.Errorf("GET %s: status %d, want 404", id, code)
		}
		if code, _, _ := do(t, "POST", srv.URL+"/v1/attempts/"+id+"/submit", ""); code != 404 {
			t.Errorf("submit %s: status %d, want 404", id, code)
		}
	}
	// Body larger than the cap is rejected.
	big := `{"pattern_id":"` + strings.Repeat("x", maxBodyBytes+10) + `"}`
	if code, _, _ := do(t, "POST", srv.URL+"/v1/attempts", big); code != 400 {
		t.Errorf("oversized body status %d, want 400", code)
	}
}

func TestCORS(t *testing.T) {
	srv := newServer(t, Config{CORSOrigin: "http://localhost:5173"})
	code, _, hdr := do(t, "OPTIONS", srv.URL+"/v1/attempts", "")
	if code != http.StatusNoContent || hdr.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Errorf("preflight code=%d origin=%q", code, hdr.Get("Access-Control-Allow-Origin"))
	}
	_, _, hdr = do(t, "GET", srv.URL+"/healthz", "")
	if hdr.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Errorf("GET lacks CORS header")
	}
}

func TestExamFlowAndStatusCodes(t *testing.T) {
	srv, tok := newAuthServer(t)
	tokA := tok("user-a")
	do := func(t *testing.T, method, url, body string) (int, string, http.Header) {
		t.Helper()
		return doToken(t, method, url, body, tokA)
	}
	if code, out, _ := do(t, "GET", srv.URL+"/v1/specs", ""); code != http.StatusOK || !strings.Contains(out, `"skd-2024"`) {
		t.Fatalf("specs: %d %s", code, out)
	}
	code, out, _ := do(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TWK"]}`)
	if code != http.StatusBadRequest {
		t.Fatalf("TWK exam = %d %s, want 400", code, out)
	}
	code, out, _ = do(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TIU"]}`)
	if code != http.StatusCreated {
		t.Fatalf("create exam = %d %s", code, out)
	}
	var c createResp
	if err := json.Unmarshal([]byte(out), &c); err != nil || len(c.Items) != 35 {
		t.Fatalf("exam body: %v items=%d", err, len(c.Items))
	}
	for _, banned := range []string{"answer_key", "steps"} {
		if strings.Contains(out, banned) {
			t.Errorf("create response leaks %q", banned)
		}
	}
	if code, _, _ := do(t, "GET", srv.URL+"/v1/exams/"+c.ID+"/result", ""); code != http.StatusConflict {
		t.Errorf("result while open = %d, want 409", code)
	}
	url := srv.URL + "/v1/exams/" + c.ID + "/responses/1"
	if code, out, _ := do(t, "PUT", url, `{"answer":"A","time_ms":5,"seq":1}`); code != http.StatusOK {
		t.Fatalf("put = %d %s", code, out)
	}
	if code, _, _ := do(t, "PUT", url, `{"answer":"Z","time_ms":5,"seq":2}`); code != http.StatusBadRequest {
		t.Errorf("bad answer = %d, want 400", code)
	}
	if code, out, _ := do(t, "POST", srv.URL+"/v1/exams/"+c.ID+"/submit", ""); code != http.StatusOK || !strings.Contains(out, `"partial":true`) {
		t.Fatalf("submit = %d %s", code, out)
	}
	if code, _, _ := do(t, "PUT", url, `{"answer":"B","time_ms":5,"seq":3}`); code != http.StatusConflict {
		t.Errorf("put after submit = %d, want 409", code)
	}
	if code, _, _ := do(t, "GET", srv.URL+"/v1/attempts/"+c.ID, ""); code != http.StatusNotFound {
		t.Errorf("drill route on an exam = %d, want 404", code)
	}
}

const testProject = "demo-project"

// newAuthServer starts the API with a verifier backed by a fake JWKS and
// returns a function that mints a valid ID token for a given sub.
func newAuthServer(t *testing.T) (*httptest.Server, func(sub string) string) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	jwks := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{{
			"kid": "k1", "kty": "RSA", "alg": "RS256", "use": "sig",
			"n": base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
			"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
		}}})
	}))
	t.Cleanup(jwks.Close)
	v := auth.New(auth.Config{
		Issuer: "https://securetoken.google.com/" + testProject, Audience: testProject, JWKSURL: jwks.URL,
	})
	mint := func(sub string) string {
		now := time.Now()
		tok := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"iss": "https://securetoken.google.com/" + testProject, "aud": testProject, "sub": sub,
			"iat": now.Add(-time.Minute).Unix(), "auth_time": now.Add(-time.Minute).Unix(),
			"exp": now.Add(time.Hour).Unix(),
		})
		tok.Header["kid"] = "k1"
		s, err := tok.SignedString(key)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	return newServer(t, Config{Auth: v}), mint
}

func TestExamRoutesRequireAuth(t *testing.T) {
	srv, mint := newAuthServer(t)
	body := `{"spec_id":"skd-2024","track":"umum","sections":["TIU"]}`
	zero := "00000000-0000-0000-0000-000000000000"
	routes := []struct{ method, path, body string }{
		{"POST", "/v1/exams", body},
		{"GET", "/v1/exams", ""},
		{"GET", "/v1/exams/" + zero, ""},
		{"PUT", "/v1/exams/" + zero + "/responses/1", `{"answer":"A","time_ms":1,"seq":1}`},
		{"POST", "/v1/exams/" + zero + "/submit", ""},
		{"GET", "/v1/exams/" + zero + "/result", ""},
	}
	for _, r := range routes {
		code, _, hdr := doToken(t, r.method, srv.URL+r.path, r.body, "")
		if code != http.StatusUnauthorized || hdr.Get("WWW-Authenticate") != "Bearer" {
			t.Errorf("%s %s without token = %d (WWW-Authenticate %q), want 401", r.method, r.path, code, hdr.Get("WWW-Authenticate"))
		}
		if code, _, _ := doToken(t, r.method, srv.URL+r.path, r.body, "not.a.jwt"); code != http.StatusUnauthorized {
			t.Errorf("%s %s with a bad token = %d, want 401", r.method, r.path, code)
		}
	}
	// Wrong scheme with a valid token.
	req, _ := http.NewRequest("POST", srv.URL+"/v1/exams", strings.NewReader(body))
	req.Header.Set("Authorization", "Basic "+mint("user-a"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Basic scheme = %d, want 401", resp.StatusCode)
	}
	// Public routes stay open.
	if code, _, _ := do(t, "GET", srv.URL+"/v1/specs", ""); code != http.StatusOK {
		t.Errorf("/v1/specs without token = %d, want 200", code)
	}
	create(t, srv, `{"pattern_id":"`+persentase.PatternID+`","n":2,"band":1}`)
}

func TestExamRoutesWithoutVerifierAre503(t *testing.T) {
	srv := newServer(t, Config{})
	if code, out, _ := doToken(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TIU"]}`, "x"); code != http.StatusServiceUnavailable || !strings.Contains(out, "auth not configured") {
		t.Errorf("exam without verifier = %d %s, want 503", code, out)
	}
	if code, _, _ := doToken(t, "GET", srv.URL+"/v1/exams", "", "x"); code != http.StatusServiceUnavailable {
		t.Errorf("list without verifier = %d, want 503", code)
	}
	create(t, srv, `{"pattern_id":"`+persentase.PatternID+`","n":2,"band":1}`)
}

func TestExamOwnership(t *testing.T) {
	srv, mint := newAuthServer(t)
	tokA, tokB := mint("owner-a"), mint("owner-b")
	code, out, _ := doToken(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TIU"]}`, tokA)
	if code != http.StatusCreated {
		t.Fatalf("create = %d %s", code, out)
	}
	var c createResp
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatal(err)
	}
	base := srv.URL + "/v1/exams/" + c.ID
	if code, _, _ := doToken(t, "GET", base, "", tokA); code != http.StatusOK {
		t.Errorf("owner GET = %d, want 200", code)
	}
	for _, r := range []struct{ method, path, body string }{
		{"GET", "", ""},
		{"PUT", "/responses/1", `{"answer":"A","time_ms":1,"seq":1}`},
		{"POST", "/submit", ""},
		{"GET", "/result", ""},
	} {
		if code, _, _ := doToken(t, r.method, base+r.path, r.body, tokB); code != http.StatusNotFound {
			t.Errorf("other user %s %s = %d, want 404", r.method, r.path, code)
		}
	}
	// B's attempts did not touch A's exam.
	if _, out, _ := doToken(t, "GET", base, "", tokA); !strings.Contains(out, `"status":"in_progress"`) {
		t.Errorf("owner exam changed: %s", out)
	}
}

func TestListExams(t *testing.T) {
	srv, mint := newAuthServer(t)
	tokA, tokB := mint("list-a"), mint("list-b")
	code, out, _ := doToken(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TIU"]}`, tokA)
	if code != http.StatusCreated {
		t.Fatalf("create = %d %s", code, out)
	}
	var c createResp
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatal(err)
	}
	code, out, _ = doToken(t, "GET", srv.URL+"/v1/exams?limit=5", "", tokA)
	if code != http.StatusOK || !strings.Contains(out, c.ID) || !strings.Contains(out, `"status":"in_progress"`) {
		t.Errorf("owner list = %d %s", code, out)
	}
	if strings.Contains(out, "answer_key") || strings.Contains(out, `"items"`) {
		t.Errorf("list must not carry items or keys: %s", out)
	}
	code, out, _ = doToken(t, "GET", srv.URL+"/v1/exams", "", tokB)
	if code != http.StatusOK || strings.Contains(out, c.ID) || !strings.Contains(out, `"exams":[]`) {
		t.Errorf("other user list = %d %s", code, out)
	}
	for _, q := range []string{"?limit=0", "?limit=51", "?limit=abc"} {
		if code, _, _ := doToken(t, "GET", srv.URL+"/v1/exams"+q, "", tokA); code != http.StatusBadRequest {
			t.Errorf("GET /v1/exams%s = %d, want 400", q, code)
		}
	}
}

func TestCORSAllowsAuthorizationHeader(t *testing.T) {
	srv := newServer(t, Config{CORSOrigin: "http://localhost:5173"})
	_, _, hdr := do(t, "OPTIONS", srv.URL+"/v1/exams", "")
	if !strings.Contains(hdr.Get("Access-Control-Allow-Headers"), "Authorization") {
		t.Errorf("Allow-Headers = %q", hdr.Get("Access-Control-Allow-Headers"))
	}
}
