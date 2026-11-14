-- Drill attempts, frozen items and responses (docs/architecture.md section 6).
-- Single local user for now: no user_id column; a later migration adds it.

create table attempt (
    id           uuid primary key default gen_random_uuid(),
    mode         text        not null check (mode in ('drill')),
    pattern_id   text        not null,
    band         smallint    not null check (band between 1 and 3),
    n            smallint    not null check (n between 1 and 20),
    started_at   timestamptz not null default now(),
    submitted_at timestamptz
);

create index attempt_started_at_idx on attempt (started_at desc);

-- item_snapshot is the frozen question (stem, options, answer key, steps).
create table attempt_item (
    attempt_id    uuid     not null references attempt (id) on delete cascade,
    pos           smallint not null check (pos >= 1),
    seed          bigint   not null,
    item_snapshot jsonb    not null,
    primary key (attempt_id, pos)
);

create table response (
    attempt_id       uuid        not null,
    pos              smallint    not null,
    answer           text        not null check (answer in ('A', 'B', 'C', 'D', 'E')),
    correct          boolean     not null,
    time_ms          integer     not null check (time_ms >= 0),
    misconception_id text,
    answered_at      timestamptz not null default now(),
    primary key (attempt_id, pos),
    foreign key (attempt_id, pos) references attempt_item (attempt_id, pos) on delete cascade
);
