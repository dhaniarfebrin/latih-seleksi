-- Simulation (exam) attempts (docs/architecture.md sections 6-8).
-- spec_body is a frozen copy of the exam spec, so old results stay correct if
-- the embedded spec file changes later.

alter table attempt drop constraint attempt_mode_check;
alter table attempt add constraint attempt_mode_check check (mode in ('drill', 'simulasi'));

alter table attempt alter column pattern_id drop not null;
alter table attempt alter column band drop not null;
alter table attempt alter column n drop not null;

alter table attempt
    add column spec_id     text,
    add column spec_body   jsonb,
    add column track       text,
    add column sections    text[],
    add column deadline_at timestamptz,
    add column status      text check (status in ('in_progress', 'submitted', 'expired'));

alter table attempt add constraint attempt_mode_columns_check check (
    (mode = 'drill'
        and pattern_id is not null and band is not null and n is not null
        and spec_id is null and status is null)
    or
    (mode = 'simulasi'
        and spec_id is not null and spec_body is not null and track is not null
        and sections is not null and deadline_at is not null and status is not null)
);

alter table attempt_item
    add column section    text,
    add column pattern_id text;

-- seq: the client's counter; the highest seq per position wins.
alter table response
    add column seq     integer not null default 0 check (seq >= 0),
    add column flagged boolean not null default false;
