CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	username TEXT NOT NULL,
	is_active BOOLEAN NOT NULL
);

CREATE INDEX idx_users_username ON users(username);
