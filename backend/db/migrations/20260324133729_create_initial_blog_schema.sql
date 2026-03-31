-- +goose Up
CREATE TABLE articles (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	title VARCHAR(255) NOT NULL,
	body TEXT NOT NULL,
	is_published BOOLEAN NOT NULL DEFAULT FALSE,
	publish_start_at DATETIME NULL,
	publish_end_at DATETIME NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);

CREATE TABLE categories (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	name VARCHAR(64) NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE KEY uq_categories_name (name)
);

CREATE TABLE article_categories (
	article_id INT UNSIGNED NOT NULL,
	category_id INT UNSIGNED NOT NULL,
	PRIMARY KEY (article_id, category_id),
	KEY idx_article_categories_category_id (category_id),
	CONSTRAINT fk_article_categories_article_id
	FOREIGN KEY (article_id) REFERENCES articles (id)
	ON DELETE CASCADE,
	CONSTRAINT fk_article_categories_category_id
	FOREIGN KEY (category_id) REFERENCES categories (id)
	ON DELETE CASCADE
);

CREATE TABLE comments (
	id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	article_id INT UNSIGNED NOT NULL,
	author_name VARCHAR(64) NOT NULL,
	body TEXT NOT NULL,
	is_visible BOOLEAN NOT NULL DEFAULT TRUE,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	KEY idx_comments_article_id_created_at (article_id, created_at),
	CONSTRAINT fk_comments_article_id
	FOREIGN KEY (article_id) REFERENCES articles (id)
	ON DELETE CASCADE
);

-- +goose Down
DROP TABLE comments;
DROP TABLE article_categories;
DROP TABLE categories;
DROP TABLE articles;
