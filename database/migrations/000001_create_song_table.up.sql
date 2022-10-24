create table "song" (
  "id" serial primary key,
  "song_name_jp" varchar(255) not null,
  "song_name_en" varchar(255),
  "artist_name_jp" varchar(255) not null,
  "artist_name_en" varchar(255),
  "image_url" text not null
);
