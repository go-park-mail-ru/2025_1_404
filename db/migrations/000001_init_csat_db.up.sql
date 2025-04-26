-- CREATE SCHEMA csat;

SET SEARCH_PATH = csat;

CREATE TABLE IF NOT EXISTS CSAT (
    id BIGINT GENERATED ALWAYS AS IDENTITY
    PRIMARY KEY,
    event TEXT NOT NULL
    constraint event_length CHECK (char_length(event) <= 32)
);

CREATE TABLE IF NOT EXISTS Question (
    id BIGINT GENERATED ALWAYS AS IDENTITY
    PRIMARY KEY,
    text TEXT NOT NULL
    constraint text_length CHECK (char_length(text) <= 512),
    csat_id BIGINT
    DEFAULT NULL
    REFERENCES CSAT (id)
    ON DELETE cascade
    ON UPDATE cascade
);

CREATE TABLE IF NOT EXISTS Answer (
    id BIGINT GENERATED ALWAYS AS IDENTITY
    PRIMARY KEY,
    rating int NOT NULL
    CONSTRAINT rating CHECK (rating >= 0 and rating <= 100),
    question_id BIGINT
    DEFAULT NULL
    REFERENCES Question (id)
    ON DELETE cascade
    ON UPDATE cascade
);
