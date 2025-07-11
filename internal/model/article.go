package model

import (
    "time"
)

type Article struct {
    ID        string    `json:"id" db:"id"`
    AuthorID  string    `json:"author_id" db:"author_id"`
    AuthorName string    `json:"author_name" db:"author_name"`
    Title     string    `json:"title" db:"title"`
    Body      string    `json:"body" db:"body"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Author struct {
    AuthorID string `json:"id" db:"id"`
    Name     string `json:"name" db:"name"`
}

type CreateArticleRequest struct {
    AuthorID string `json:"author_id" binding:"required"`
    Title    string `json:"title" binding:"required,min=1,max=200"`
    Body     string `json:"body" binding:"required,min=1"`
}

type GetArticlesRequest struct {
    Query    string `form:"query"`
    AuthorName string `form:"author_name"`
    Page     int    `form:"page,default=1"`
    Limit    int    `form:"limit,default=10"`
}

type GetArticlesResponse struct {
    Articles []Article `json:"articles"`
    Total    int       `json:"total"`
    Page     int       `json:"page"`
    Limit    int       `json:"limit"`
}