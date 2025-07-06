// package handler

// import (
//     "net/http"

//     "article-service/internal/model"
//     "article-service/internal/service"

//     "github.com/gin-gonic/gin"
// )

// type ArticleHandler struct {
//     service service.ArticleService
// }

// func NewArticleHandler(service service.ArticleService) *ArticleHandler {
//     return &ArticleHandler{
//         service: service,
//     }
// }

// func (h *ArticleHandler) CreateArticle(c *gin.Context) {
//     var req model.CreateArticleRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     article, err := h.service.CreateArticle(c.Request.Context(), &req)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
//         return
//     }

//     c.JSON(http.StatusCreated, gin.H{"data": article})
// }

// func (h *ArticleHandler) GetArticles(c *gin.Context) {
//     var req model.GetArticlesRequest
//     if err := c.ShouldBindQuery(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     response, err := h.service.GetArticles(c.Request.Context(), &req)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"data": response})
// }

package handler

import (
    "net/http"

    "article-service/internal/model"
    "article-service/internal/service"

    "github.com/gin-gonic/gin"
)

type ArticleHandler struct {
    service service.ArticleService
}

func NewArticleHandler(service service.ArticleService) *ArticleHandler {
    return &ArticleHandler{
        service: service,
    }
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    var req model.CreateArticleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    article, err := h.service.CreateArticle(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create article"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"data": article})
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
    var req model.GetArticlesRequest
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := h.service.GetArticles(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"data": response})
}