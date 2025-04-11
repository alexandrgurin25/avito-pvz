
INSERT INTO cities (name) VALUES 
    ('Москва'),
    ('Санкт-Петербург'),
    ('Казань')
ON CONFLICT (name) DO NOTHING;