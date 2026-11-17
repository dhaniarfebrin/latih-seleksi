package httpx

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"latihseleksi/api/internal/auth"
	"latihseleksi/api/internal/bank"
	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/exam"
	"latihseleksi/api/internal/export"
	"latihseleksi/api/internal/gen/catalog"
	"latihseleksi/api/internal/gen/persentase"
	"latihseleksi/api/internal/insight"
	"latihseleksi/api/internal/review"
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
	items, err := bank.Load()
	if err != nil {
		t.Fatal(err)
	}
	examSvc := exam.New(st, specs, exam.NewGeneratedSource(catalog.All()), items)
	reviews := review.New(st)
	srv := httptest.NewServer(NewHandler(drill.New(st, catalog.All(), items, reviews), examSvc, insight.New(st, catalog.All(), reviews), export.New(st), reviews, st, cfg))
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

func create(t *testing.T, srv *httptest.Server, token, body string) createResp {
	t.Helper()
	code, out, _ := doToken(t, "POST", srv.URL+"/v1/attempts", body, token)
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
	for _, id := range []string{"TWK.nasionalisme", "TWK.integritas", "TWK.bela-negara", "TWK.pilar-negara", "TWK.bahasa-negara"} {
		if have[id] == 0 {
			t.Errorf("pattern %s missing in %s", id, out)
		}
	}
}

