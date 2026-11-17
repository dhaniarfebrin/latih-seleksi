package store

import (
	"context"
	"strconv"
	"testing"
	"time"

	"latihseleksi/api/internal/question"
)

func TestExportUserHidesKeyOfRunningAttempt(t *testing.T) {
	st := openTest(t)
	ctx := context.Background()
	uid, err := st.UpsertUser(ctx, "export-test-"+strconv.FormatInt(time.Now().UnixNano(), 36))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = st.DeleteUser(ctx, uid) })

	now := time.Now()
	items := []question.Item{{Pos: 1, Seed: 1, Snapshot: []byte(`{"options":[{"key":"A"}]}`), Section: "TIU", PatternID: "px"}}
	running, err := st.CreateExam(ctx, question.Exam{UserID: uid, SpecID: "s", SpecBody: []byte(`{}`), Track: "umum", Sections: []string{"TIU"}, StartedAt: now, DeadlineAt: now.Add(time.Hour)}, items)
	if err != nil {
		t.Fatal(err)
	}
	d, err := st.ExportUser(ctx, uid)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Attempts) != 1 || d.Attempts[0].ID != running.ID || len(d.Attempts[0].Items) != 1 {
		t.Fatalf("export = %+v", d)
	}
	if d.Attempts[0].Items[0].Snapshot != nil {
		t.Error("running attempt exported its snapshot")
	}
	if _, err := st.pool.Exec(ctx, `update attempt set status = 'submitted', submitted_at = now() where id = $1`, running.ID); err != nil {
		t.Fatal(err)
	}
	d, err = st.ExportUser(ctx, uid)
	if err != nil {
		t.Fatal(err)
	}
	if d.Attempts[0].Items[0].Snapshot == nil {
		t.Error("closed attempt lacks its snapshot")
	}
}
