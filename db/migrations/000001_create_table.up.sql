
CREATE TABLE IF NOT EXISTS cities (
   id SMALLSERIAL PRIMARY KEY,
   name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS pvz (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   city_id SMALLINT NOT NULL REFERENCES cities(id),
   created_at TIMESTAMP NOT NULL 
);

CREATE TABLE IF NOT EXISTS categories (
    id SMALLSERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS receivings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pvz_id UUID NOT NULL REFERENCES pvz(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('in_progress', 'closed')) DEFAULT 'in_progress',
    start_time TIMESTAMP NOT NULL DEFAULT NOW(),
    end_time TIMESTAMP
);


CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    receiving_id UUID NOT NULL REFERENCES receivings(id),
    category_id SMALLINT NOT NULL REFERENCES categories(id),
    added_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- CREATE TABLE IF NOT EXISTS users (
--     id SERIAL PRIMARY KEY,
--     email TEXT NOT NULL UNIQUE,
--     password_hash TEXT NOT NULL,
--     role VARCHAR(8) NOT NULL CHECK (role IN ('client', 'moderator', 'employee')),
--     pvz_id UUID REFERENCES pvz(id),
--     created_at TIMESTAMP NOT NULL DEFAULT NOW()
-- );