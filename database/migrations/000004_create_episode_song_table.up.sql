create table "episode_song" (
  "id" serial primary key,
  "episode_id" int not null,
  "song_id" int not null,

  constraint fk_episode_song_episode_id foreign key(episode_id) references episode(id),
  constraint fk_episode_song_song_id foreign key(song_id) references song(id)
)
