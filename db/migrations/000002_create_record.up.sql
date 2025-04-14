INSERT INTO cities (id, name) VALUES 
    (1, 'Москва'),
    (2, 'Санкт-Петербург'),
    (3, 'Казань')
ON CONFLICT (id) DO NOTHING;

INSERT INTO categories (id, name) VALUES 
    (1, 'электроника'),
    (2, 'одежда'),
    (3, 'обувь')
ON CONFLICT (id) DO NOTHING;
