SET SEARCH_PATH = kvartirum;

UPDATE OfferStatus SET name = 'Черновик' WHERE id = 1;
UPDATE OfferStatus SET name = 'Активный' WHERE id = 2;
UPDATE OfferStatus SET name = 'Завершенный' WHERE id = 3;