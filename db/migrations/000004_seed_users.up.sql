INSERT INTO users (username, is_active)
SELECT
    'user_' || g::text,
    CASE WHEN random() < 0.8 THEN true ELSE false END
FROM generate_series(1, 200) g;
