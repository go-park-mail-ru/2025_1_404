SET SEARCH_PATH = csat;

INSERT INTO CSAT (event) VALUES
    ('create_offer'),
    ('edit_offer'),
    ('view_offers');

INSERT INTO Question (text, csat_id) VALUES
    ('Насколько вы удовлетворены процессом создания объявления?', (SELECT id FROM CSAT WHERE event = 'create_offer')),
    ('Насколько вы удовлетворены процессом обновления объявления?', (SELECT id FROM CSAT WHERE event = 'edit_offer')),
    ('Насколько вы удовлетворены удобством Квартирума?', (SELECT id FROM CSAT WHERE event = 'view_offers'));



INSERT INTO Answer (rating, question_id) VALUES
    (1, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (1, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (3, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (4, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (5, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (5, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (5, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом обновления объявления?')),
    (4, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом обновления объявления?')),
    (4, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены удобством Квартирума?')),
    (5, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены удобством Квартирума?')),
    (5, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены удобством Квартирума?'));

    