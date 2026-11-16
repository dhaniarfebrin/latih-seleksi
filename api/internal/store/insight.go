package store

import (
	"context"

	"github.com/jackc/pgx/v5"

	"latihseleksi/api/internal/insight"
)

// ResponseFacts returns the user's latest answered responses across drill and
// simulation attempts. Drill items have no section, so it is the first segment
// of the pattern id ("TIU.numerik.x" is TIU); drill items without their own
// pattern use the attempt's. Weighted (TKP) items, whose options carry a
// weight, are flagged so the dashboard can skip them. This relies on the
// snapshot weight key being omitted when 0: bank validation requires every
// TKP option weight to be 1-5, so a weighted item always has the key.
func (s *Store) ResponseFacts(ctx context.Context, userID string, limit int) ([]insight.Fact, error) {
	rows, err := s.pool.Query(ctx,
		`select a.id::text, a.mode, a.started_at,
		        coalesce(i.section, split_part(coalesce(i.pattern_id, a.pattern_id, ''), '.', 1)),
		        coalesce(i.pattern_id, a.pattern_id, ''),
		        r.correct,
		        coalesce(exists (select 1 from jsonb_array_elements(i.item_snapshot->'options') o where o ? 'weight'), false),
		        r.time_ms, coalesce(r.misconception_id, '')
		 from response r
		 join attempt a on a.id = r.attempt_id
		 left join attempt_item i on i.attempt_id = r.attempt_id and i.pos = r.pos
		 where a.user_id = $1
		 order by r.answered_at desc, r.attempt_id, r.pos
		 limit $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, func(r pgx.CollectableRow) (insight.Fact, error) {
		var f insight.Fact
		err := r.Scan(&f.AttemptID, &f.Mode, &f.StartedAt, &f.Section, &f.PatternID, &f.Correct, &f.Weighted, &f.TimeMS, &f.MisconceptionID)
		return f, err
	})
}
