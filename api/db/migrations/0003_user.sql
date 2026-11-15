-- Users (docs/architecture.md section 11). Only the auth provider's subject is
-- stored: no email, no name. Simulation attempts belong to a user; drill
-- attempts stay anonymous (user_id null).

create table app_user (
    id         uuid primary key default gen_random_uuid(),
    auth_sub   text        not null unique,
    created_at timestamptz not null default now()
);

alter table attempt
    add column user_id uuid references app_user (id) on delete cascade;

-- not valid: simulation rows from before this migration have no owner and are
-- left as they are; every new row is checked.
alter table attempt add constraint attempt_simulasi_user_check
    check (mode <> 'simulasi' or user_id is not null) not valid;

create index attempt_user_id_idx on attempt (user_id) where user_id is not null;
