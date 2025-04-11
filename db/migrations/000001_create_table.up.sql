
CREATE TABLE IF NOT EXISTS cities(
   id SERIAL PRIMARY KEY,
   name TEXT
);

CREATE TABLE IF NOT EXISTS pvz(
   id         UUID     PRIMARY KEY,
   city       INT      NOT NULL,
   created_at TIMESTAMP DEFAULT NOW(),
   FOREIGN KEY (city) REFERENCES cities(id)
);
