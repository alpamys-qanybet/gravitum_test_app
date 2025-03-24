CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	surname VARCHAR(100) NULL,
	inserted_at timestamptz DEFAULT NOW(),
	updated_at timestamptz NULL
);