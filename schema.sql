CREATE TABLE categories(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL,
    category_id INTEGER REFERENCES categories(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);