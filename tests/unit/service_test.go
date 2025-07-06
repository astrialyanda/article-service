package unit

import (
    "context"
    "errors"
    "testing"

    "article-service/internal/model"
    "article-service/internal/service"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockArticleRepository struct {
    mock.Mock
}

func (m *MockArticleRepository) Create(ctx context.Context, article *model.Article) error {
    args := m.Called(ctx, article)
    return args.Error(0)
}

func (m *MockArticleRepository) GetList(ctx context.Context, req *model.GetArticlesRequest) ([]model.Article, int, error) {
    args := m.Called(ctx, req)
    return args.Get(0).([]model.Article), args.Int(1), args.Error(2)
}

func TestArticleService_CreateArticle(t *testing.T) {
    mockRepo := new(MockArticleRepository)
    articleService := service.NewArticleService(mockRepo)

    req := &model.CreateArticleRequest{
        Title:  "Test Article",
        Body:   "Test Body",
        AuthorID: "Test Author",
    }

    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Article")).Return(nil)

    article, err := articleService.CreateArticle(context.Background(), req)

    assert.NoError(t, err)
    assert.Equal(t, req.Title, article.Title)
    assert.Equal(t, req.Body, article.Body)
    assert.Equal(t, req.AuthorID, article.AuthorID)
    mockRepo.AssertExpectations(t)
}

func TestArticleService_CreateArticle_Error(t *testing.T) {
    mockRepo := new(MockArticleRepository)
    articleService := service.NewArticleService(mockRepo)

    req := &model.CreateArticleRequest{
        Title:  "Test Article",
        Body:   "Test Body",
        AuthorID: "Test Author",
    }

    mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Article")).
        Return(errors.New("database error"))

    article, err := articleService.CreateArticle(context.Background(), req)

    assert.Error(t, err)
    assert.Nil(t, article)
    mockRepo.AssertExpectations(t)
}

func TestArticleService_GetArticles(t *testing.T) {
    mockRepo := new(MockArticleRepository)
    articleService := service.NewArticleService(mockRepo)

    req := &model.GetArticlesRequest{
        Page:  1,
        Limit: 10,
    }

    expectedArticles := []model.Article{
        {ID: "1", Title: "Article 1", Body: "Body 1", AuthorID: "author-1"},
        {ID: "2", Title: "Article 2", Body: "Body 2", AuthorID: "author-2"},
    }

    mockRepo.On("GetList", mock.Anything, req).Return(expectedArticles, 2, nil)

    response, err := articleService.GetArticles(context.Background(), req)

    assert.NoError(t, err)
    assert.Equal(t, expectedArticles, response.Articles)
    assert.Equal(t, 2, response.Total)
    assert.Equal(t, 1, response.Page)
    assert.Equal(t, 10, response.Limit)
    mockRepo.AssertExpectations(t)
}