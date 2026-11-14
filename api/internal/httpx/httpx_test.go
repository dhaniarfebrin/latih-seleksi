package httpx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"latihseleksi/api/internal/drill"
	"latihseleksi/api/internal/gen/kecepatan"
	"latihseleksi/api/internal/gen/kerjasama"
	"latihseleksi/api/internal/gen/persentase"
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
	srv := httptest.NewServer(NewHandler(drill.New(st), cfg))
	t.Cleanup(srv.Close)
	return srv
}

func do(t *testing.T, method, url, body string) (int, string, http.Header) {
	t.Helper()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
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
	for _, id := range []string{persentase.PatternID, kerjasama.PatternID, kecepatan.PatternID} {
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
