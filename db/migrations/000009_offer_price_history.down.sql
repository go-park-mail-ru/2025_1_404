SET SEARCH_PATH = kvartirum;

DROP INDEX IF EXISTS idx_offer_price_history_offer_id_date;

ALTER TABLE OfferPriceHistory
    RENAME COLUMN recorded_at TO date;
