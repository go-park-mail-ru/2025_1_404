SET SEARCH_PATH = kvartirum;

INSERT INTO offertype (name) VALUES
	('Продажа'),
	('Аренда');
	
INSERT INTO renttype (name) VALUES
	('Посуточно'),
	('Долгосрок');

INSERT INTO PurchaseType (name) VALUES
    ('Новостройка'),
    ('Вторичка');

INSERT INTO propertytype (name) VALUES
	('Апартаменты'),
	('Дом'),
	('Квартира');

INSERT INTO OfferStatus (name) VALUES
	('Активный'),
	('Черновик'),
	('Завершенный');

INSERT INTO offerrenovation (name) VALUES
	('Современный ремонт'),
	('Косметический ремонт'),
	('Черновая отделка'),
	('Нужен полный ремонт'),
	('Нужен частичный ремонт'),
	('Улучшенная черновая');

INSERT INTO HousingComplexClass (name) VALUES
	('Комфорт'),
	('Бизнес'),
	('Элит'),
	('Эконом');

INSERT INTO metroline (name, color) VALUES
    ('Арбатско-Покровская', '0033A0'),
    ('Большая кольцевая линия', '82C0C0'),
    ('Бутовская', 'A1A2A3'),
    ('Замоскворецкая', '007D3C'),
    ('Калининская', 'FFD702'),
    ('Калужско-Рижская', 'FFA300'),
    ('Каховская', 'A1A2A3'),
    ('Кольцевая', '894E35'),
    ('Люблинско-Дмитровская', '9EC862'),
    ('МЦД-1', 'F6A600'),
    ('МЦД-2', '0078BE'),
    ('МЦД-3', 'DEA62C'),
    ('МЦД-4', 'AD1D33'),
    ('МЦД-5', 'ADB3B8'),
    ('МЦК', 'FF8642'),
    ('Некрасовская', 'DE64A1'),
    ('Рублёво-Архангельская', '78C596'),
    ('Серпуховско-Тимирязевская', 'A1A2A3'),
    ('Сокольническая', 'EF161E'),
    ('Таганско-Краснопресненская', '97005E'),
    ('Троицкая', '009A49'),
    ('Солнцевская', 'FFD702'),
    ('Филевская', '0078BE');

