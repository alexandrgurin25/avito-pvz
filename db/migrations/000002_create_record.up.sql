INSERT INTO cities (name) VALUES 
    ('Москва'),
    ('Санкт-Петербург'),
    ('Казань')
ON CONFLICT (name) DO NOTHING; 

INSERT INTO categories (name) VALUES 
    ('электроника'),
    ('одежда'),
    ('обувь')
ON CONFLICT (name) DO NOTHING;
