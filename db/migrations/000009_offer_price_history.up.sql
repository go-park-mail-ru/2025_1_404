SET SEARCH_PATH = kvartirum;

ALTER TABLE OfferPriceHistory
    RENAME COLUMN date TO recorded_at;

CREATE INDEX IF NOT EXISTS idx_offer_price_history_offer_id_date
    ON OfferPriceHistory (offer_id, recorded_at DESC);