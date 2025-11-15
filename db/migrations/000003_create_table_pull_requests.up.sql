CREATE TABLE pull_requests (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	author_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	status TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	merged_at TIMESTAMPTZ
);

CREATE TABLE pull_request_reviewers (
	id SERIAL PRIMARY KEY,
	pull_request_id TEXT NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
	reviewer_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_pull_reqeusts ON pull_requests(author_id);
