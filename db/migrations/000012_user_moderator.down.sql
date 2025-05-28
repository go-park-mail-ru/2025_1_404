SET SEARCH_PATH = kvartirum;

DELETE FROM kvartirum.Users
WHERE email = 'moderator@kvartirum.online';

ALTER TABLE Users
DROP COLUMN role;