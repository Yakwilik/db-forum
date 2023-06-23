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
        path bigint[] DEFAULT ARRAY []::INTEGER[]
);


CREATE TABLE IF NOT EXISTS
    votes (
        user_nickname citext COLLATE "ucs_basic" NOT NULL  REFERENCES users(nickname),
        thread_id bigserial NOT NULL REFERENCES threads(id),
        voice int NOT NULL,
        CONSTRAINT user_thread_key unique (user_nickname, thread_id)
);

CREATE UNLOGGED TABLE IF NOT EXISTS user_forums
(

    nickname citext COLLATE "ucs_basic" NOT NULL REFERENCES users (nickname),
    forum    citext                     NOT NULL REFERENCES forums (slug),
    fullname text,
    about    text,
    email    citext,
    CONSTRAINT user_forum_key unique (nickname, forum)
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


-- триггер для создания пути в дереве комментариев
CREATE OR REPLACE FUNCTION update_path() RETURNS TRIGGER AS $$
BEGIN
    NEW.path = (SELECT path FROM posts WHERE id = new.parent) || new.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_path_trigger
    BEFORE INSERT OR UPDATE ON posts
    FOR EACH ROW EXECUTE PROCEDURE update_path();


-- обновление числа постов в форуме
CREATE OR REPLACE FUNCTION increment_forum_posts() RETURNS TRIGGER AS $$
BEGIN
    UPDATE forums
    SET posts = posts + 1
    WHERE slug = NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER increment_forum_posts_trigger
    AFTER INSERT ON posts
    FOR EACH ROW EXECUTE PROCEDURE increment_forum_posts();


-- обновление числа веток в форуме
CREATE OR REPLACE FUNCTION increment_forum_threads() RETURNS TRIGGER AS $$
BEGIN
    UPDATE forums
    SET threads = threads + 1
    WHERE slug = NEW.forum;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER increment_forum_threads_trigger
    AFTER INSERT ON threads
    FOR EACH ROW EXECUTE PROCEDURE increment_forum_threads();

-- обновление состояния сообщения при изменении
CREATE OR REPLACE FUNCTION update_is_edited() RETURNS TRIGGER AS $$
BEGIN
    IF OLD.message <> NEW.message THEN
        NEW.is_edited = true;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_update_message
    BEFORE UPDATE OF message ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_is_edited();


--- Пользователи форума

CREATE OR REPLACE FUNCTION function_update_user_forum()
    RETURNS TRIGGER AS

-- info posts, threads
$$
DECLARE
    _nickname citext;
    _fullname text;
    _about    text;
    _email    citext;
BEGIN
    SELECT u.nickname, u.fullname, u.about, u.email
    FROM users u
    WHERE u.nickname = NEW.author
    INTO _nickname, _fullname, _about, _email;

    INSERT INTO user_forums (nickname, fullname, about, email, forum)
    VALUES (_nickname, _fullname, _about, _email, NEW.forum)
    ON CONFLICT DO NOTHING;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_user_forum
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE function_update_user_forum();

CREATE TRIGGER update_users_forum
    AFTER INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE function_update_user_forum();



------------ Индексы

CREATE INDEX IF NOT EXISTS forum_slug_hash ON forums USING hash (slug);
CREATE INDEX IF NOT EXISTS forum_user_hash ON forums USING hash ("user");


CREATE INDEX IF NOT EXISTS idx_posts_id ON posts USING hash (id);
CREATE INDEX idx_posts_id_btree ON posts using btree (id);
CREATE INDEX IF NOT EXISTS post_thread ON posts USING hash (thread_id);
CREATE INDEX IF NOT EXISTS post_parent ON posts (thread_id, id, (path[1]), parent);
CREATE INDEX IF NOT EXISTS post_path_1_path ON posts ((path[1]), path);
CREATE INDEX IF NOT EXISTS post_thread_path ON posts (thread_id, path);


CREATE INDEX idx_posts_created ON posts using btree (created);
CREATE INDEX idx_posts_path ON posts using btree (path);

CREATE INDEX idx_forums_slug ON forums USING HASH (slug);
CREATE INDEX idx_threads_id_hash ON threads USING hash (id);
CREATE INDEX idx_threads_created ON threads USING btree (created);
CREATE INDEX idx_users_nickname ON users USING HASH (nickname);
CREATE INDEX idx_users_nickname_btree ON users USING btree (nickname varchar_pattern_ops);
CREATE INDEX idx_users_email ON users USING HASH (email);