CREATE TABLE pull_requests (
	id SERIAL PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	status TEXT NOT NULL
);

CREATE TABLE pull_request_reviewers (
	id SERIAL PRIMARY KEY NOT NULL,
	pull_request_id INTEGER NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
	reviewer_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_pull_reqeusts ON pull_requests(author_id);
