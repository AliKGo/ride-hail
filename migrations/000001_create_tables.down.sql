begin;

drop table if exists location_history;
drop table if exists driver_sessions;
drop index if exists idx_drivers_status;
drop table if exists drivers;
drop table if exists driver_status;

commit;
