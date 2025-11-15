INSERT INTO pull_requests (name, author_id, status)
SELECT
    'PR_' || u.id || '_' || g::text,
    u.id,
    CASE WHEN random() < 0.5 THEN 'open' ELSE 'closed' END
FROM users u
JOIN generate_series(1, 3) g ON random() < 0.7;
