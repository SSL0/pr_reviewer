CREATE TABLE users (
	id TEXT PRIMARY KEY,
	username TEXT NOT NULL,
	team_name TEXT NOT NULL REFERENCES teams(name) ON DELETE CASCADE,
	is_active BOOLEAN NOT NULL
);

CREATE INDEX idx_users_username ON users(username);
