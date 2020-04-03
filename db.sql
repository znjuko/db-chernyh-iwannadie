DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS threadUF CASCADE;
DROP TABLE IF EXISTS threads;
DROP TABLE IF EXISTS messageTU CASCADE;
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS counter;
DROP TABLE IF EXISTS voteThreads;

DROP FUNCTION IF EXISTS user_counter();
DROP FUNCTION IF EXISTS forum_counter();
DROP FUNCTION IF EXISTS thread_counter();
DROP FUNCTION IF EXISTS message_counter();

DROP TRIGGER IF EXISTS add_user ON users;
DROP TRIGGER IF EXISTS add_forum ON forums;
DROP TRIGGER IF EXISTS add_thread ON threads;
DROP TRIGGER IF EXISTS add_message ON messages;
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE users
(
    u_id     BIGSERIAL PRIMARY KEY,
    nickname CITEXT COLLATE "C"  NOT NULL UNIQUE,
    fullname VARCHAR(100) NOT NULL,
    email    CITEXT       NOT NULL UNIQUE,
    about    TEXT
);

CREATE TABLE forums
(
    f_id  BIGSERIAL PRIMARY KEY,
    u_id  BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    slug  CITEXT UNIQUE NOT NULL,
    title TEXT
);

CREATE TABLE threadUF
(
    t_id BIGSERIAL PRIMARY KEY,
    slug    CITEXT UNIQUE,
    u_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    f_id BIGINT NOT NULL REFERENCES forums ON DELETE CASCADE
--     CONSTRAINT unq_forum_slug UNIQUE(slug)
);

CREATE TABLE threads
(
    t_id    BIGINT    NOT NULL REFERENCES threadUF ON DELETE CASCADE,
    date    TIMESTAMP WITH TIME ZONE,
    message TEXT,
    title   TEXT,
    votes   BIGINT DEFAULT 0
);

CREATE TABLE voteThreads
(
    t_id BIGINT NOT NULL REFERENCES threadUF ON DELETE CASCADE ,
    u_id BIGINT NOT NULL,
    counter INT DEFAULT 0
);

CREATE TABLE messageTU
(
    m_id BIGSERIAL PRIMARY KEY,
    u_id BIGINT NOT NULL REFERENCES users ON DELETE CASCADE,
    t_id BIGINT NOT NULL REFERENCES threadUF ON DELETE CASCADE
);

CREATE TABLE messages
(
    m_id    BIGINT  NOT NULL REFERENCES messageTU ON DELETE CASCADE,
    date    TIMESTAMP WITH TIME ZONE,
    message TEXT,
    edit    BOOLEAN DEFAULT false,
    parent  BIGINT
);

CREATE TABLE counter
(
    users    BIGINT DEFAULT 0,
    forums   BIGINT DEFAULT 0,
    threads  BIGINT DEFAULT 0,
    messages BIGINT DEFAULT 0
);

CREATE FUNCTION user_counter()
    RETURNS TRIGGER
AS
$$
BEGIN
    UPDATE counter SET users = users + 1;
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE FUNCTION forum_counter()
    RETURNS TRIGGER
AS
$$
BEGIN
    UPDATE counter SET forums = forums + 1;
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE FUNCTION thread_counter()
    RETURNS TRIGGER
AS
$$
BEGIN
    UPDATE counter SET threads = threads + 1;
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE FUNCTION message_counter()
    RETURNS TRIGGER
AS
$$
BEGIN
    UPDATE counter SET messages = messages + 1;
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;

CREATE TRIGGER add_user
    BEFORE INSERT
    ON users
    FOR EACH ROW
EXECUTE PROCEDURE user_counter();

CREATE TRIGGER add_forum
    BEFORE INSERT
    ON forums
    FOR EACH ROW
EXECUTE PROCEDURE forum_counter();

CREATE TRIGGER add_thread
    BEFORE INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE thread_counter();

CREATE TRIGGER add_message
    BEFORE INSERT
    ON messages
    FOR EACH ROW
EXECUTE PROCEDURE message_counter();

INSERT INTO counter (users, forums, threads, messages)
VALUES (0, 0, 0, 0);


