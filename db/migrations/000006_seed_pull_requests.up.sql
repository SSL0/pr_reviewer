
INSERT INTO pull_requests (id, name, author_id, status) VALUES
('pr-1','Add login feature','u1','open'),
('pr-2','Fix payment bug','u5','merged'),
('pr-3','Update CI pipeline','u6','open'),
('pr-4','Frontend redesign','u3','closed'),
('pr-5','Mobile app bugfix','u8','open'),
('pr-6','Data migration','u9','open');

INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id) VALUES
('pr-1','u2'), ('pr-1','u3'),
('pr-2','u6'), ('pr-2','u7'),
('pr-3','u1'), ('pr-3','u7'),
('pr-4','u5'), ('pr-4','u8'),
('pr-5','u9'), ('pr-5','u10'),
('pr-6','u11'), ('pr-6','u12');
