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

### 2. Таблица OfferImages
Содержит информацию о изображениях предложений недвижимости.

### 3. Таблица OfferPriceHistory
Содержит информацию о истории цен предложений недвижимости.

### 4. Таблица RentType
Содержит информацию о типах аренды (посуточно, долгосрочно).

### 5. Таблица OfferType
Содержит информацию о типах предложений (продажа, аренда).

### 6. Таблица PropertyType
Содержит информацию о типах недвижимости (квартира, дом, апартаменты).

### 7. Таблица PurchaseType
Содержит информацию о типах недвижимости (новостройка, вторичное жилье).

### 8. Таблица OfferStatus
Содержит информацию о статусах предложений (активно, снято с продажи, черновик).

### 9. Таблица OfferRenovation
Содержит информацию о типах ремонта (без ремонта, косметический, евро).

### 10. Таблица MetroStation
Содержит информацию о станциях метро.

### 11. Таблица MetroLine
Содержит информацию о ветках метро.

### 12. Таблица User
Содержит информацию о пользователях.

### 13. Таблица Image
Содержит информацию о изображениях.

### 14. Таблица UserNotification
Содержит информацию о уведомлениях пользователей.

### 15. Таблица UserReview
Содержит информацию об отзывах о пользователях.

### 16. Таблица UserOfferFavourites
Содержит информацию о избранных предложениях пользователей.

### 17. Таблица Chat
Содержит информацию о чатах.

### 18. Таблица ChatMessage
Содержит информацию о сообщениях в чатах.

### 19. Таблица HousingComplex
Содержит информацию о жилых комплексах.

### 20. Таблица HousingComplexClass
Содержит информацию о классах жилых комплексов (комфорт, бизнес, элит, ...).

### 21. Таблица HousingComplexReview
Содержит информацию об отзывах о жилых комплексах.

### 22. Таблица HousingComplexImages
Содержит информацию о изображениях жилых комплексов.