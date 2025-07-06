package repository

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "time"

    "article-service/internal/model"

    _ "github.com/lib/pq"
)

type ArticleRepository interface {
    Create(ctx context.Context, article *model.Article) error
    GetList(ctx context.Context, req *model.GetArticlesRequest) ([]model.Article, int, error)
}

type articleRepository struct {
    db *sql.DB
}

func NewDB(databaseURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    // Set connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    return db, nil
}

func NewArticleRepository(db *sql.DB) ArticleRepository {
    return &articleRepository{
        db: db,
    }
}

func (r *articleRepository) Create(ctx context.Context, article *model.Article) error {
    query := `
        INSERT INTO articles (author_id, title, body, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at
    `
    now := time.Now()
    err := r.db.QueryRowContext(ctx, query, article.AuthorID, article.Title, article.Body, now).
        Scan(&article.ID, &article.CreatedAt)
    if err != nil {
        return fmt.Errorf("failed to create article: %w", err)
    }
    return nil
}

func (r *articleRepository) GetList(ctx context.Context, req *model.GetArticlesRequest) ([]model.Article, int, error) {
    var conditions []string
    var args []interface{}
    argIndex := 1

    if req.Query != "" {
        conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR body ILIKE $%d)", argIndex, argIndex))
        args = append(args, "%"+req.Query+"%")
        argIndex++
    }

    if req.AuthorID != "" {
        conditions = append(conditions, fmt.Sprintf("author_id = $%d", argIndex))
        args = append(args, req.AuthorID)
        argIndex++
    }

    whereClause := ""
    if len(conditions) > 0 {
        whereClause = "WHERE " + strings.Join(conditions, " AND ")
    }

    countQuery := fmt.Sprintf("SELECT COUNT(*) FROM articles %s", whereClause)
    var total int
    err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get total count: %w", err)
    }

    offset := (req.Page - 1) * req.Limit
    query := fmt.Sprintf(`
        SELECT id, author_id, title, body, created_at
        FROM articles
        %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, argIndex, argIndex+1)

    args = append(args, req.Limit, offset)

    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get articles: %w", err)
    }
    defer rows.Close()

    var articles []model.Article
    for rows.Next() {
        var article model.Article
        err := rows.Scan(&article.ID, &article.AuthorID, &article.Title, &article.Body, &article.CreatedAt)
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan article: %w", err)
        }
        articles = append(articles, article)
    }

    if err := rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("rows iteration error: %w", err)
    }

    return articles, total, nil
}