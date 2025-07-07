package service

import (
    "context"
    "fmt"

    "article-service/internal/model"
    "article-service/internal/repository"
    "github.com/google/uuid"
)

type ArticleService interface {
    CreateArticle(ctx context.Context, req *model.CreateArticleRequest) (*model.Article, error)
    GetArticles(ctx context.Context, req *model.GetArticlesRequest) (*model.GetArticlesResponse, error)
}

type articleService struct {
    repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
    return &articleService{
        repo: repo,
    }
}

func (s *articleService) CreateArticle(ctx context.Context, req *model.CreateArticleRequest) (*model.Article, error) {
    article := &model.Article{
        ID:       uuid.NewString(),
        AuthorID: req.AuthorID,
        Title:    req.Title,
        Body:     req.Body,
    }

    if err := s.repo.Create(ctx, article); err != nil {
        return nil, fmt.Errorf("failed to create article: %w", err)
    }

    return article, nil
}

func (s *articleService) GetArticles(ctx context.Context, req *model.GetArticlesRequest) (*model.GetArticlesResponse, error) {
    // Validate and set defaults
    if req.Page < 1 {
        req.Page = 1
    }
    if req.Limit < 1 || req.Limit > 100 {
        req.Limit = 10
    }

    articles, total, err := s.repo.GetList(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to get articles: %w", err)
    }

    return &model.GetArticlesResponse{
        Articles: articles,
        Total:    total,
        Page:     req.Page,
        Limit:    req.Limit,
    }, nil
}
