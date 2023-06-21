CREATE EXTENSION IF NOT EXISTS citext;
-- SET TIME ZONE 'Europe/Moscow';

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
        id bigserial PRIMARY KEY NOT NULL,
        title text NOT NULL,
        author citext NOT NULL REFERENCES users (nickname),
        forum citext NOT NULL REFERENCES forums (slug),
        message text NOT NULL,
        slug      citext,
        created timestamp with time zone NOT NULL DEFAULT now(),
        votes integer DEFAULT 0
);

CREATE TABLE IF NOT EXISTS
    posts (
        id bigserial PRIMARY KEY NOT NULL UNIQUE,
        parent int default 0,
        author citext NOT NULL REFERENCES users(nickname),
        message text NOT NULL,
        is_edited bool DEFAULT false,
        forum citext REFERENCES forums(slug),
        thread_id bigserial REFERENCES threads(id),
        created timestamp with time zone DEFAULT now(),
        path text
);

CREATE TABLE IF NOT EXISTS
    votes (
        user_nickname citext COLLATE "ucs_basic" NOT NULL  REFERENCES users(nickname),
        thread_id bigserial NOT NULL REFERENCES threads(id),
        voice int NOT NULL,
        CONSTRAINT user_thread_key unique (user_nickname, thread_id)
);


------------
-- триггер для обновления количества голосов у потока при вставке голоса
CREATE OR REPLACE FUNCTION update_thread_on_vote_insert() RETURNS TRIGGER AS $$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_votes_on_insert
    AFTER INSERT ON votes
    FOR EACH ROW
    EXECUTE FUNCTION update_thread_on_vote_insert();

-- триггер для обновления количества голосов у потока при обновлении голоса
CREATE OR REPLACE FUNCTION update_thread_on_vote_update() RETURNS TRIGGER AS $$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice - OLD.voice
    WHERE id = NEW.thread_id;
    RETURN NEW;
END; $$ LANGUAGE plpgsql;

CREATE TRIGGER update_votes_on_update
    AFTER UPDATE ON votes
    FOR EACH ROW
    EXECUTE FUNCTION update_thread_on_vote_update();


CREATE OR REPLACE FUNCTION update_path() RETURNS TRIGGER AS $$
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.path := NEW.id::text;
    ELSE
        NEW.path := (SELECT path FROM posts WHERE id = NEW.parent) || '.' || NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_path_trigger
    BEFORE INSERT OR UPDATE ON posts
    FOR EACH ROW EXECUTE PROCEDURE update_path();
