SET SEARCH_PATH = kvartirum;

ALTER TABLE offerpayment
DROP CONSTRAINT offerpayment_offer_id_fkey,
ADD CONSTRAINT offerpayment_offer_id_fkey
FOREIGN KEY (offer_id) REFERENCES offer(id)
ON DELETE CASCADE ON UPDATE CASCADE;