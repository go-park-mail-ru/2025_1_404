# [Код Go](https://github.com/go-park-mail-ru/2025_1_404/tree/dev_db)

# ER диаграмма
```mermaid
erDiagram
    Offer {
        id BIGINT
        seller_id BIGINT
        offer_type_id INT
        metro_station_id INT
        rent_type_id INT
        purchase_type_id INT 
        property_type_id INT
        offer_status_id INT
        renovation_id INT
        complex_id INT
        price INT
        description TEXT
        floor INT
        total_floors INT
        rooms INT
        address TEXT
        flat INT
        area INT
        ceiling_height INT
        created_at TIMESTAMP
        updated_at TIMESTAMP
    }
    
    OfferImages {
        id BIGINT
        offer_id BIGINT
        image_id BIGINT
    }
    
    OfferPriceHistory {
        id BIGINT
        offer_id BIGINT
        price INT
        date INT
    }
    
    RentType {
        id INT
        name TEXT
    }
    
    OfferType {
        id INT
        name TEXT
    }
    
    PropertyType {
        id INT
        name TEXT
    }
    
    PurchaseType {
        id INT
        name TEXT
    }
    
    OfferStatus {
        id INT
        name TEXT
    }
    
    OfferRenovation {
        id INT
        name TEXT
    }
    
    MetroStation {
        id INT
        metro_line_id INT
        name TEXT
    }
    
    MetroLine {
        id INT
        name TEXT
        color INT
    }
    
    User {
        id BIGINT
        image_id BIGINT
        first_name TEXT
        last_name TEXT
        email TEXT
        password TEXT
        last_notification_id INT
        created_at TIMESTAMP
        updated_at TIMESTAMP
    }
    
    Image {
        id BIGINT
        uuid TEXT
        created_at TIMESTAMP
        updated_at TIMESTAMP
    }
    
    UserNotification {
        id BIGINT
        user_id BIGINT
        message TEXT
        redirect_uri TEXT
        created_at TIMESTAMP
    }
    
    UserReview {
        id BIGINT
        reviewer_id BIGINT
        seller_id BIGINT
        rating INT
        comment TEXT
        created_at TIMESTAMP
        updated_at TIMESTAMP
    }
    
    UserOfferFavourites {
        id BIGINT
        user_id BIGINT
        offer_id BIGINT
        created_at TIMESTAMP
    }
    
    Chat {
        id BIGINT
        offer_id BIGINT
        customer_id BIGINT
        created_at TIMESTAMP
        updated_at TIMESTAMP
    }
    
    ChatMessage {
        id BIGINT
        chat_id BIGINT
        user_id BIGINT
        message TEXT
        created_at TIMESTAMP
        updated_at TIMESTAMP
    }
    
    HousingComplex {
        id BIGINT
        class_id INT
        name TEXT
        developer TEXT
        created_at TIMESTAMP
        updated_at TIMESTAMP
    }
    
    HousingComplexClass {
        id INT
        name TEXT
    }
    
    HousingComplexReview {
        id BIGINT
        user_id BIGINT
        housing_complex_id BIGINT
        rating INT
        comment TEXT
        created_at TIMESTAMP
        updated_at TIMESTAMP
    }
    
    HousingComplexImages {
        id BIGINT
        housing_complex_id BIGINT
        image_id BIGINT
        created_at TIMESTAMP
        updated_at TIMESTAMP
    }
    
    Offer ||--o{ OfferImages : "has"
    Offer ||--o{ OfferPriceHistory : "has"
    Offer ||--o{ Chat : "has"
    Offer ||--o{ UserOfferFavourites : "has"
    OfferImages ||--|| Image : "has"

    RentType ||--o{ Offer : "has"
    PropertyType ||--o{ Offer : "has"
    PurchaseType ||--o{ Offer : "has"
    OfferType ||--o{ Offer : "has"
    OfferStatus ||--o{ Offer : "has"
    OfferRenovation ||--o{ Offer : "has"

    User ||--|| Image : "has"
    User ||--o{ UserOfferFavourites : "has"
    User ||--o{ Offer : "has"
    User ||--o{ Chat : "has"
    User ||--o{ ChatMessage : "has"
    User ||--o{ HousingComplexReview : "has"
    User ||--o{ UserNotification : "has"
    User ||--o{ UserReview : "has"
    
    
    Chat ||--o{ ChatMessage : "has"
    
    HousingComplex ||--o{ Offer : "has"
    HousingComplex ||--o{ HousingComplexImages : "has"
    HousingComplex ||--o{ HousingComplexReview : "has"
    HousingComplexClass ||--o{ HousingComplex : "has"
    HousingComplexImages ||--|| Image : "has"
    
    MetroLine ||--o{ MetroStation : "has"
    MetroStation ||--o{ Offer : "has"
```

