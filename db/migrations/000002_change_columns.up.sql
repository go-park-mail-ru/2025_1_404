SET SEARCH_PATH = kvartirum;

ALTER TABLE HousingComplex
ADD COLUMN phone_number TEXT NOT NULL
    CONSTRAINT phone_number_length CHECK (char_length(phone_number) <= 20),
ADD COLUMN address TEXT NOT NULL
    CONSTRAINT address_length CHECK (char_length(address) <= 512),
ADD COLUMN description TEXT NOT NULL
    CONSTRAINT description_length CHECK (char_length(description) <= 2048);

ALTER TABLE OFFER
ADD COLUMN area_fraction INT NOT NULL DEFAULT 0
    CONSTRAINT area_fraction_check CHECK (area_fraction >= 0 AND area_fraction <= 99),
ADD COLUMN ceiling_height_fraction INT NOT NULL DEFAULT 0
    CONSTRAINT ceiling_height_fraction_check CHECK (ceiling_height_fraction >= 0 AND ceiling_height_fraction <= 99);

ALTER TABLE MetroLine
ALTER COLUMN color SET DATA TYPE TEXT;

ALTER TABLE MetroLine
ALTER COLUMN color SET NOT NULL;
