-- User reports on a question of a closed simulation (FR-03). One report per
-- (attempt, pos); a new report replaces the old one. Rows go with the attempt,
-- and so with the user (DELETE /v1/me). Read through SQL only: no admin UI yet.

create table question_report (
    attempt_id uuid        not null,
    pos        smallint    not null,
    reason     text        not null check (reason in ('kunci_salah', 'soal_ambigu', 'typo', 'lainnya')),
    note       text        not null default '' check (length(note) <= 500),
    created_at timestamptz not null default now(),
    primary key (attempt_id, pos),
    foreign key (attempt_id, pos) references attempt_item (attempt_id, pos) on delete cascade
);