func TestFlowHidesAndRevealsKey(t *testing.T) {
	srv, mint := newAuthServer(t)
	tok := mint("flow-user")
	do := func(t *testing.T, method, url, body string) (int, string, http.Header) {
		t.Helper()
		return doToken(t, method, url, body, tok)
	}
	c := create(t, srv, tok, `{"pattern_id":"`+persentase.PatternID+`"}`) // defaults n=8 band=2
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

func TestTWKDrillFlow(t *testing.T) {
	srv, mint := newAuthServer(t)
	tok := mint("twk-flow-user")
	c := create(t, srv, tok, `{"pattern_id":"TWK.bela-negara","band":2}`)
	if len(c.Items) == 0 || len(c.Items) > 8 {
		t.Fatalf("items = %d", len(c.Items))
	}
	code, out, _ := doToken(t, "PUT", srv.URL+"/v1/attempts/"+c.ID+"/responses/1", `{"answer":"A","time_ms":1500}`, tok)
	if code != 200 || !strings.Contains(out, "answer_key") || !strings.Contains(out, "steps") {
		t.Errorf("PUT status = %d body=%s", code, out)
	}
}

func TestBadRequestsAndNotFound(t *testing.T) {
	srv, mint := newAuthServer(t)
	tok := mint("bad-requests-user")
	do := func(t *testing.T, method, url, body string) (int, string, http.Header) {
		t.Helper()
		return doToken(t, method, url, body, tok)
	}
	p := persentase.PatternID
	for _, body := range []string{
		`{"pattern_id":"` + p + `","n":0}`,
		`{"pattern_id":"` + p + `","n":21}`,
		`{"pattern_id":"` + p + `","band":0}`,
		`{"pattern_id":"` + p + `","band":4}`,
		`{"pattern_id":"unknown"}`,
		`{"pattern_id":"TWK.unknown"}`,
		`{"pattern_id":"` + p + `","extra":1}`,
		`not json`,
	} {
		if code, out, _ := do(t, "POST", srv.URL+"/v1/attempts", body); code != 400 {
			t.Errorf("POST %s: status %d, want 400 (%s)", body, code, out)
		}
	}

	c := create(t, srv, tok, `{"pattern_id":"`+p+`","n":2,"band":1}`)
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
	code, out, _ := do(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["XYZ"]}`)
	if code != http.StatusBadRequest {
		t.Fatalf("unknown section exam = %d %s, want 400", code, out)
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
}

func TestDrillRequiresAuth(t *testing.T) {
	srv, mint := newAuthServer(t)
	zero := "00000000-0000-0000-0000-000000000000"
	for _, r := range []struct{ method, path, body string }{
		{"POST", "/v1/attempts", `{"pattern_id":"` + persentase.PatternID + `"}`},
		{"GET", "/v1/attempts/" + zero, ""},
		{"PUT", "/v1/attempts/" + zero + "/responses/1", `{"answer":"A","time_ms":1}`},
		{"POST", "/v1/attempts/" + zero + "/submit", ""},
		{"DELETE", "/v1/me", ""},
		{"GET", "/v1/me/dashboard", ""},
	} {
		code, _, hdr := doToken(t, r.method, srv.URL+r.path, r.body, "")
		if code != http.StatusUnauthorized || hdr.Get("WWW-Authenticate") != "Bearer" {
			t.Errorf("%s %s without token = %d (WWW-Authenticate %q), want 401", r.method, r.path, code, hdr.Get("WWW-Authenticate"))
		}
		if code, _, _ := doToken(t, r.method, srv.URL+r.path, r.body, "not.a.jwt"); code != http.StatusUnauthorized {
			t.Errorf("%s %s with a bad token = %d, want 401", r.method, r.path, code)
		}
	}
	// The pattern list stays public.
	if code, _, _ := do(t, "GET", srv.URL+"/v1/patterns", ""); code != http.StatusOK {
		t.Errorf("/v1/patterns without token = %d, want 200", code)
	}
	// A valid token works.
	create(t, srv, mint("drill-auth-user"), `{"pattern_id":"`+persentase.PatternID+`","n":2,"band":1}`)
}

func TestExamRoutesWithoutVerifierAre503(t *testing.T) {
	srv := newServer(t, Config{})
	if code, out, _ := doToken(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TIU"]}`, "x"); code != http.StatusServiceUnavailable || !strings.Contains(out, "auth not configured") {
		t.Errorf("exam without verifier = %d %s, want 503", code, out)
	}
	if code, _, _ := doToken(t, "GET", srv.URL+"/v1/exams", "", "x"); code != http.StatusServiceUnavailable {
		t.Errorf("list without verifier = %d, want 503", code)
	}
	if code, _, _ := doToken(t, "POST", srv.URL+"/v1/attempts", `{"pattern_id":"`+persentase.PatternID+`"}`, "x"); code != http.StatusServiceUnavailable {
		t.Errorf("drill without verifier = %d, want 503", code)
	}
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

func TestDrillOwnership(t *testing.T) {
	srv, mint := newAuthServer(t)
	tokA, tokB := mint("drill-owner-a"), mint("drill-owner-b")
	c := create(t, srv, tokA, `{"pattern_id":"`+persentase.PatternID+`","n":2,"band":1}`)
	base := srv.URL + "/v1/attempts/" + c.ID
	if code, _, _ := doToken(t, "GET", base, "", tokA); code != http.StatusOK {
		t.Errorf("owner GET = %d, want 200", code)
	}
	for _, r := range []struct{ method, path, body string }{
		{"GET", "", ""},
		{"PUT", "/responses/1", `{"answer":"A","time_ms":1}`},
		{"POST", "/submit", ""},
	} {
		if code, _, _ := doToken(t, r.method, base+r.path, r.body, tokB); code != http.StatusNotFound {
			t.Errorf("other user %s %s = %d, want 404", r.method, r.path, code)
		}
	}
	// B's calls changed nothing for A.
	if _, out, _ := doToken(t, "GET", base, "", tokA); strings.Contains(out, "answer_key") || strings.Contains(out, "submitted_at") {
		t.Errorf("owner attempt changed: %s", out)
	}
}

func TestReportAttemptItem(t *testing.T) {
	srv, mint := newAuthServer(t)
	suffix := strconv.FormatInt(time.Now().UnixNano(), 36)
	tokA, tokB := mint("drill-report-a-"+suffix), mint("drill-report-b-"+suffix)
	c := create(t, srv, tokA, `{"pattern_id":"`+persentase.PatternID+`","n":2,"band":1}`)
	base := srv.URL + "/v1/attempts/" + c.ID
	url := base + "/items/1/report"
	body := `{"reason":"typo","note":"salah ketik"}`
	if code, out, _ := doToken(t, "POST", url, body, tokA); code != http.StatusBadRequest {
		t.Errorf("report on unanswered = %d %s, want 400", code, out)
	}
	if code, out, _ := doToken(t, "PUT", base+"/responses/1", `{"answer":"`+c.Items[0].Options[0].Key+`","time_ms":5}`, tokA); code != http.StatusOK {
		t.Fatalf("answer = %d %s", code, out)
	}
	for _, tc := range []struct {
		name, url, body, tok string
		want                 int
	}{
		{"ok", url, body, tokA, http.StatusNoContent},
		{"replace", url, `{"reason":"lainnya"}`, tokA, http.StatusNoContent},
		{"bad reason", url, `{"reason":"x"}`, tokA, http.StatusBadRequest},
		{"long note", url, `{"reason":"typo","note":"` + strings.Repeat("x", 501) + `"}`, tokA, http.StatusBadRequest},
		{"unanswered pos", base + "/items/2/report", body, tokA, http.StatusBadRequest},
		{"pos out of range", base + "/items/999/report", body, tokA, http.StatusNotFound},
		{"pos not a number", base + "/items/x/report", body, tokA, http.StatusNotFound},
		{"other user", url, body, tokB, http.StatusNotFound},
		{"no token", url, body, "", http.StatusUnauthorized},
	} {
		if code, out, _ := doToken(t, "POST", tc.url, tc.body, tc.tok); code != tc.want {
			t.Errorf("%s = %d %s, want %d", tc.name, code, out, tc.want)
		}
	}
}

func TestDeleteMe(t *testing.T) {
	srv, mint := newAuthServer(t)
	tokA, tokB := mint("delete-a"), mint("delete-b")
	drillA := create(t, srv, tokA, `{"pattern_id":"`+persentase.PatternID+`","n":2,"band":1}`)
	if code, out, _ := doToken(t, "PUT", srv.URL+"/v1/attempts/"+drillA.ID+"/responses/1", `{"answer":"A","time_ms":5}`, tokA); code != http.StatusOK {
		t.Fatalf("answer = %d %s", code, out)
	}
	code, out, _ := doToken(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TIU"]}`, tokA)
	if code != http.StatusCreated {
		t.Fatalf("create exam = %d %s", code, out)
	}
	var examA createResp
	if err := json.Unmarshal([]byte(out), &examA); err != nil {
		t.Fatal(err)
	}
	drillB := create(t, srv, tokB, `{"pattern_id":"`+persentase.PatternID+`","n":2,"band":1}`)

	if code, _, _ := doToken(t, "DELETE", srv.URL+"/v1/me", "", tokA); code != http.StatusNoContent {
		t.Fatalf("DELETE /v1/me = %d, want 204", code)
	}
	// The same token recreates an empty user: A's old data is gone.
	if code, _, _ := doToken(t, "GET", srv.URL+"/v1/attempts/"+drillA.ID, "", tokA); code != http.StatusNotFound {
		t.Errorf("deleted drill attempt = %d, want 404", code)
	}
	if code, _, _ := doToken(t, "GET", srv.URL+"/v1/exams/"+examA.ID, "", tokA); code != http.StatusNotFound {
		t.Errorf("deleted exam = %d, want 404", code)
	}
	if _, out, _ := doToken(t, "GET", srv.URL+"/v1/exams", "", tokA); !strings.Contains(out, `"exams":[]`) {
		t.Errorf("exam list after delete = %s", out)
	}
	// Repeating the call is idempotent.
	if code, _, _ := doToken(t, "DELETE", srv.URL+"/v1/me", "", tokA); code != http.StatusNoContent {
		t.Errorf("second DELETE /v1/me = %d, want 204", code)
	}
	// B is untouched.
	if code, _, _ := doToken(t, "GET", srv.URL+"/v1/attempts/"+drillB.ID, "", tokB); code != http.StatusOK {
		t.Errorf("other user's attempt = %d, want 200", code)
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

func TestCORSAllowsDelete(t *testing.T) {
	srv := newServer(t, Config{CORSOrigin: "http://localhost:5173"})
	_, _, hdr := do(t, "OPTIONS", srv.URL+"/v1/me", "")
	if !strings.Contains(hdr.Get("Access-Control-Allow-Methods"), "DELETE") {
		t.Errorf("Allow-Methods = %q", hdr.Get("Access-Control-Allow-Methods"))
	}
}

func TestCORSAllowsAuthorizationHeader(t *testing.T) {
	srv := newServer(t, Config{CORSOrigin: "http://localhost:5173"})
	_, _, hdr := do(t, "OPTIONS", srv.URL+"/v1/exams", "")
	if !strings.Contains(hdr.Get("Access-Control-Allow-Headers"), "Authorization") {
		t.Errorf("Allow-Headers = %q", hdr.Get("Access-Control-Allow-Headers"))
	}
}

func TestDashboard(t *testing.T) {
	srv, mint := newAuthServer(t)
	tokA, tokB := mint("dash-a"), mint("dash-b")
	type dash struct {
		Sections []struct {
			Section  string `json:"section"`
			Answered int    `json:"answered"`
		} `json:"sections"`
		Patterns []struct {
			PatternID string `json:"pattern_id"`
			Section   string `json:"section"`
			Answered  int    `json:"answered"`
		} `json:"patterns"`
		Sessions []struct {
			Mode string `json:"mode"`
		} `json:"sessions"`
		Queue []struct {
			Reason string `json:"reason"`
		} `json:"queue"`
	}
	get := func(tok string) (string, dash) {
		code, out, _ := doToken(t, "GET", srv.URL+"/v1/me/dashboard", "", tok)
		if code != http.StatusOK {
			t.Fatalf("dashboard = %d %s", code, out)
		}
		var d dash
		if err := json.Unmarshal([]byte(out), &d); err != nil {
			t.Fatal(err)
		}
		return out, d
	}

	// A new user gets empty arrays, not null, and a queue of untried patterns.
	raw, d := get(tokB)
	if !strings.Contains(raw, `"sections":[]`) || !strings.Contains(raw, `"patterns":[]`) || !strings.Contains(raw, `"sessions":[]`) {
		t.Errorf("new user dashboard = %s", raw)
	}
	if len(d.Queue) != 3 || d.Queue[0].Reason != "belum_dicoba" {
		t.Errorf("new user queue = %+v", d.Queue)
	}

	// A drill answer and a simulation answer both count; B's data stays out.
	drill := create(t, srv, tokA, `{"pattern_id":"`+persentase.PatternID+`","n":2,"band":1}`)
	if code, out, _ := doToken(t, "PUT", srv.URL+"/v1/attempts/"+drill.ID+"/responses/1", `{"answer":"A","time_ms":5}`, tokA); code != http.StatusOK {
		t.Fatalf("drill answer = %d %s", code, out)
	}
	code, out, _ := doToken(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TIU"]}`, tokA)
	if code != http.StatusCreated {
		t.Fatalf("create exam = %d %s", code, out)
	}
	var ex createResp
	if err := json.Unmarshal([]byte(out), &ex); err != nil {
		t.Fatal(err)
	}
	if code, out, _ := doToken(t, "PUT", srv.URL+"/v1/exams/"+ex.ID+"/responses/1", `{"answer":"A","time_ms":7,"seq":1}`, tokA); code != http.StatusOK {
		t.Fatalf("exam answer = %d %s", code, out)
	}
	_, d = get(tokA)
	if len(d.Sections) != 1 || d.Sections[0].Section != "TIU" || d.Sections[0].Answered != 2 {
		t.Errorf("sections = %+v, want TIU with 2 answers", d.Sections)
	}
	if len(d.Sessions) != 2 {
		t.Errorf("sessions = %+v, want drill and simulasi", d.Sessions)
	}

	// A TKP answer has a weight, not right or wrong: the dashboard ignores it.
	code, out, _ = doToken(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TKP"]}`, tokA)
	if code != http.StatusCreated {
		t.Fatalf("create TKP exam = %d %s", code, out)
	}
	var tkp createResp
	if err := json.Unmarshal([]byte(out), &tkp); err != nil {
		t.Fatal(err)
	}
	if code, out, _ := doToken(t, "PUT", srv.URL+"/v1/exams/"+tkp.ID+"/responses/1", `{"answer":"A","time_ms":9,"seq":1}`, tokA); code != http.StatusOK {
		t.Fatalf("TKP answer = %d %s", code, out)
	}
	raw, d = get(tokA)
	if len(d.Sections) != 1 || d.Sections[0].Answered != 2 || len(d.Sessions) != 2 || strings.Contains(raw, "TKP") {
		t.Errorf("TKP answer leaked into the dashboard: %s", raw)
	}
	found := false
	for _, p := range d.Patterns {
		if p.PatternID == persentase.PatternID && p.Section == "TIU" {
			found = true
		}
	}
	if !found {
		t.Errorf("patterns = %+v, want %s", d.Patterns, persentase.PatternID)
	}
	if _, db := get(tokB); len(db.Sections)+len(db.Patterns)+len(db.Sessions) != 0 {
		t.Errorf("other user sees data: %+v", db)
	}
}

// A full SKD simulation: 110 items from the generator and the bank; TKP is
// scored by the weight of the chosen option, shown once the exam is closed.
func TestFullSimulationScoresTKPByWeight(t *testing.T) {
	srv, mint := newAuthServer(t)
	tok := mint("full-sim")
	code, out, _ := doToken(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TWK","TIU","TKP"]}`, tok)
	if code != http.StatusCreated {
		t.Fatalf("create = %d %s", code, out)
	}
	var ex struct {
		ID    string `json:"id"`
		Items []struct {
			Pos     int    `json:"pos"`
			Section string `json:"section"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(out), &ex); err != nil || len(ex.Items) != 110 {
		t.Fatalf("exam body: %v items=%d", err, len(ex.Items))
	}
	if strings.Contains(out, `"weight"`) {
		t.Error("open exam exposes weights")
	}
	for _, it := range ex.Items {
		if it.Section != "TKP" {
			continue
		}
		path := srv.URL + "/v1/exams/" + ex.ID + "/responses/" + strconv.Itoa(it.Pos)
		if code, out, _ := doToken(t, "PUT", path, `{"answer":"A","time_ms":5,"seq":1}`, tok); code != http.StatusOK {
			t.Fatalf("answer %d = %d %s", it.Pos, code, out)
		}
	}
	code, out, _ = doToken(t, "POST", srv.URL+"/v1/exams/"+ex.ID+"/submit", "", tok)
	if code != http.StatusOK {
		t.Fatalf("submit = %d %s", code, out)
	}
	var res struct {
		Partial  bool `json:"partial"`
		Sections []struct {
			ID       string `json:"id"`
			Score    int    `json:"score"`
			Answered int    `json:"answered"`
		} `json:"sections"`
	}
	if err := json.Unmarshal([]byte(out), &res); err != nil || res.Partial || len(res.Sections) != 3 {
		t.Fatalf("result: %v %s", err, out)
	}
	code, out, _ = doToken(t, "GET", srv.URL+"/v1/exams/"+ex.ID, "", tok)
	if code != http.StatusOK {
		t.Fatalf("get = %d %s", code, out)
	}
	var closed struct {
		Items []struct {
			Section  string `json:"section"`
			Response *struct {
				Weight *int `json:"weight"`
			} `json:"response"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(out), &closed); err != nil {
		t.Fatal(err)
	}
	sum := 0
	for _, it := range closed.Items {
		if it.Section == "TKP" && it.Response != nil && it.Response.Weight != nil {
			sum += *it.Response.Weight
		}
	}
	tkp := res.Sections[2]
	if tkp.ID != "TKP" || tkp.Answered != 45 || sum < 45 || tkp.Score != sum {
		t.Errorf("TKP section = %+v, sum of shown weights = %d", tkp, sum)
	}
}

func TestReportItem(t *testing.T) {
	srv, mint := newAuthServer(t)
	suffix := strconv.FormatInt(time.Now().UnixNano(), 36)
	tokA, tokB := mint("report-a-"+suffix), mint("report-b-"+suffix)
	code, out, _ := doToken(t, "POST", srv.URL+"/v1/exams", `{"spec_id":"skd-2024","track":"umum","sections":["TIU"]}`, tokA)
	if code != http.StatusCreated {
		t.Fatalf("create = %d %s", code, out)
	}
	var c createResp
	if err := json.Unmarshal([]byte(out), &c); err != nil {
		t.Fatal(err)
	}
	url := srv.URL + "/v1/exams/" + c.ID + "/items/1/report"
	body := `{"reason":"typo","note":"salah ketik"}`
	if code, _, _ := doToken(t, "POST", url, body, tokA); code != http.StatusConflict {
		t.Errorf("report on open exam = %d, want 409", code)
	}
	if code, _, _ := doToken(t, "POST", srv.URL+"/v1/exams/"+c.ID+"/submit", "", tokA); code != http.StatusOK {
		t.Fatalf("submit = %d", code)
	}
	for _, tc := range []struct {
		name, url, body, tok string
		want                 int
	}{
		{"ok", url, body, tokA, http.StatusNoContent},
		{"replace", url, `{"reason":"lainnya"}`, tokA, http.StatusNoContent},
		{"bad reason", url, `{"reason":"x"}`, tokA, http.StatusBadRequest},
		{"long note", url, `{"reason":"typo","note":"` + strings.Repeat("x", 501) + `"}`, tokA, http.StatusBadRequest},
		{"unknown field", url, `{"reason":"typo","extra":1}`, tokA, http.StatusBadRequest},
		{"pos out of range", srv.URL + "/v1/exams/" + c.ID + "/items/999/report", body, tokA, http.StatusNotFound},
		{"pos not a number", srv.URL + "/v1/exams/" + c.ID + "/items/x/report", body, tokA, http.StatusNotFound},
		{"other user", url, body, tokB, http.StatusNotFound},
		{"no token", url, body, "", http.StatusUnauthorized},
	} {
		if code, out, _ := doToken(t, "POST", tc.url, tc.body, tc.tok); code != tc.want {
			t.Errorf("%s = %d %s, want %d", tc.name, code, out, tc.want)
		}
	}
}

func TestExportMe(t *testing.T) {
	srv, mint := newAuthServer(t)
	suffix := strconv.FormatInt(time.Now().UnixNano(), 36)
	subA, subB := "export-a-"+suffix, "export-b-"+suffix
	tokA, tokB := mint(subA), mint(subB)
	url := srv.URL + "/v1/me/export"
	if code, _, _ := do(t, "GET", url, ""); code != http.StatusUnauthorized {
		t.Errorf("no token = %d, want 401", code)
	}

	closed := create(t, srv, tokA, `{"pattern_id":"`+persentase.PatternID+`","n":2,"band":1}`)
	open := create(t, srv, tokA, `{"pattern_id":"`+persentase.PatternID+`","n":2,"band":1}`)
	for _, c := range []createResp{closed, open} {
		if code, out, _ := doToken(t, "PUT", srv.URL+"/v1/attempts/"+c.ID+"/responses/1", `{"answer":"`+c.Items[0].Options[0].Key+`","time_ms":5}`, tokA); code != http.StatusOK {
			t.Fatalf("answer = %d %s", code, out)
		}
	}
	if code, out, _ := doToken(t, "POST", srv.URL+"/v1/attempts/"+closed.ID+"/submit", "", tokA); code != http.StatusOK {
		t.Fatalf("submit = %d %s", code, out)
	}

	code, out, hdr := doToken(t, "GET", url, "", tokA)
	if code != http.StatusOK || !strings.Contains(hdr.Get("Content-Disposition"), "attachment") {
		t.Fatalf("export = %d %v %s", code, hdr, out)
	}
	if strings.Contains(out, subA) {
		t.Error("export leaks the auth subject")
	}
	var d struct {
		ExportedAt string `json:"exported_at"`
		Attempts   []struct {
			ID    string `json:"id"`
			Items []struct {
				Snapshot json.RawMessage `json:"snapshot"`
				Response *struct {
					Answer string `json:"answer"`
				} `json:"response"`
			} `json:"items"`
		} `json:"attempts"`
	}
	if err := json.Unmarshal([]byte(out), &d); err != nil {
		t.Fatal(err)
	}
	if d.ExportedAt == "" || len(d.Attempts) != 2 {
		t.Fatalf("export = %s", out)
	}
	for _, a := range d.Attempts {
		wantSnap := a.ID == closed.ID
		for i, it := range a.Items {
			if (len(it.Snapshot) > 0) != wantSnap {
				t.Errorf("attempt %s item %d snapshot present = %v, want %v", a.ID, i, len(it.Snapshot) > 0, wantSnap)
			}
		}
		if len(a.Items) != 2 || a.Items[0].Response == nil || a.Items[1].Response != nil {
			t.Errorf("attempt %s items = %+v", a.ID, a.Items)
		}
	}

	code, out, _ = doToken(t, "GET", url, "", tokB)
	if code != http.StatusOK || strings.Contains(out, closed.ID) || !strings.Contains(out, `"attempts":[]`) {
		t.Errorf("other user export = %d %s", code, out)
	}
}

// A drill submit schedules reviews once; due facts come back through ULANG with
// their real TWK pattern, and a due TIU pattern heads the dashboard queue.
func TestReviewFlow(t *testing.T) {
	srv, mint := newAuthServer(t)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	sub := "review-" + strconv.FormatInt(time.Now().UnixNano(), 36)
	tok := mint(sub)
	t.Cleanup(func() { doToken(t, "DELETE", srv.URL+"/v1/me", "", tok) })

	count := func(table string) int {
		t.Helper()
		var n int
		err := pool.QueryRow(ctx, `select count(*) from `+table+` t join app_user u on u.id = t.user_id where u.auth_sub = $1`, sub).Scan(&n)
		if err != nil {
			t.Fatal(err)
		}
		return n
	}
	answerAndSubmit := func(c createResp) {
		t.Helper()
		for _, it := range c.Items {
			if code, out, _ := doToken(t, "PUT", fmt.Sprintf("%s/v1/attempts/%s/responses/%d", srv.URL, c.ID, it.Pos), `{"answer":"A","time_ms":1000}`, tok); code != http.StatusOK {
				t.Fatalf("answer = %d %s", code, out)
			}
		}
		if code, out, _ := doToken(t, "POST", srv.URL+"/v1/attempts/"+c.ID+"/submit", "", tok); code != http.StatusOK {
			t.Fatalf("submit = %d %s", code, out)
		}
	}
	reviews := func() map[string]int {
		t.Helper()
		code, out, _ := doToken(t, "GET", srv.URL+"/v1/me/reviews", "", tok)
		if code != http.StatusOK {
			t.Fatalf("reviews = %d %s", code, out)
		}
		for _, banned := range []string{"answer_key", "snapshot", "stem"} {
			if strings.Contains(out, banned) {
				t.Errorf("reviews leaks %q: %s", banned, out)
			}
		}
		var m map[string]int
		if err := json.Unmarshal([]byte(out), &m); err != nil {
			t.Fatal(err)
		}
		return m
	}

	twk := create(t, srv, tok, `{"pattern_id":"TWK.bela-negara","n":3}`)
	answerAndSubmit(twk)
	doToken(t, "POST", srv.URL+"/v1/attempts/"+twk.ID+"/submit", "", tok) // a second submit must add nothing
	if count("review_card") != 3 || count("review_log") != 3 {
		t.Fatalf("after submit twice: %d cards, %d logs, want 3 each", count("review_card"), count("review_log"))
	}
	var future int
	if err := pool.QueryRow(ctx, `select count(*) from review_card c join app_user u on u.id = c.user_id where u.auth_sub = $1 and c.due > now()`, sub).Scan(&future); err != nil || future != 3 {
		t.Errorf("cards due in the future = %d (%v), want 3", future, err)
	}
	if m := reviews(); m["due_count"] != 0 {
		t.Errorf("reviews before due = %v", m)
	}
	if code, _, _ := doToken(t, "POST", srv.URL+"/v1/attempts", `{"pattern_id":"ULANG"}`, tok); code != http.StatusBadRequest {
		t.Errorf("ULANG with nothing due = %d, want 400", code)
	}

	tiu := create(t, srv, tok, `{"pattern_id":"`+persentase.PatternID+`","n":4,"band":1}`)
	answerAndSubmit(tiu)
	if count("review_card") != 4 {
		t.Errorf("cards after TIU drill = %d, want 4", count("review_card"))
	}
	if _, err := pool.Exec(ctx, `update review_card set due = now() - interval '1 day' where user_id = (select id from app_user where auth_sub = $1)`, sub); err != nil {
		t.Fatal(err)
	}
	if m := reviews(); m["due_count"] != 4 || m["facts_due"] != 3 || m["patterns_due"] != 1 {
		t.Errorf("reviews = %v, want 3 facts + 1 pattern", m)
	}

	code, out, _ := doToken(t, "GET", srv.URL+"/v1/me/dashboard", "", tok)
	var dash struct {
		Queue []struct {
			PatternID string `json:"pattern_id"`
			Reason    string `json:"reason"`
		} `json:"queue"`
	}
	if err := json.Unmarshal([]byte(out), &dash); err != nil || code != http.StatusOK || len(dash.Queue) == 0 ||
		dash.Queue[0].Reason != "jatuh_tempo" || dash.Queue[0].PatternID != persentase.PatternID {
		t.Errorf("dashboard queue = %s", out)
	}

	rev := create(t, srv, tok, `{"pattern_id":"ULANG","n":8}`)
	if len(rev.Items) != 3 {
		t.Fatalf("ULANG items = %d, want the 3 due facts", len(rev.Items))
	}
	answerAndSubmit(rev)
	if count("review_log") != 7 {
		t.Errorf("logs after ULANG = %d, want 7", count("review_log"))
	}
	code, out, _ = doToken(t, "GET", srv.URL+"/v1/me/dashboard", "", tok)
	if code != http.StatusOK {
		t.Fatalf("dashboard = %d %s", code, out)
	}
	var d2 struct {
		Patterns []struct {
			PatternID string `json:"pattern_id"`
			Section   string `json:"section"`
			Answered  int    `json:"answered"`
		} `json:"patterns"`
	}
	if err := json.Unmarshal([]byte(out), &d2); err != nil {
		t.Fatal(err)
	}
	twkAnswers := 0
	for _, p := range d2.Patterns {
		if p.PatternID == "ULANG" {
			t.Errorf("dashboard attributes answers to ULANG: %s", out)
		}
		if p.PatternID == "TWK.bela-negara" && p.Section == "TWK" {
			twkAnswers = p.Answered
		}
	}
	if twkAnswers != 6 { // 3 from the first drill, 3 from the review drill
		t.Errorf("TWK.bela-negara answered = %d, want 6 in %s", twkAnswers, out)
	}

	code, out, _ = doToken(t, "GET", srv.URL+"/v1/me/export", "", tok)
	var ex struct {
		Cards []struct{ Kind, Key string } `json:"review_cards"`
		Log   []struct{ Kind, Key string } `json:"review_log"`
	}
	if err := json.Unmarshal([]byte(out), &ex); err != nil || code != http.StatusOK || len(ex.Cards) != 4 || len(ex.Log) != 7 {
		t.Errorf("export = %d %v cards %d log %d", code, err, len(ex.Cards), len(ex.Log))
	}
}
