CREATE TABLE IF NOT EXISTS cities (
   id SMALLSERIAL PRIMARY KEY,
   name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS pvz (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   city_id SMALLINT NOT NULL REFERENCES cities(id),
   created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'Europe/Moscow') 
);

CREATE TABLE IF NOT EXISTS categories (
    id SMALLSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE 
);

CREATE TABLE IF NOT EXISTS receivings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pvz_id UUID NOT NULL REFERENCES pvz(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('in_progress', 'close')) DEFAULT 'in_progress',
    start_time TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'Europe/Moscow'),
    end_time TIMESTAMP
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    receiving_id UUID NOT NULL REFERENCES receivings(id),
    category_id SMALLINT NOT NULL REFERENCES categories(id),
    added_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'Europe/Moscow')
);


CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role VARCHAR(10) NOT NULL CHECK (role IN ('moderator', 'employee')),
    pvz_id UUID REFERENCES pvz(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Эти индексы НУЖНЫ, т.к. они не связаны с FK:
CREATE INDEX IF NOT EXISTS idx_receivings_status ON receivings(status); -- Ускоряет фильтрацию по статусу
CREATE INDEX IF NOT EXISTS idx_receivings_time_range ON receivings(start_time, end_time); -- Для временных запросов
CREATE INDEX IF NOT EXISTS idx_products_added_at ON products(added_at); -- Для аналитики по дате
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role); -- Ускоряет фильтрацию по роли
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at); -- Для отчетов по регистрации