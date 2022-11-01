alter table "song" add column "released_at_episode" int;
alter table "song" add constraint fk_released_at_episode foreign key(released_at_episode) references episode(id)
