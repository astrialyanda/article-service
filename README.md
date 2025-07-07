# Article Service

An article management service built with Go, featuring a REST API and PostgreSQL access.

## Features

- **RESTful API**:
  - `POST /api/v1/articles` - Create new articles
  - `GET /api/v1/articles` - Get articles with filtering and pagination

- **Advanced Query Capabilities**:
  - Full-text search in title and body
  - Author name filtering
  - Pagination support
  - Sorting by latest first

## Executing the program

### Using Docker Compose

```bash
# Clone the repository
git clone <your-repo-url>
cd article-service

# Start the database service
docker-compose up -d --build db

# Perform migration in the container
docker-compose exec db psql -U user -d articledb -f /migrations/001_create_articles_table.sql

# Create user database
docker-compose exec db psql -U user -d articledb
    # in the psql terminal
    CREATE DATABASE "user";

# Start all services
docker-compose up --build

# The API will be available at http://localhost:8080
```

## API Documentation

### Create Article

```bash
curl -X POST http://localhost:8080/api/v1/articles \
  -H "Content-Type: application/json" \
  -d '{
    "title": "My First Article",
    "body": "This is the content of my article...",
    "author_id": "author-123"
  }'
```

### Get Articles

Accessing get articles endpoint should be preceded by seeding the author table, can be done by manually seed on the psql terminal.

```bash
# Get all articles
curl http://localhost:8080/api/v1/articles

# Search articles by keyword and filter by author name
curl "http://localhost:8080/api/v1/articles?query=Article&author_name=John%20Doe&page=1&limit=10"
```

## Database Schema

```sql
CREATE TABLE IF NOT EXISTS authors (
    id TEXT PRIMARY KEY,
    name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS articles (
    id TEXT PRIMARY KEY,
    author_id TEXT NOT NULL REFERENCES authors(id),
    title VARCHAR(200) NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_articles_created_at ON articles(created_at DESC);
CREATE INDEX idx_articles_author_id ON articles(author_id);
CREATE INDEX idx_articles_title_gin ON articles USING gin(to_tsvector('english', title));
CREATE INDEX idx_articles_body_gin ON articles USING gin(to_tsvector('english', body));
```

## Architecture

The service follows clean architecture principles:

- **Handler Layer**: HTTP request/response handling
- **Service Layer**: Business logic
- **Repository Layer**: Data access (with SQL JOIN for author name)
- **Model Layer**: Data structures and DTOs

## Performance Considerations

1. **Database Optimization**:
   - GIN indexes for full-text search
   - B-tree indexes for sorting and filtering
   - Connection pooling

2. **Pagination**:
   - Limit-offset pagination
   - Configurable page sizes (max 100)
   - Total count optimization

## Testing

```bash
# Run unit tests
go test ./tests/unit/...

# Run integration tests
go test ./tests/integration/...

## Environment Variables

- `PORT`: Server port (default: 8080)
- `DATABASE_URL`: PostgreSQL connection string
