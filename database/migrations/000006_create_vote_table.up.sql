create table "vote" (
  "id" serial primary key,
  "ip_address" varchar(15) not null,
  "episode_song_id" int not null,
  "created_at" timestamptz default now() not null,

  constraint fk_vote_episode_song_id foreign key(episode_song_id) references episode_song(id)
)
