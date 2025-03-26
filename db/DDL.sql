-- Сброс схемы
DROP SCHEMA kvartirum CASCADE;
CREATE SCHEMA kvartirum;

-- Указываем схему
SET SEARCH_PATH = kvartirum;

-- Создание таблицы Image
DROP TABLE IF EXISTS Image;
CREATE TABLE IF NOT EXISTS Image (
                                     id BIGINT GENERATED ALWAYS AS IDENTITY
                                     PRIMARY KEY,
                                     uuid TEXT NOT NULL
                                     DEFAULT NULL
                                     CONSTRAINT uuid_length CHECK (char_length(uuid) <= 64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы Users
DROP TABLE IF EXISTS Users;
CREATE TABLE IF NOT EXISTS Users (
                                     id BIGINT GENERATED ALWAYS AS IDENTITY
                                     PRIMARY KEY,
                                     image_id BIGINT
                                     DEFAULT NULL
                                     REFERENCES Image (id)
    ON DELETE cascade
    ON UPDATE cascade,
    first_name TEXT NOT NULL
    CONSTRAINT first_name_length CHECK (char_length(first_name) <= 32),
    last_name TEXT NOT NULL
    CONSTRAINT last_name_length CHECK (char_length(last_name) <= 32),
    email TEXT NOT NULL
    CONSTRAINT email_length CHECK (char_length(email) <= 32),
    password TEXT NOT NULL
    CONSTRAINT password_length CHECK (char_length(password) <= 64),
    last_notification_id INTEGER
    DEFAULT NULL,
    token_version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы OfferType
DROP TABLE IF EXISTS OfferType;
CREATE TABLE IF NOT EXISTS OfferType (
                                         id INT GENERATED ALWAYS AS IDENTITY
                                         PRIMARY KEY,
                                         name TEXT NOT NULL
                                         CONSTRAINT name_length CHECK (char_length(name) <= 32)
    );

-- Создание таблицы MetroLine
DROP TABLE IF EXISTS MetroLine;
CREATE TABLE IF NOT EXISTS MetroLine (
                                         id INT GENERATED ALWAYS AS IDENTITY
                                         PRIMARY KEY,
                                         name TEXT NOT NULL
                                         CONSTRAINT name_length CHECK (char_length(name) <= 32),
    color INT NOT NULL
    );

-- Создание таблицы MetroStation
DROP TABLE IF EXISTS MetroStation;
CREATE TABLE IF NOT EXISTS MetroStation (
                                            id INT GENERATED ALWAYS AS IDENTITY
                                            PRIMARY KEY,
                                            metro_line_id INT NOT NULL
                                            REFERENCES MetroLine (id)
    ON DELETE cascade
    ON UPDATE cascade,
    name TEXT NOT NULL
    CONSTRAINT name_length CHECK (char_length(name) <= 32)
    );

-- Создание таблицы RentType ()
DROP TABLE IF EXISTS RentType;
CREATE TABLE IF NOT EXISTS RentType (
                                        id INT GENERATED ALWAYS AS IDENTITY
                                        PRIMARY KEY,
                                        name TEXT NOT NULL
                                        CONSTRAINT name_length CHECK (char_length(name) <= 32)
    );

-- Создание таблицы PurchaseType ()
DROP TABLE IF EXISTS PurchaseType;
CREATE TABLE IF NOT EXISTS PurchaseType (
                                            id INT GENERATED ALWAYS AS IDENTITY
                                            PRIMARY KEY,
                                            name TEXT NOT NULL
                                            CONSTRAINT name_length CHECK (char_length(name) <= 32)
    );

-- Создание таблицы PropertyType ()
DROP TABLE IF EXISTS PropertyType;
CREATE TABLE IF NOT EXISTS PropertyType (
                                            id INT GENERATED ALWAYS AS IDENTITY
                                            PRIMARY KEY,
                                            name TEXT NOT NULL
                                            CONSTRAINT name_length CHECK (char_length(name) <= 32)
    );

-- Создание таблицы OfferStatus ()
DROP TABLE IF EXISTS OfferStatus;
CREATE TABLE IF NOT EXISTS OfferStatus (
                                           id INT GENERATED ALWAYS AS IDENTITY
                                           PRIMARY KEY,
                                           name TEXT NOT NULL
                                           CONSTRAINT name_length CHECK (char_length(name) <= 32)
    );

-- Создание таблицы OfferRenovation ()
DROP TABLE IF EXISTS OfferRenovation;
CREATE TABLE IF NOT EXISTS OfferRenovation (
                                               id INT GENERATED ALWAYS AS IDENTITY
                                               PRIMARY KEY,
                                               name TEXT NOT NULL
                                               CONSTRAINT name_length CHECK (char_length(name) <= 32)
    );

-- Создание таблицы HousingComplexClass()
DROP TABLE IF EXISTS HousingComplexClass;
CREATE TABLE IF NOT EXISTS HousingComplexClass (
                                                   id INT GENERATED ALWAYS AS IDENTITY
                                                   PRIMARY KEY,
                                                   name TEXT NOT NULL
                                                   CONSTRAINT name_length CHECK (char_length(name) <= 16)
    );

-- Создание таблицы HousingComplex
DROP TABLE IF EXISTS HousingComplex;
CREATE TABLE IF NOT EXISTS HousingComplex (
                                              id BIGINT GENERATED ALWAYS AS IDENTITY
                                              PRIMARY KEY,
                                              class_id INT NOT NULL
                                              REFERENCES HousingComplexClass (id)
    ON DELETE cascade
    ON UPDATE cascade,
    name TEXT NOT NULL
    CONSTRAINT name_length CHECK (char_length(name) <= 32),
    developer TEXT NOT NULL
    CONSTRAINT developer_length CHECK (char_length(developer) <= 32),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы HousingComplexReview
DROP TABLE IF EXISTS HousingComplexReview;
CREATE TABLE IF NOT EXISTS HousingComplexReview (
                                                    id BIGINT GENERATED ALWAYS AS IDENTITY
                                                    PRIMARY KEY,
                                                    user_id BIGINT NOT NULL
                                                    REFERENCES Users (id)
    ON DELETE cascade
    ON UPDATE cascade,
    housing_complex_id BIGINT NOT NULL
    REFERENCES HousingComplex (id)
    ON DELETE cascade
    ON UPDATE cascade,
    rating INT NOT NULL
    CONSTRAINT rating CHECK (rating >= 0 AND rating <= 5),
    comment TEXT
    DEFAULT NULL
    CONSTRAINT comment_length CHECK (char_length(comment) <= 128),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы HousingComplexImages
DROP TABLE IF EXISTS HousingComplexImages;
CREATE TABLE IF NOT EXISTS HousingComplexImages (
                                                    id BIGINT GENERATED ALWAYS AS IDENTITY
                                                    PRIMARY KEY,
                                                    housing_complex_id BIGINT NOT NULL
                                                    REFERENCES HousingComplex (id)
    ON DELETE cascade
    ON UPDATE cascade,
    image_id BIGINT NOT NULL
    REFERENCES Image (id)
    ON DELETE cascade
    ON UPDATE cascade,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы Offer
DROP TABLE IF EXISTS Offer;
CREATE TABLE IF NOT EXISTS Offer (
                                     id BIGINT GENERATED ALWAYS AS IDENTITY
                                     PRIMARY KEY,
                                     seller_id BIGINT NOT NULL
                                     REFERENCES Users (id)
    ON DELETE cascade
    ON UPDATE cascade,
    offer_type_id INT NOT NULL
    REFERENCES OfferType (id)
    ON DELETE cascade
    ON UPDATE cascade,
    metro_station_id INT
    REFERENCES MetroStation (id)
    ON DELETE cascade
    ON UPDATE cascade
    DEFAULT NULL,
    rent_type_id INT
    REFERENCES RentType (id)
    ON DELETE cascade
    ON UPDATE cascade
    DEFAULT NULL,
    purchase_type_id INT
    REFERENCES PurchaseType (id)
    ON DELETE cascade
    ON UPDATE cascade
    DEFAULT NULL,
    property_type_id INT NOT NULL
    REFERENCES PropertyType (id)
    ON DELETE cascade
    ON UPDATE cascade,
    offer_status_id INT NOT NULL
    REFERENCES OfferStatus (id)
    ON DELETE cascade
    ON UPDATE cascade,
    renovation_id INT NOT NULL
    REFERENCES OfferRenovation (id)
    ON DELETE cascade
    ON UPDATE cascade,
    complex_id INT
    REFERENCES HousingComplex (id)
    ON DELETE cascade
    ON UPDATE cascade
    DEFAULT NULL,
    price INT NOT NULL
    CONSTRAINT price CHECK (price >= 0),
    description TEXT
    DEFAULT NULL
    CONSTRAINT description_length CHECK (char_length(description) <= 512),
    floor INT NOT NULL
    CONSTRAINT floor CHECK (floor >= 0 AND floor <= 100),
    total_floors INT NOT NULL
    CONSTRAINT total_floors CHECK (total_floors >= 0 AND total_floors <= 100),
    rooms INT NOT NULL
    CONSTRAINT rooms CHECK (rooms >= 0 AND rooms <= 100),
    address text
    CONSTRAINT address_length CHECK (char_length(address) <= 512),
    flat INT NOT NULL
    CONSTRAINT flat CHECK (flat >= 0 AND flat <= 1000),
    area INT NOT NULL
    CONSTRAINT area CHECK (area >= 0 AND area <= 1000),
    ceiling_height INT NOT NULL
    CONSTRAINT ceiling_height CHECK (ceiling_height >= 0 AND ceiling_height <= 100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы OfferPriceHistory
DROP TABLE IF EXISTS OfferPriceHistory;
CREATE TABLE IF NOT EXISTS OfferPriceHistory (
                                                 id BIGINT GENERATED ALWAYS AS IDENTITY
                                                 PRIMARY KEY,
                                                 offer_id BIGINT NOT NULL
                                                 REFERENCES Offer (id)
    ON DELETE cascade
    ON UPDATE cascade,
    price INT NOT NULL
    CONSTRAINT price CHECK (price >= 0),
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы UserOfferFavourites
DROP TABLE IF EXISTS UserOfferFavourites;
CREATE TABLE IF NOT EXISTS UserOfferFavourites (
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

-- Создание таблицы UserNotification
DROP TABLE IF EXISTS UserNotification;
CREATE TABLE IF NOT EXISTS UserNotification (
                                                id BIGINT GENERATED ALWAYS AS IDENTITY
                                                PRIMARY KEY,
                                                user_id BIGINT NOT NULL
                                                REFERENCES Users (id)
    ON DELETE cascade
    ON UPDATE cascade,
    message TEXT
    CONSTRAINT message_length CHECK (char_length(message) <= 64),
    redirect_uri TEXT
    DEFAULT NULL
    CONSTRAINT redirect_uri_length CHECK (char_length(redirect_uri) <= 64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы UserReview
DROP TABLE IF EXISTS UserReview;
CREATE TABLE IF NOT EXISTS UserReview (
                                          id BIGINT GENERATED ALWAYS AS IDENTITY
                                          PRIMARY KEY,
                                          reviewer_id BIGINT NOT NULL
                                          REFERENCES Users (id)
    ON DELETE cascade
    ON UPDATE cascade,
    seller_id BIGINT NOT NULL
    REFERENCES Users (id)
    ON DELETE cascade
    ON UPDATE cascade,
    rating INT NOT NULL
    CONSTRAINT rating CHECK (rating >= 0 AND rating <= 5),
    comment TEXT
    DEFAULT NULL
    CONSTRAINT comment_length CHECK (char_length(comment) <= 128),
    redirect_uri TEXT
    DEFAULT NULL
    CONSTRAINT redirect_uri_length CHECK (char_length(redirect_uri) <= 64),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы OfferImages
DROP TABLE IF EXISTS OfferImages;
CREATE TABLE IF NOT EXISTS OfferImages (
                                           id BIGINT GENERATED ALWAYS AS IDENTITY
                                           PRIMARY KEY,
                                           offer_id BIGINT NOT NULL
                                           REFERENCES Offer (id)
    ON DELETE cascade
    ON UPDATE cascade,
    image_id BIGINT NOT NULL
    REFERENCES Image (id)
    ON DELETE cascade
    ON UPDATE cascade,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы Chat
DROP TABLE IF EXISTS Chat;
CREATE TABLE IF NOT EXISTS Chat (
                                    id BIGINT GENERATED ALWAYS AS IDENTITY
                                    PRIMARY KEY,
                                    offer_id BIGINT NOT NULL
                                    REFERENCES Offer (id)
    ON DELETE cascade
    ON UPDATE cascade,
    customer_id BIGINT NOT NULL
    REFERENCES Users (id)
    ON DELETE cascade
    ON UPDATE cascade,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Создание таблицы ChatMessage
DROP TABLE IF EXISTS ChatMessage;
CREATE TABLE IF NOT EXISTS ChatMessage (
                                           id BIGINT GENERATED ALWAYS AS IDENTITY
                                           PRIMARY KEY,
                                           chat_id BIGINT NOT NULL
                                           REFERENCES Chat (id)
    ON DELETE cascade
    ON UPDATE cascade,
    user_id BIGINT NOT NULL
    REFERENCES Users (id)
    ON DELETE cascade
    ON UPDATE cascade,
    message TEXT NOT NULL
    CONSTRAINT message_length CHECK (char_length(message) <= 128),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
    );

-- Функция для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Добавляем триггер для таблицы Users
CREATE TRIGGER set_updated_at_users
    BEFORE UPDATE ON Users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Добавляем триггер для таблицы Offer
CREATE TRIGGER set_updated_at_offer
    BEFORE UPDATE ON Offer
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Добавляем триггер для таблицы Chat
CREATE TRIGGER set_updated_at_chat
    BEFORE UPDATE ON Chat
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Добавляем триггер для таблицы ChatMessage
CREATE TRIGGER set_updated_at_chat_message
    BEFORE UPDATE ON ChatMessage
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();