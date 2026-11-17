-- FSRS review scheduling (PRD FR-11). review_card is the current memory state
-- of one item per user: kind 'fact' = TWK bank entry id, kind 'pattern' = TIU
-- pattern id. review_log keeps every review for later calibration; its primary
-- key makes a drill's review idempotent (attempt_id, kind, key).

create table review_card (
    user_id     uuid             not null references app_user (id) on delete cascade,
    kind        text             not null check (kind in ('fact', 'pattern')),
    key         text             not null,
    stability   double precision not null,
    difficulty  double precision not null,
    due         timestamptz      not null,
    last_review timestamptz      not null,
    reps        integer          not null,
    lapses      integer          not null,
    primary key (user_id, kind, key)
);

create index review_card_due_idx on review_card (user_id, kind, due);

create table review_log (
    user_id           uuid             not null references app_user (id) on delete cascade,
    attempt_id        uuid             not null,
    kind              text             not null check (kind in ('fact', 'pattern')),
    key               text             not null,
    grade             smallint         not null check (grade in (1, 3)),
    reviewed_at       timestamptz      not null,
    elapsed_days      integer          not null,
    stability_before  double precision,
    difficulty_before double precision,
    primary key (attempt_id, kind, key)
);

create index review_log_user_idx on review_log (user_id, reviewed_at);
