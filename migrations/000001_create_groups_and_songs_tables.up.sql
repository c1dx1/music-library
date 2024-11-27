CREATE TABLE groups
(
    id         SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE songs
(
    id           SERIAL PRIMARY KEY,
    group_id     INT NOT NULL,
    song_name    VARCHAR(255) NOT NULL,
    release_date DATE NOT NULL,
    text         TEXT NOT NULL,
    link         VARCHAR(255) NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups (id) ON DELETE CASCADE
);
CREATE INDEX idx_songs_group_id ON songs(group_id);