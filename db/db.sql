CREATE EXTENSION IF NOT EXISTS citext;
SET TIME ZONE 'Europe/Moscow';

CREATE TABLE IF NOT EXISTS
    users (
        nickname citext COLLATE "ucs_basic" NOT NULL UNIQUE PRIMARY KEY,
        fullname text NOT NULL,
        about text NOT NULL,
        email citext NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS
    forums (
        title text NOT NULL,
        "user" citext NOT NULL REFERENCES users (nickname),
        slug citext NOT NULL PRIMARY KEY,
        posts int DEFAULT 0,
        threads int DEFAULT 0
);

CREATE TABLE IF NOT EXISTS
    threads (
        title text NOT NULL,
        author citext NOT NULL REFERENCES users (nickname),
        message text NOT NULL,
        created timestamp NOT NULL
);

