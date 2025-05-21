SET SEARCH_PATH = kvartirum;

CREATE TABLE IF NOT EXISTS kvartirum.OfferPayment (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    offer_id BIGINT NOT NULL REFERENCES Offer (id),
    type INT NOT NULL,
    yookassa_id TEXT NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    is_paid boolean NOT NULL DEFAULT false
);
