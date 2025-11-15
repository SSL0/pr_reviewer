CREATE TABLE teams (
	name TEXT PRIMARY KEY
);

CREATE INDEX idx_teams_name ON teams(name);
