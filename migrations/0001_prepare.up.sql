CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE SCHEMA IF NOT EXISTS finstar;

CREATE TABLE IF NOT EXISTS finstar.users(
	id uuid DEFAULT uuid_generate_v4(),
 	name varchar NOT NUll,
 	PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS finstar.accounts(
	id uuid DEFAULT uuid_generate_v4(),
 	user_id uuid,
  	balance numeric default 0 CHECK (balance >= 0),
  	PRIMARY KEY (id),
  	FOREIGN KEY (user_id) REFERENCES finstar.users (id)
);

--seed data
WITH u AS (
	INSERT INTO finstar.users(name)
	VALUES('John Smith')
	RETURNING id
)
INSERT INTO finstar.accounts(user_id)
SELECT id FROM u;

WITH u AS (
	INSERT INTO finstar.users(name)
	VALUES('Jane Smith')
	RETURNING id
)
INSERT INTO finstar.accounts(user_id)
SELECT id FROM u;

WITH u AS (
	INSERT INTO finstar.users(name)
	VALUES('Alex Smith')
	RETURNING id
)
INSERT INTO finstar.accounts(user_id)
SELECT id FROM u;