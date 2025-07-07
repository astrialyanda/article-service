-- +goose Up
CREATE TABLE IF NOT EXISTS articles (
    id TEXT PRIMARY KEY,
    author_id TEXT NOT NULL,
    title VARCHAR(200) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS authors (
    id TEXT PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE INDEX idx_articles_created_at ON articles(created_at DESC);
CREATE INDEX idx_articles_author_id ON articles(author_id);
CREATE INDEX idx_articles_title_gin ON articles USING gin(to_tsvector('english', title));
CREATE INDEX idx_articles_body_gin ON articles USING gin(to_tsvector('english', body));

