create table "episode" (
  "id" serial primary key,
  "episode" integer not null,
  "episode_name" varchar(255),
  "episode_date" timestamptz not null
);
