SET SEARCH_PATH = kvartirum;

ALTER TABLE kvartirum.Offer
    ADD COLUMN promotes_until TIMESTAMP WITH TIME ZONE;
