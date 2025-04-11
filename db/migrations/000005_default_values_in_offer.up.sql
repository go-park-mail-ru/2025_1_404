SET SEARCH_PATH = kvartirum;

ALTER TABLE Offer
ALTER COLUMN price SET DEFAULT 0,
ALTER COLUMN floor SET DEFAULT 1,
ALTER COLUMN total_floors SET DEFAULT 1,
ALTER COLUMN rooms SET DEFAULT 1,
ALTER COLUMN flat SET DEFAULT 1,
ALTER COLUMN area SET DEFAULT 1,
ALTER COLUMN ceiling_height SET DEFAULT 3,
ALTER COLUMN address SET DEFAULT NULL;

CREATE OR REPLACE FUNCTION set_offer_defaults()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.offer_type_id IS NULL THEN
        NEW.offer_type_id := (SELECT id FROM kvartirum.OfferType WHERE name = 'Продажа' LIMIT 1);
    END IF;
    
    IF NEW.property_type_id IS NULL THEN
        NEW.property_type_id := (SELECT id FROM kvartirum.PropertyType WHERE name = 'Квартира' LIMIT 1);
    END IF;
    
    IF NEW.offer_status_id IS NULL THEN
        NEW.offer_status_id := (SELECT id FROM kvartirum.OfferStatus WHERE name = 'Черновик' LIMIT 1);
    END IF;
    
    IF NEW.renovation_id IS NULL THEN
        NEW.renovation_id := (SELECT id FROM kvartirum.OfferRenovation WHERE name = 'Черновая отделка' LIMIT 1);
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_offer_defaults
BEFORE INSERT ON Offer
FOR EACH ROW
EXECUTE FUNCTION set_offer_defaults();