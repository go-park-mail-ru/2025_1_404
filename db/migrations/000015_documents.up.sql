SET SEARCH_PATH = kvartirum;

CREATE TABLE offer_documents (
                                 id SERIAL PRIMARY KEY,
                                 offer_id INT NOT NULL REFERENCES offers(id) ON DELETE CASCADE,
                                 url TEXT NOT NULL,
                                 name TEXT NOT NULL,
                                 created_at TIMESTAMP DEFAULT now()
);
