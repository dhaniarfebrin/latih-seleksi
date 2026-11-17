package store

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5"

	"latihseleksi/api/internal/export"
)

// ExportUser returns the user's attempts, items, responses and reports. The
// frozen snapshot holds the answer key, so it is read only for closed attempts
// (drill submitted, simulation not in progress); the key is never sent before
// submit.
func (s *Store) ExportUser(ctx context.Context, userID string) (export.Data, error) {
	var d export.Data
	if err := s.pool.QueryRow(ctx, `select created_at from app_user where id = $1`, userID).Scan(&d.Account.CreatedAt); err != nil {
		return export.Data{}, err
	}

	rows, err := s.pool.Query(ctx,
		`select id::text, mode, coalesce(pattern_id, ''), band, n, coalesce(spec_id, ''),
		        coalesce(track, ''), coalesce(status, ''), started_at, submitted_at
		 from attempt where user_id = $1
		 order by started_at, id`, userID)
	if err != nil {
		return export.Data{}, err
	}
	d.Attempts, err = pgx.CollectRows(rows, func(r pgx.CollectableRow) (export.Attempt, error) {
		var a export.Attempt
		var band, n *int16
		err := r.Scan(&a.ID, &a.Mode, &a.PatternID, &band, &n, &a.SpecID, &a.Track, &a.Status, &a.StartedAt, &a.SubmittedAt)
		if band != nil {
			v := int(*band)
			a.Band = &v
		}
		if n != nil {
			v := int(*n)
			a.N = &v
		}
		a.Items = []export.Item{}
		return a, err
	})
	if err != nil {
		return export.Data{}, err
	}
	byID := make(map[string]int, len(d.Attempts))
	for i, a := range d.Attempts {
		byID[a.ID] = i
	}

	irows, err := s.pool.Query(ctx,
		`select a.id::text, i.pos, coalesce(i.section, ''), coalesce(i.pattern_id, ''),
		        case when (a.mode = 'drill' and a.submitted_at is not null)
		               or (a.mode = 'simulasi' and a.status <> 'in_progress')
		             then i.item_snapshot end,
		        r.answer, r.correct, r.time_ms, coalesce(r.misconception_id, ''), r.flagged, r.answered_at,
		        q.reason, q.note, q.created_at
		 from attempt a
		 join attempt_item i on i.attempt_id = a.id
		 left join response r on r.attempt_id = i.attempt_id and r.pos = i.pos
		 left join question_report q on q.attempt_id = i.attempt_id and q.pos = i.pos
		 where a.user_id = $1
		 order by a.started_at, a.id, i.pos`, userID)
	if err != nil {
		return export.Data{}, err
	}
	defer irows.Close()
	for irows.Next() {
		var attemptID string
		var it export.Item
		var snap []byte
		var answer, reason, note *string
		var correct, flagged *bool
		var timeMS *int32
		var misc string
		var answeredAt, reportAt *time.Time
		if err := irows.Scan(&attemptID, &it.Pos, &it.Section, &it.PatternID, &snap,
			&answer, &correct, &timeMS, &misc, &flagged, &answeredAt, &reason, &note, &reportAt); err != nil {
			return export.Data{}, err
		}
		if snap != nil {
			it.Snapshot = json.RawMessage(snap)
		}
		if answer != nil {
			it.Response = &export.Response{Answer: *answer, Correct: *correct, TimeMS: int(*timeMS),
				MisconceptionID: misc, Flagged: *flagged, AnsweredAt: *answeredAt}
		}
		if reason != nil {
			it.Report = &export.Report{Reason: *reason, Note: *note, CreatedAt: *reportAt}
		}
		i := byID[attemptID]
		d.Attempts[i].Items = append(d.Attempts[i].Items, it)
	}
	if err := irows.Err(); err != nil {
		return export.Data{}, err
	}
	return d, nil
}
