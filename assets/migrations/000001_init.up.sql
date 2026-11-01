CREATE TABLE blogger (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  first_name  TEXT     NOT NULL,
  last_name   TEXT     NOT NULL,
  email       TEXT     NOT NULL,
  age         INTEGER  NOT NULL,
  gender      TEXT     NOT NULL CHECK (gender IN ('Male', 'Female')),
  bio         TEXT     NOT NULL,
  avatar      INTEGER  NOT NULL DEFAULT 0,
  create_date DATETIME NOT NULL
);

CREATE TABLE post (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  blogger_id  INTEGER  NOT NULL REFERENCES blogger (id) ON DELETE CASCADE,
  title       TEXT     NOT NULL,
  content     TEXT     NOT NULL,
  likes       INTEGER  NOT NULL DEFAULT 0,
  dislikes    INTEGER  NOT NULL DEFAULT 0,
  create_date DATETIME NOT NULL
);

CREATE TABLE user (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  username    TEXT     NOT NULL UNIQUE,
  password    TEXT     NOT NULL,
  email       TEXT     NOT NULL UNIQUE,
  avatar      INTEGER  NOT NULL DEFAULT 0,
  create_date DATETIME
);

CREATE TABLE comment (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER  NOT NULL REFERENCES user (id) ON DELETE CASCADE,
  post_id     INTEGER  NOT NULL REFERENCES post (id) ON DELETE CASCADE,
  content     TEXT     NOT NULL,
  likes       INTEGER  NOT NULL DEFAULT 0,
  dislikes    INTEGER  NOT NULL DEFAULT 0,
  create_date DATETIME NOT NULL
);

CREATE INDEX idx_post_blogger_id ON post (blogger_id);
CREATE INDEX idx_comment_post_id ON comment (post_id);
CREATE INDEX idx_comment_user_id ON comment (user_id);
