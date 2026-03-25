SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE comments;
TRUNCATE TABLE article_categories;
TRUNCATE TABLE articles;
TRUNCATE TABLE categories;

SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO articles (id, title, body, is_published, publish_start_at)
WITH RECURSIVE seq AS (
  SELECT 1 AS n
  UNION ALL
  SELECT n + 1
  FROM seq
  WHERE n < 100
)
SELECT
  n,
  CONCAT('かさまし記事', n),
  '記事の件数を増やすためのかさまし記事です。内容は適当です。',
  TRUE,
  NOW() - INTERVAL 1 DAY
FROM seq;

INSERT INTO categories (id, name)
VALUES
  (1, 'Go'),
  (2, 'Docker'),
  (3, 'Cloudflare'),
  (4, 'AWS'),
  (5, 'Raspberry Pi');

INSERT INTO article_categories (article_id, category_id)
VALUES
  (1, 1),
  (1, 2),
  (2, 3),
  (2, 4),
  (3, 1),
  (3, 5);
