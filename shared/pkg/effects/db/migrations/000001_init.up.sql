BEGIN;

# ─────────────────────────────── job scheduler ───────────────────────────────
CREATE TABLE import_statuses(status) AS
  VALUES
    ('NOT_STARTED'),
    ('SCAN'),
    ('METADATA'),
    ('DEDUPE'),
    ('CONVERT'),
    ('CLEANUP');
    ('DONE'),
    ('FAILED');

CREATE TABLE job_kinds(kind) AS
  VALUES
    ('SCAN'),
    ('METADATA'),
    ('METADATA_HASH'),
    ('METADATA_TIMESTAMP'),
    ('METADATA_LIVE'),
    ('METADATA_MONITOR'),
    ('DEDUPE'),
    ('CONVERT'),
    ('CONVERT_VIDEO'),
    ('CONVERT_IMAGE'),
    ('CONVERT_MONITOR'),
    ('CLEANUP');

CREATE TABLE imports (
  id         UUID     PRIMARY KEY,
  opts       TEXT     NOT NULL,
  status     TEXT     NOT NULL REFERENCES import_statuses(status),
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE jobs (
  id         UUID     PRIMARY KEY,
  import_id  UUID     NOT NULL REFERENCES imports(id),
  kind       TEXT     NOT NULL REFERENCES job_kinds(kind),
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


# ─────────────────────────────── files & paths ───────────────────────────────
CREATE TABLE file_kinds(kind) AS
  VALUES
    ('VIDEO'),
    ('IMAGE');

CREATE TABLE paths (
  id        UUID PRIMARY KEY,
  import_id UUID NOT NULL REFERENCES imports(id),
  path      TEXT NOT NULL
);

CREATE TABLE path_metadata (
  path_id   UUID     PRIMARY KEY REFERENCES paths(id),
  path      TEXT     NOT NULL,
  kind      TEXT     NOT NULL REFERENCES file_kinds(kind),
  timestamp DATETIME,
  init_hash BLOB,
  live_id   BLOB,
);

CREATE TABLE files (
  file_id   UUID     PRIMARY KEY
  import_id UUID     NOT NULL REFERENCES imports(id),
  kind      TEXT     NOT NULL REFERENCES file_kinds(kind),
  timestamp DATETIME NOT NULL,
  init_hash BLOB     NOT NULL UNIQUE,
  live_id   BLOB,
);


# ────────────────────────────────── albums ───────────────────────────────────
CREATE TABLE albums (
  id   UUID PRIMARY KEY,
  name TEXT NOT NULL
);

CREATE TABLE album_albums (
  album_id       UUID PRIMARY KEY REFERENCES albums(id),
  child_album_id UUID REFERENCES albums(id),
);

CREATE TABLE album_files (
  album_id     UUID PRIMARY KEY REFERENCES albums(id),
  file_id      UUID PRIMARY KEY REFERENCES files(id),
  is_highlight BOOL NOT NULL DEFAULT FALSE,
);

COMMIT;
