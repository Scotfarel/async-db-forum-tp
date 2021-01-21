/*                 TABLES                   */
/* ---------------------------------------- */

/*-------USERS-------*/
CREATE TABLE users (
    id          SERIAL          PRIMARY KEY,

    email       VARCHAR(80)     NOT NULL UNIQUE,
    nickname    VARCHAR(80)     NOT NULL UNIQUE,
    fullname    VARCHAR(80)     NOT NULL,
    about       TEXT
);

CREATE INDEX ON users (lower(nickname));


/*------FORUMS-------*/
CREATE TABLE forums (
    id         SERIAL           PRIMARY KEY,

    slug       VARCHAR(80)      NOT NULL UNIQUE,
    admin      INTEGER          NOT NULL,
    title      VARCHAR(120)     NOT NULL,
    threads    INTEGER          NOT NULL DEFAULT 0,
    posts      INTEGER          NOT NULL DEFAULT 0,

    FOREIGN KEY (admin) REFERENCES "users" (id)
);

CREATE INDEX ON forums (admin);
CREATE INDEX ON forums (lower(slug));


/*----FORUMS-USERS-----*/
/* denormalization, trying to get more rps :D */
CREATE TABLE forums_users (
    id          SERIAL      PRIMARY KEY,

    user_id     INTEGER     NOT NULL,
    forum_id    INTEGER     NOT NULL,

    FOREIGN KEY (user_id) REFERENCES  "users" (id),
    FOREIGN KEY (forum_id) REFERENCES "forums" (id),
    CONSTRAINT unq_forums_users UNIQUE (forum_id, user_id)
);

CREATE INDEX ON forums_users (forum_id, user_id);


/*------THREADS------*/
CREATE TABLE threads (
    id          SERIAL                         NOT NULL PRIMARY KEY,

    author     INTEGER                         NOT NULL,
    created    TIMESTAMP (6) WITH TIME ZONE    NOT NULL,
    forum      INTEGER                         NOT NULL,
    message    TEXT                            NOT NULL,
    slug       VARCHAR(80)                     UNIQUE,
    title      VARCHAR(120)                    NOT NULL,
    votes      INTEGER                         DEFAULT 0,

    FOREIGN KEY (forum)     REFERENCES  "forums"    (id),
    FOREIGN KEY (author)    REFERENCES  "users"     (id)
);

CREATE INDEX ON threads (forum, created);
CREATE INDEX ON threads (author);
CREATE INDEX ON threads (lower(slug));


/*-------POSTS-------*/
CREATE TABLE posts (
    id          SERIAL                          NOT NULL PRIMARY KEY,

    author      INTEGER                         NOT NULL,
    forum       INTEGER                         NOT NULL,
    created     TIMESTAMP (6) WITH TIME ZONE    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    message     TEXT                            NOT NULL,
    isEdited    BOOLEAN                         DEFAULT FALSE ,
    path        INTEGER[]                       NOT NULL,
    parent      INTEGER,
    thread      INTEGER                         NOT NULL,

    FOREIGN KEY (author)   REFERENCES  "users"      (id),
    FOREIGN KEY (thread)   REFERENCES  "threads"    (id),
    FOREIGN KEY (forum)    REFERENCES  "forums"     (id)
);

CREATE INDEX ON posts ((path[1]));
CREATE INDEX ON posts (author);
CREATE INDEX ON posts (thread, created);
CREATE INDEX ON posts (thread, path);
CREATE INDEX ON posts (thread, (array_length(path, 1)));
CREATE INDEX ON posts (forum);


/*-------VOTES------*/
CREATE TABLE votes (
    id        SERIAL      NOT NULL PRIMARY KEY,

    thread    INTEGER     NOT NULL,
    author    INTEGER     NOT NULL,
    vote      INTEGER     NOT NULL,

    FOREIGN KEY (thread)    REFERENCES  "threads"   (id),
    FOREIGN KEY (author)    REFERENCES  "users"     (id)
);

CREATE INDEX ON votes (author);
CREATE INDEX ON votes (thread);
CREATE INDEX ON votes (thread, author);


/*                 TRIGGERS                 */
/* ---------------------------------------- */


/* THREAD VOTE INCREMENT FUNCTION */
CREATE OR REPLACE FUNCTION thread_vote_inc()
RETURNS TRIGGER AS
$thread_vote_inc$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.vote
    WHERE id = NEW.thread;
    RETURN NEW;
END;
$thread_vote_inc$ LANGUAGE plpgsql;

/* THREAD VOTE INCREMENT TRIGGER */
CREATE TRIGGER thread_vote_inc
    AFTER INSERT
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE thread_vote_inc();


/* UPDATE FORUMS PATH FUNCTION */
CREATE OR REPLACE FUNCTION update_forums_path()
RETURNS TRIGGER AS
$update_forums_path$
DECLARE found_id integer;
BEGIN
    NEW.path = (SELECT path FROM posts WHERE id = NEW.parent) || NEW.id;
    UPDATE forums
    SET posts = posts + 1
    WHERE id = NEW.forum;
    RETURN NEW;
END;
$update_forums_path$ LANGUAGE plpgsql;

/* UPDATE FORUMS PATH TRIGGER */
CREATE TRIGGER update_forums_path
    BEFORE INSERT
    ON posts
    FOR EACH ROW
EXECUTE PROCEDURE update_forums_path();


/* ADD COUNT UPDATE FUNCTION */
CREATE OR REPLACE FUNCTION add_count_update()
RETURNS TRIGGER AS
$add_count_update$
BEGIN
    UPDATE threads
    SET votes = votes - OLD.vote + NEW.vote
    WHERE id = NEW.thread;
    RETURN NEW;
END;
$add_count_update$ LANGUAGE plpgsql;

/* ADD COUNT UPDATE TRIGGER */
CREATE TRIGGER add_count_update
    BEFORE UPDATE
    ON votes
    FOR EACH ROW
EXECUTE PROCEDURE add_count_update();



/* ADD USER TO FORUM FUNCTION */
CREATE OR REPLACE FUNCTION add_forum_user()
RETURNS TRIGGER AS
$add_forum_user$
BEGIN
    INSERT INTO forums_users (user_id, forum_id)
    VALUES (NEW.author, NEW.forum)
    ON CONFLICT DO NOTHING;

    UPDATE forums
    SET threads = threads + 1
    WHERE id = NEW.forum;

    RETURN NEW;
END;
$add_forum_user$ LANGUAGE plpgsql;

/* ADD USER TO FORUM TRIGGER */
CREATE TRIGGER add_forum_user
    AFTER INSERT
    ON threads
    FOR EACH ROW
EXECUTE PROCEDURE add_forum_user();