INSERT INTO metrostation (metro_line_id, name) VALUES
    ((SELECT id FROM metroline WHERE name = 'Калининская'), 'Новокосино'),
    ((SELECT id FROM metroline WHERE name = 'Калининская'), 'Новогиреево'),
    ((SELECT id FROM metroline WHERE name = 'Калининская'), 'Перово'),
    ((SELECT id FROM metroline WHERE name = 'Калининская'), 'Шоссе Энтузиастов'),
    ((SELECT id FROM metroline WHERE name = 'Калининская'), 'Авиамоторная'),
    ((SELECT id FROM metroline WHERE name = 'Калининская'), 'Площадь Ильича'),
    ((SELECT id FROM metroline WHERE name = 'Калининская'), 'Марксистская'),
    ((SELECT id FROM metroline WHERE name = 'Калининская'), 'Третьяковская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Ховрино'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Беломорская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Речной вокзал'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Водный стадион'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Войковская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Сокол'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Аэропорт'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Динамо'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Белорусская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Маяковская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Тверская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Театральная'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Новокузнецкая'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Павелецкая'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Автозаводская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Технопарк'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Коломенская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Каширская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Кантемировская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Царицыно'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Орехово'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Домодедовская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Красногвардейская'),
    ((SELECT id FROM metroline WHERE name = 'Замоскворецкая'), 'Алма-Атинская'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Медведково'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Бабушкинская'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Свиблово'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Ботанический сад'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'ВДНХ'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Алексеевская'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Рижская'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Проспект Мира'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Сухаревская'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Тургеневская'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Китай-город'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Третьяковская'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Октябрьская'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Шаболовская'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Ленинский проспект'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Академическая'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Профсоюзная'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Новые Черемушки'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Калужская'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Беляево'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Коньково'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Теплый Стан'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Ясенево'),
    ((SELECT id FROM metroline WHERE name = 'Калужско-Рижская'), 'Новоясеневская'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Бульвар Рокоссовского'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Черкизовская'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Преображенская площадь'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Сокольники'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Красносельская'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Комсомольская'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Красные ворота'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Чистые пруды'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Лубянка'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Охотный ряд'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Библиотека им.Ленина'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Кропоткинская'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Парк культуры'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Фрунзенская'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Спортивная'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Воробьевы горы'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Университет'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Проспект Вернадского'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Юго-Западная'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Тропарево'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Румянцево'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Саларьево'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Филатов луг'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Прокшино'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Ольховая'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Новомосковская (Коммунарка)'),
    ((SELECT id FROM metroline WHERE name = 'Сокольническая'), 'Потапово'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Щелковская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Первомайская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Измайловская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Партизанская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Семеновская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Электрозаводская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Бауманская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Курская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Площадь Революции'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Арбатская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Смоленская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Киевская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Парк Победы'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Славянский бульвар'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Кунцевская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Молодежная'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Крылатское'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Строгино'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Мякинино'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Волоколамская'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Митино'),
    ((SELECT id FROM metroline WHERE name = 'Арбатско-Покровская'), 'Пятницкое шоссе'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Кунцевская'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Пионерская'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Филевский парк'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Багратионовская'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Фили'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Кутузовская'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Студенческая'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Киевская'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Смоленская'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Арбатская'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Александровский сад'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Деловой центр (Выставочная)'),
    ((SELECT id FROM metroline WHERE name = 'Филевская'), 'Москва-Сити'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Алтуфьево'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Бибирево'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Отрадное'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Владыкино'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Петровско-Разумовская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Тимирязевская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Дмитровская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Савёловская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Менделеевская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Цветной бульвар'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Чеховская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Боровицкая'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Полянка'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Серпуховская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Тульская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Нагатинская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Нагорная'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Нахимовский проспект'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Севастопольская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Чертановская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Южная'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Пражская'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Улица Академика Янгеля'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Аннино'),
    ((SELECT id FROM metroline WHERE name = 'Серпуховско-Тимирязевская'), 'Бульвар Дмитрия Донского'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Планерная'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Сходненская'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Тушинская'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Спартак'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Щукинская'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Октябрьское поле'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Полежаевская'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Беговая'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Улица 1905 года'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Баррикадная'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Пушкинская'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Кузнецкий мост'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Китай-город'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Таганская'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Пролетарская'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Волгоградский проспект'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Текстильщики'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Кузьминки'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Рязанский проспект'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Выхино'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Лермонтовский проспект'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Жулебино'),
    ((SELECT id FROM metroline WHERE name = 'Таганско-Краснопресненская'), 'Котельники'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Новослободская'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Проспект Мира'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Комсомольская'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Курская'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Таганская'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Павелецкая'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Добрынинская'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Октябрьская'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Парк культуры'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Киевская'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Краснопресненская'),
    ((SELECT id FROM metroline WHERE name = 'Кольцевая'), 'Белорусская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Физтех'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Лианозово'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Яхромская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Селигерская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Верхние Лихоборы'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Окружная'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Петровско-Разумовская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Фонвизинская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Бутырская '),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Марьина Роща'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Достоевская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Трубная'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Сретенский бульвар'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Чкаловская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Римская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Крестьянская застава'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Дубровка'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Кожуховская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Печатники'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Волжская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Люблино'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Братиславская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Марьино'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Борисово'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Шипиловская'),
    ((SELECT id FROM metroline WHERE name = 'Люблинско-Дмитровская'), 'Зябликово'),
    ((SELECT id FROM metroline WHERE name = 'Каховская'), 'Каширская'),
    ((SELECT id FROM metroline WHERE name = 'Каховская'), 'Варшавская'),
    ((SELECT id FROM metroline WHERE name = 'Каховская'), 'Каховская'),
    ((SELECT id FROM metroline WHERE name = 'Бутовская'), 'Бунинская аллея'),
    ((SELECT id FROM metroline WHERE name = 'Бутовская'), 'Улица Горчакова'),
    ((SELECT id FROM metroline WHERE name = 'Бутовская'), 'Бульвар Адмирала Ушакова'),
    ((SELECT id FROM metroline WHERE name = 'Бутовская'), 'Улица Скобелевская'),
    ((SELECT id FROM metroline WHERE name = 'Бутовская'), 'Улица Старокачаловская'),
    ((SELECT id FROM metroline WHERE name = 'Бутовская'), 'Лесопарковая'),
    ((SELECT id FROM metroline WHERE name = 'Бутовская'), 'Битцевский Парк'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Деловой центр'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Пыхтино'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Парк Победы'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Минская'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Ломоносовский проспект'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Раменки'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Мичуринский проспект'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Озёрная'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Говорово'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Солнцево'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Боровское шоссе'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Новопеределкино'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Рассказовка'),
    ((SELECT id FROM metroline WHERE name = 'Солнцевская'), 'Аэропорт Внуково'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Окружная'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Владыкино'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Ботанический сад'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Ростокино'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Белокаменная'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Бульвар Рокоссовского'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Локомотив'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Измайлово'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Соколиная Гора'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Шоссе Энтузиастов'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Андроновка'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Нижегородская'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Новохохловская'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Угрешская'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Дубровка'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Автозаводская'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'ЗИЛ'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Верхние Котлы'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Крымская'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Площадь Гагарина'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Лужники'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Кутузовская'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Деловой центр'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Шелепиха'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Хорошево'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Зорге'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Панфиловская'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Стрешнево'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Балтийская'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Коптево'),
    ((SELECT id FROM metroline WHERE name = 'МЦК'), 'Лихоборы'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Мичуринский проспект'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Проспект Вернадского'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Новаторская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Воронцовская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Зюзино'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Каховская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Варшавская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Каширская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Кленовый бульвар'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Нагатинский Затон'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Печатники'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Текстильщики'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Нижегородская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Авиамоторная'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Лефортово'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Электрозаводская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Сокольники'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Рижская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Марьина Роща'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Савёловская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Петровский парк'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'ЦСКА'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Хорошевская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Шелепиха'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Деловой центр'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Народное Ополчение'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Мнёвники'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Терехово'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Кунцевская'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Давыдково'),
    ((SELECT id FROM metroline WHERE name = 'Большая кольцевая линия'), 'Аминьевская'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Косино'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Улица Дмитриевского'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Лухмановская'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Некрасовка'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Юго-Восточная'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Окская'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Стахановская'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Нижегородская'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Лефортово'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Электрозаводская'),
    ((SELECT id FROM metroline WHERE name = 'Некрасовская'), 'Авиамоторная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Лобня'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Шереметьевская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Хлебниково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Водники'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Долгопрудная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Новодачная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Марк'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Лианозово'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Бескудниково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Дегунино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Окружная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Тимирязевская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Савёловская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Белорусская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Беговая'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Тестовская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Фили'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Славянский бульвар'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Кунцевская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Рабочий Посёлок'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Сетунь'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Немчиновка'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Сколково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Баковка'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-1'), 'Одинцово'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Нахабино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Аникеевка'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Опалиха'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Красногорская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Павшино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Пенягино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Волоколамская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Трикотажная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Тушинская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Покровское-Стрешнево'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Стрешнево'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Красный Балтиец'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Гражданская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Дмитровская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Рижская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Площадь трёх вокзалов'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Курская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Москва Товарная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Калитники'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Новохохловская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Текстильщики'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Люблино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Депо'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Перерва'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Курьяново'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Москворечье'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Царицыно'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Покровское'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Красный строитель'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Битца'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Бутово'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Щербинка'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Остафьево'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Силикатная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Подольск'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-2'), 'Марьина Роща'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Крюково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Малино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Фирсановская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Сходня'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Подрезково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Новоподрезково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Молжаниново'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Химки'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Левобережная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Ховрино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Грачёвская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Моссельмаш'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Лихоборы'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Петровско-Разумовская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Останкино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Рижская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Митьково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Электрозаводская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Сортировочная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Авиамоторная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Андроновка'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Перово'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Плющево'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Вешняки'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Выхино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Косино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Ухтомская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Люберцы I'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Панки'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Томилино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Красково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Малаховка'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Удельная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Быково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Ильинская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Отдых'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Кратово'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Есенинская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Фабричная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Раменское'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-3'), 'Ипподром'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Железнодорожная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Ольгино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Кучино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Салтыковская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Никольское'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Реутово'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Новогиреево'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Кусково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Чухлинка'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Нижегородская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Серп и Молот'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Курская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Площадь трёх вокзалов'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Рижская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Марьина Роща'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Савёловская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Белорусская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Ермакова Роща'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Тестовская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Кутузовская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Поклонная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Минская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Матвеевское'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Аминьевская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Очаково I'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Мещерская'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Солнечная'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Новопеределкино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Переделкино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Мичуринец'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Внуково'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Лесной Городок'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Толстопальцево'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Кокошкино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Санино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Крёкшино'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Победа'),
    ((SELECT id FROM metroline WHERE name = 'МЦД-4'), 'Апрелевка'),
    ((SELECT id FROM metroline WHERE name = 'Троицкая'), 'Новаторская'),
    ((SELECT id FROM metroline WHERE name = 'Троицкая'), 'Университет дружбы народов'),
    ((SELECT id FROM metroline WHERE name = 'Троицкая'), 'Генерала Тюленева'),
    ((SELECT id FROM metroline WHERE name = 'Троицкая'), 'Тютчевская'),
    ((SELECT id FROM metroline WHERE name = 'Троицкая'), 'Корниловская'),
    ((SELECT id FROM metroline WHERE name = 'Троицкая'), 'Коммунарка'),
    ((SELECT id FROM metroline WHERE name = 'Троицкая'), 'Новомосковская');