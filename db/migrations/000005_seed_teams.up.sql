INSERT INTO teams (name)
SELECT 'team_' || g::text
FROM generate_series(1, 20) g;

INSERT INTO team_members (team_id, user_id)
SELECT
    t.id,
    u.id
FROM teams t
JOIN LATERAL (
    SELECT id
    FROM users
    ORDER BY random()
    LIMIT (1 + floor(random() * 3)::int)
) u ON true;