# Таблицы БД
### 1. Таблица Offer
Содержит информацию о предложениях недвижимости.

{id} -> {seller_id, offer_type_id, metro_station_id, rent_type_id, purchase_type_id, property_type_id, offer_status_id, renovation_id, complex_id, price, description, floor, total_floors, rooms, address, flat, area, ceiling_height, created_at, updated_at}

### 2. Таблица OfferImages
Содержит информацию о изображениях предложений недвижимости.

{id} -> {offer_id, image_id}
{offer_id, image_id} -> {id}

### 3. Таблица OfferPriceHistory
Содержит информацию о истории цен предложений недвижимости.

{id} -> {offer_id, price, date}

### 4. Таблица RentType
Содержит информацию о типах аренды (посуточно, долгосрочно).

{id} -> {name}
{name} -> {id}

### 5. Таблица OfferType
Содержит информацию о типах предложений (продажа, аренда).

{id} -> {name}
{name} -> {id}

### 6. Таблица PropertyType
Содержит информацию о типах недвижимости (квартира, дом, апартаменты).

{id} -> {name}
{name} -> {id}

### 7. Таблица PurchaseType
Содержит информацию о типах недвижимости (новостройка, вторичное жилье).

{id} -> {name}
{name} -> {id}

### 8. Таблица OfferStatus
Содержит информацию о статусах предложений (активно, снято с продажи, черновик).

{id} -> {name}
{name} -> {id}

### 9. Таблица OfferRenovation
Содержит информацию о типах ремонта (без ремонта, косметический, евро).

{id} -> {name}
{name} -> {id}

### 10. Таблица MetroStation
Содержит информацию о станциях метро.

{id} -> {metro_line_id, name}
{metro_line_id, name} -> {id}

### 11. Таблица MetroLine
Содержит информацию о ветках метро.

{id} -> {name, color}
{name} -> {id, color}

### 12. Таблица User
Содержит информацию о пользователях.

{id} -> {image_id, first_name, last_name, email, password, last_notification_id, created_at, updated_at}
{email} -> {id, image_id, first_name, last_name, password, last_notification_id, created_at, updated_at}

### 13. Таблица Image
Содержит информацию о изображениях.

{id} -> {uuid, created_at, updated_at}
{uuid} -> {id, created_at, updated_at}

### 14. Таблица UserNotification
Содержит информацию о уведомлениях пользователей.

{id} -> {user_id, message, redirect_uri, created_at}

### 15. Таблица UserReview
Содержит информацию об отзывах о пользователях.

{id} -> {reviewer_id, seller_id, rating, comment, created_at, updated_at}

### 16. Таблица UserOfferFavourites
Содержит информацию о избранных предложениях пользователей.

{id} -> {user_id, offer_id, created_at}
{user_id, offer_id} -> {id, created_at}

### 17. Таблица Chat
Содержит информацию о чатах.

{id} -> {offer_id, customer_id, created_at, updated_at}
{offer_id, customer_id} -> {id, created_at, updated_at}

### 18. Таблица ChatMessage
Содержит информацию о сообщениях в чатах.

{id} -> {chat_id, user_id, message, created_at, updated_at}

### 19. Таблица HousingComplex
Содержит информацию о жилых комплексах.

{id} -> {class_id, name, developer, created_at, updated_at}

### 20. Таблица HousingComplexClass
Содержит информацию о классах жилых комплексов (комфорт, бизнес, элит, ...).

{id} -> {name}
{name} -> {id}

### 21. Таблица HousingComplexReview
Содержит информацию об отзывах о жилых комплексах.

{id} -> {user_id, housing_complex_id, rating, comment, created_at, updated_at}

### 22. Таблица HousingComplexImages
Содержит информацию о изображениях жилых комплексов.

{id} -> {housing_complex_id, image_id, created_at, updated_at}
{housing_complex_id, image_id} -> {id, image_id, created_at, updated_at}

## БД находится в 3-ей нормальной форме Бойса-Кодда
Для доказательства того, что БД находится в 3-ей нормальной форме Бойса-Кодда, рассмотрим функциональные зависимости в таблицах БД. Они расписаны в предыдущем блоке для каждой таблицы.

### 1 нормальная форма
Все отношения удовлетворяют 1НФ, так как каждый кортеж содержит только одно значение для каждого атрибута.

### 2 нормальная форма
Все отношения удовлетворяют 2НФ, так как:
- удовлетворяют 1НФ
- все атрибуты, не входящие в первичный ключ, полностью функционально зависят от первичного ключа

### 3 нормальная форма
Все отношения удовлетворяют 3НФ, так как:
- удовлетворяют 2НФ
- нет транзитивных зависимостей

### 3 нормальная форма Бойса-Кодда
Все отношения удовлетворяют 3НФ Бойса-Кодда, так как:
- удовлетворяют 3НФ
- все левые части функциональных зависимостей являются ключами
