-- Drill attempts now belong to a user. not valid: anonymous drill rows from
-- before this migration are left as they are (unreachable, user_id null); every
-- new row is checked. Same pattern as attempt_simulasi_user_check in 0003.

alter table attempt add constraint attempt_drill_user_check
    check (mode <> 'drill' or user_id is not null) not valid;
