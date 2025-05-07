SET SEARCH_PATH = kvartirum;

ALTER TABLE HousingComplex
ADD COLUMN metro_station_id INT NOT NULL,
ADD CONSTRAINT fk_metro_station
  FOREIGN KEY (metro_station_id)
  REFERENCES MetroStation (id)
  ON DELETE CASCADE
  ON UPDATE CASCADE;

ALTER TABLE Offer
ADD COLUMN longitude TEXT NOT NULL
    CONSTRAINT longitude_length CHECK (char_length(longitude) <= 32),
ADD COLUMN latitude TEXT NOT NULL
    CONSTRAINT latitude_length CHECK (char_length(latitude) <= 32);


ALTER TABLE Image
DROP CONSTRAINT IF EXISTS uuid_length;

-- Добавляем новое ограничение с большей длиной
ALTER TABLE Image
ADD CONSTRAINT uuid_length CHECK (char_length(uuid) <= 256);

CREATE TABLE IF NOT EXISTS Likes(
     id BIGINT GENERATED ALWAYS AS IDENTITY
    PRIMARY KEY,
    user_id BIGINT NOT NULL
    REFERENCES Users (id)
    ON DELETE cascade
    ON UPDATE cascade,
    offer_id BIGINT NOT NULL
    REFERENCES Offer (id)
    ON DELETE cascade
    ON UPDATE cascade,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
); 