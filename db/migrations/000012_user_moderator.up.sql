SET SEARCH_PATH = kvartirum;

ALTER TABLE Users
ADD role TEXT NOT NULL DEFAULT 'user'
CONSTRAINT role_len CHECK (char_length(role) <= 16);

INSERT INTO kvartirum.Users (
    first_name, last_name, email, password, token_version, role
) VALUES ('Модератор', 'Модератович', 'moderator@kvartirum.online', '$2a$12$SWBPmsNvHmSgeeRAmye17O/TQ/f4OlfcxvTm71Jefo0VvjjpYuGZS', 1, 'moderator');