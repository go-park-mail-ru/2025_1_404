SET SEARCH_PATH = csat;

INSERT INTO CSAT (event) VALUES
    ('create_offer');

INSERT INTO Question (text, csat_id) VALUES
    ('Насколько вы удовлетворены процессом создания объявления?', (SELECT id FROM CSAT WHERE event = 'create_offer'));

INSERT INTO Answer (rating, question_id) VALUES
    (1, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (1, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (3, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (4, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (5, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?')),
    (5, (SELECT id FROM Question WHERE text = 'Насколько вы удовлетворены процессом создания объявления?'));