CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS Teams (
  id      BIGSERIAL NOT NULL PRIMARY KEY,
  slug    CITEXT UNIQUE,
  title   TEXT UNIQUE,
  country CHAR(2) DEFAULT 'NO'
);

CREATE TABLE IF NOT EXISTS Players (
  id       BIGSERIAL NOT NULL PRIMARY KEY,
  fullname TEXT,
  nickname CITEXT    NOT NULL UNIQUE,
  country  CHAR(2) DEFAULT 'NO',
  teamSlug CITEXT NOT NULL,
  team     BIGINT REFERENCES Teams (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Games (
  id    BIGSERIAL   NOT NULL PRIMARY KEY,
  slug  CITEXT      NOT NULL UNIQUE,
  title TEXT UNIQUE NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS games_slug_idx
  ON Games (slug);

CREATE TABLE IF NOT EXISTS Tournaments (
  id     BIGSERIAL NOT NULL PRIMARY KEY,
  title  TEXT UNIQUE,
  slug   CITEXT    NOT NULL UNIQUE,
  stages JSONB
);

CREATE TABLE IF NOT EXISTS Matches (
  id         BIGSERIAL                NOT NULL PRIMARY KEY,
  title      TEXT UNIQUE,
  slug       CITEXT                   NOT NULL UNIQUE,
  team1      BIGINT                   NOT NULL REFERENCES Teams (id) ON DELETE CASCADE,
  team2      BIGINT                   NOT NULL REFERENCES Teams (id) ON DELETE CASCADE,
  tournament BIGINT REFERENCES Tournaments (id) ON DELETE CASCADE,
  date       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE SUBSCRIPTION sub CONNECTION 'hostaddr=192.168.1.30 port=5432 dbname=shard1-master user=postgres password=yCJvclmS6jZag68j' PUBLICATION pub;