package integration

import (
    "bytes"
	"context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "article-service/internal/handler"
    "article-service/internal/model"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockArticleService struct {
    mock.Mock
}

func (m *MockArticleService) CreateArticle(ctx context.Context, req *model.CreateArticleRequest) (*model.Article, error) {
    args := m.Called(ctx, req)
    return args.Get(0).(*model.Article), args.Error(1)
}

func (m *MockArticleService) GetArticles(ctx context.Context, req *model.GetArticlesRequest) (*model.GetArticlesResponse, error) {
    args := m.Called(ctx, req)
    return args.Get(0).(*model.GetArticlesResponse), args.Error(1)
}

func TestCreateArticle_Integration(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    mockService := new(MockArticleService)
    articleHandler := handler.NewArticleHandler(mockService)
    
    router := gin.New()
    router.POST("/articles", articleHandler.CreateArticle)

    req := &model.CreateArticleRequest{
        AuthorID: "author-123",
        Title:    "Test Article",
        Body:     "Test Body",
    }

    expectedArticle := &model.Article{
        ID:      "1",
        AuthorID: "author-123",
        Title:   "Test Article",
        Body:    "Test Body",
    }

    mockService.On("CreateArticle", mock.Anything, req).Return(expectedArticle, nil)

    body, _ := json.Marshal(req)
    w := httptest.NewRecorder()
    request, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(body))
    request.Header.Set("Content-Type", "application/json")

    router.ServeHTTP(w, request)

    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    
    data := response["data"].(map[string]interface{})
    assert.Equal(t, float64(1), data["id"])
    assert.Equal(t, req.Title, data["title"])
    assert.Equal(t, req.Body, data["body"])
    assert.Equal(t, req.AuthorID, data["author"])
    
    mockService.AssertExpectations(t)
}
