SET
SEARCH_PATH = kvartirum;

CREATE ROLE user_role
    WITH
    LOGIN
    PASSWORD 'password'
    NOINHERIT;

-- Даем право использовать схему kvartirum
GRANT
USAGE
ON
SCHEMA
kvartirum TO user_role;

-- Даем право подключаться к базе данных kvartirum
GRANT CONNECT
ON DATABASE kvartirum TO user_role;

-- Даем права только на SELECT
GRANT
SELECT
ON TABLE
    HousingComplexImages,
    HousingComplexReview,
    HousingComplex,
    HousingComplexClass,
    OfferRenovation,
    OfferStatus,
    PropertyType,
    PurchaseType,
    RentType,
    MetroStation,
    MetroLine,
    OfferType,
    ChatMessage,
    Chat,
    OfferImages,
    UserReview,
    UserNotification,
    UserOfferFavourites,
    OfferPriceHistory,
    Offer,
    Users,
    Image,
    Views
    TO user_role;

-- Даем права на INSERT, UPDATE, DELETE
GRANT
INSERT,
UPDATE,
DELETE
ON TABLE
    ChatMessage,
    Chat,
    OfferImages,
    UserReview,
    UserNotification,
    UserOfferFavourites,
    OfferPriceHistory,
    Offer,
    Users,
    Image,
    Views
    TO user_role;
