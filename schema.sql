CREATE TABLE IF NOT EXISTS blocks
(
    `index`         INTEGER PRIMARY KEY,
    `hash`          VARCHAR(64) NOT NULL,
    `previous_hash` VARCHAR(64) NOT NULL,
    `timestamp`     INTEGER DEFAULT 0,
    `data`          TEXT        NOT NULL
);