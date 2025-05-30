SET SEARCH_PATH = kvartirum;

ALTER TABLE Offer
    DROP COLUMN verified,
    DROP COLUMN comment;
