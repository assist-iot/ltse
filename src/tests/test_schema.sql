CREATE SCHEMA example2;
SET search_path = example2;

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  email VARCHAR(50) NOT NULL UNIQUE,
  password VARCHAR(50) NOT NULL
);

INSERT INTO users (name, email, password) VALUES
  ('Juan', 'juan@example.com', 'password1'),
  ('María', 'maria@example.com', 'password2'),
  ('Pedro', 'pedro@example.com', 'password3');

CREATE TABLE products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description TEXT,
  price NUMERIC(8,2) NOT NULL
);

INSERT INTO products (name, description, price) VALUES
  ('Product 1', 'Product description 1', 10.99),
  ('Product 2', 'Product description 2', 15.99),
  ('Product 3', 'Product description 3', 20.99);

RESET search_path;


