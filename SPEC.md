# Go RT Framework - Technical Specification

## Framework Name

**Go RT** (Go Ruby-like Toolkit / Go RealTime / Go Rails-inspired Toolkit)

---

## 🎯 Core Philosophy

- **Convention over Configuration** - Opinionated structure like Rails
- **Developer Happiness** - Intuitive commands and generators
- **Performance First** - Leverage Go's speed and concurrency
- **Batteries Included** - Everything you need out of the box

---

## 📊 Database Support

### Primary Support (v1.0)

- **PostgreSQL** (default, recommended)
- **MySQL/MariaDB**
- **SQLite** (development/testing)

### ORM Layer

- Custom query builder inspired by ActiveRecord
- Support for associations: `has_many`, `belongs_to`, `has_one`, `many_to_many`
- Automatic timestamps: `created_at`, `updated_at`
- Soft deletes: `deleted_at` (optional)
- Validations at model level
- Query scopes and callbacks

### Migration System

- Timestamped migration files: `20251205143022_create_users.go`
- DSL for schema definition
- Rollback support
- Schema versioning table

---

## 🎮 CLI Commands

### Project Management

```bash
gort new <app_name>              # Create new application
gort server                      # Start development server (hot reload)
gort console                     # Interactive REPL with app context
gort version                     # Show framework version
```

### Generators

```bash
gort generate model <name> [fields]
  # Example: gort generate model User name:string email:string:unique age:integer

gort generate controller <name> [actions]
  # Example: gort generate controller Posts index show create update destroy

gort generate migration <name>
  # Example: gort generate migration add_role_to_users role:string

gort generate scaffold <name> [fields]
  # Full CRUD: model + migration + controller + views + routes

gort generate middleware <name>
  # Example: gort generate middleware Auth
```

### Database Commands

```bash
gort db:create                   # Create database
gort db:drop                     # Drop database
gort db:migrate                  # Run pending migrations
gort db:rollback [steps]         # Rollback migrations
gort db:reset                    # Drop, create, and migrate
gort db:seed                     # Run seed file
gort db:schema:dump              # Export schema
```

### Other Commands

```bash
gort routes                      # List all routes
gort test                        # Run test suite
gort build                       # Build production binary
gort deploy                      # Deploy helpers (future)
```

---

## 📁 Project Structure

```
myapp/
├── app/
│   ├── controllers/
│   │   ├── application_controller.go
│   │   └── posts_controller.go
│   ├── models/
│   │   ├── post.go
│   │   └── user.go
│   ├── views/
│   │   ├── layouts/
│   │   │   └── application.html
│   │   └── posts/
│   │       ├── index.html
│   │       ├── show.html
│   │       └── _form.html (partial)
│   └── middleware/
│       └── auth.go
├── config/
│   ├── application.go           # App configuration
│   ├── database.yml             # DB config (multi-env)
│   ├── routes.go                # Route definitions
│   └── environments/
│       ├── development.go
│       ├── test.go
│       └── production.go
├── db/
│   ├── migrations/
│   │   └── 20251205143022_create_users.go
│   ├── seeds.go
│   └── schema.go (auto-generated)
├── public/
│   ├── assets/
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   └── uploads/
├── test/
│   ├── controllers/
│   ├── models/
│   └── integration/
├── tmp/
│   └── pids/
├── vendor/
├── main.go                      # Application entry point
├── go.mod
├── go.sum
└── README.md
```

---

## 🎨 Code Examples

### Model File (`app/models/user.go`)

```go
package models

import (
    "time"
    "github.com/yourusername/gort"
)

type User struct {
    gort.Model
    ID        int64     `gort:"primary_key"`
    Name      string    `gort:"type:varchar(255);not_null"`
    Email     string    `gort:"type:varchar(255);unique;not_null"`
    Age       int       `gort:"type:integer"`
    CreatedAt time.Time `gort:"autoCreateTime"`
    UpdatedAt time.Time `gort:"autoUpdateTime"`

    // Associations
    Posts []Post `gort:"has_many"`
}

// Validations
func (u *User) Validate() error {
    return gort.Validate(u,
        gort.ValidatePresence("Name", "Email"),
        gort.ValidateFormat("Email", gort.EmailRegex),
        gort.ValidateLength("Name", gort.Min(2), gort.Max(100)),
        gort.ValidateUniqueness("Email"),
    )
}

// Callbacks
func (u *User) BeforeSave() error {
    // Normalize email
    u.Email = strings.ToLower(u.Email)
    return nil
}

// Scopes
func (u *User) Active() *gort.Query {
    return gort.Where("age > ?", 18)
}
```

### Controller File (`app/controllers/posts_controller.go`)

```go
package controllers

import (
    "myapp/app/models"
    "github.com/yourusername/gort"
)

type PostsController struct {
    gort.Controller
}

// GET /posts
func (c *PostsController) Index() {
    var posts []models.Post
    gort.Model(&models.Post{}).FindAll(&posts)

    c.Render("posts/index", gort.H{
        "posts": posts,
        "title": "All Posts",
    })
}

// GET /posts/:id
func (c *PostsController) Show() {
    id := c.Params["id"]
    var post models.Post

    if err := gort.Model(&models.Post{}).Find(id, &post); err != nil {
        c.NotFound()
        return
    }

    c.Render("posts/show", gort.H{"post": post})
}

// POST /posts
func (c *PostsController) Create() {
    var post models.Post

    if err := c.BindJSON(&post); err != nil {
        c.BadRequest(err.Error())
        return
    }

    if err := post.Save(); err != nil {
        c.UnprocessableEntity(gort.H{"errors": err})
        return
    }

    c.RedirectTo("/posts/" + post.ID)
}

// PUT /posts/:id
func (c *PostsController) Update() {
    id := c.Params["id"]
    var post models.Post

    if err := gort.Model(&models.Post{}).Find(id, &post); err != nil {
        c.NotFound()
        return
    }

    if err := c.BindJSON(&post); err != nil {
        c.BadRequest(err.Error())
        return
    }

    if err := post.Update(); err != nil {
        c.UnprocessableEntity(gort.H{"errors": err})
        return
    }

    c.JSON(gort.H{"post": post})
}

// DELETE /posts/:id
func (c *PostsController) Destroy() {
    id := c.Params["id"]
    var post models.Post

    if err := gort.Model(&models.Post{}).Find(id, &post); err != nil {
        c.NotFound()
        return
    }

    post.Destroy()
    c.RedirectTo("/posts")
}
```

### Routes File (`config/routes.go`)

```go
package config

import (
    "myapp/app/controllers"
    "github.com/yourusername/gort"
)

func Routes(r *gort.Router) {
    // Root route
    r.Root("home#index")

    // RESTful resources
    r.Resources("posts", &controllers.PostsController{})
    // Generates:
    // GET    /posts          -> Index
    // GET    /posts/:id      -> Show
    // POST   /posts          -> Create
    // PUT    /posts/:id      -> Update
    // DELETE /posts/:id      -> Destroy

    // Custom routes
    r.Get("/about", "pages#about")
    r.Post("/login", "auth#login")
    r.Delete("/logout", "auth#logout")

    // Nested resources
    r.Resources("users", func(r *gort.Router) {
        r.Resources("posts", &controllers.PostsController{})
    })
    // Generates: /users/:user_id/posts, etc.

    // Namespace
    r.Namespace("api/v1", func(r *gort.Router) {
        r.Resources("posts", &controllers.API.V1.PostsController{})
    })

    // Middleware
    r.Use(gort.Logger())
    r.Use(gort.Recovery())
    r.Use(middleware.Auth())
}
```

### Migration File (`db/migrations/20251205143022_create_users.go`)

```go
package migrations

import "github.com/yourusername/gort/migration"

func init() {
    migration.Register("20251205143022", Up, Down)
}

func Up(m *migration.Migration) {
    m.CreateTable("users", func(t *migration.Table) {
        t.ID()
        t.String("name").NotNull()
        t.String("email").Unique().NotNull()
        t.Integer("age")
        t.Timestamps()
    })

    m.AddIndex("users", "email")
}

func Down(m *migration.Migration) {
    m.DropTable("users")
}
```

### View File (`app/views/posts/index.html`)

```html
{{ define "posts/index" }} {{ template "layouts/application" . }} {{ define
"content" }}
<h1>{{ .title }}</h1>

<div class="posts">
  {{ range .posts }}
  <article class="post">
    <h2><a href="/posts/{{ .ID }}">{{ .Title }}</a></h2>
    <p>{{ .Body }}</p>
    <small>Posted {{ .CreatedAt | formatDate }}</small>
  </article>
  {{ end }}
</div>

<a href="/posts/new" class="btn">New Post</a>
{{ end }} {{ end }}
```

### Config File (`config/database.yml`)

```yaml
development:
  adapter: postgresql
  database: myapp_development
  host: localhost
  port: 5432
  username: postgres
  password: postgres
  pool: 5

test:
  adapter: postgresql
  database: myapp_test
  host: localhost
  port: 5432
  username: postgres
  password: postgres

production:
  adapter: postgresql
  database: myapp_production
  host: <%= ENV["DB_HOST"] %>
  port: 5432
  username: <%= ENV["DB_USER"] %>
  password: <%= ENV["DB_PASSWORD"] %>
  pool: 25
```

---

## 🔥 Framework Features (Rails Comparison)

| Feature                   | Rails                    | Go RT                         | Status     |
| ------------------------- | ------------------------ | ----------------------------- | ---------- |
| **CLI Generator**         | `rails new`              | `gort new`                    | ✅ Planned |
| **Model Generators**      | `rails g model`          | `gort generate model`         | ✅ Planned |
| **Controller Generators** | `rails g controller`     | `gort generate controller`    | ✅ Planned |
| **Migrations**            | ActiveRecord migrations  | Timestamp-based Go migrations | ✅ Planned |
| **ORM**                   | ActiveRecord             | Custom ORM with AR-like API   | ✅ Planned |
| **Routing DSL**           | `resources`, `namespace` | Identical syntax              | ✅ Planned |
| **View Templates**        | ERB                      | Go templates                  | ✅ Planned |
| **Asset Pipeline**        | Sprockets                | Modern Go asset handling      | 🔄 Future  |
| **Testing Framework**     | Minitest/RSpec           | Built-in Go testing + helpers | ✅ Planned |
| **Console/REPL**          | `rails console`          | `gort console`                | 🔄 Future  |
| **Hot Reload**            | Built-in                 | Air integration               | ✅ Planned |
| **Background Jobs**       | ActiveJob/Sidekiq        | Go routines + queue           | 🔄 Future  |
| **WebSockets**            | ActionCable              | Built-in channels             | 🔄 Future  |
| **Mailers**               | ActionMailer             | Email helpers                 | 🔄 Future  |
| **Authentication**        | Devise                   | Built-in auth generator       | 🔄 Future  |

---

## 🎯 Core Dependencies

- **CLI**: `spf13/cobra` - Command line interface
- **Router**: Custom (inspired by Chi/Gin but more Rails-like)
- **Database**: `database/sql` with custom layer
- **Templates**: `html/template` with helpers
- **Config**: `spf13/viper` - Configuration management
- **Hot Reload**: `cosmtrek/air` - Development auto-reload
- **Testing**: Standard `testing` package with helpers
- **Validation**: Custom validation engine

---

## 🎁 Additional Awesome Features

### Developer Experience

- **Interactive Generators**: Prompt for fields if not provided
- **Code Snippets**: VS Code extension with gort snippets
- **Live Reload**: Browser auto-refresh on file changes
- **Error Pages**: Beautiful error pages in development (like Better Errors)
- **SQL Query Display**: Show executed queries in dev mode
- **Debug Toolbar**: Request/response inspector in browser

### Performance Features

- **Query Optimization**: Automatic N+1 query detection and warnings
- **Connection Pooling**: Configurable database connection pools
- **Response Caching**: Built-in HTTP caching with ETags
- **Asset Compilation**: Minify and bundle CSS/JS for production
- **Memory Profiling**: Built-in memory leak detection
- **Concurrency**: Easy goroutine management for parallel processing

### Modern Web Features

- **Turbo/Hotwire**: Server-rendered SPA experience (like Turbo)
- **LiveView**: Real-time updates without JavaScript (Phoenix LiveView style)
- **Progressive Web App**: PWA manifest and service worker generators
- **GraphQL Support**: Optional GraphQL API layer
- **gRPC Support**: Generate gRPC services alongside REST

### CLI Enhancements

- **Autocomplete**: Shell completion for zsh/bash
- **Project Templates**: `gort new --template=api` (API-only, full-stack, etc.)
- **Undo Generator**: `gort destroy model User` (reverse generators)
- **Diff Preview**: Preview file changes before generating
- **Annotations**: `gort annotate` - add schema comments to models

### Database Goodies

- **Multi-DB Support**: Connect to multiple databases
- **Read Replicas**: Automatic read/write splitting
- **Database Tasks**: `gort db:structure:dump`, `gort db:schema:load`
- **Fixtures**: Load test data from YAML files
- **Database Encryption**: Transparent column encryption

### Authentication Options

- **OAuth Integration**: Google, GitHub, Facebook login
- **Two-Factor Auth**: TOTP/SMS 2FA support
- **Magic Links**: Passwordless authentication
- **API Keys**: Generate and manage API keys
- **Session Stores**: Redis, database, or memory sessions

### API Features

- **Auto Documentation**: Generate OpenAPI/Swagger docs from routes
- **API Versioning**: `/api/v1`, `/api/v2` namespace support
- **GraphQL Playground**: Interactive GraphQL explorer
- **Webhooks**: Easy webhook handler generation
- **Rate Limiting**: Per-route, per-user, per-IP limits

### Testing Tools

- **Factory Pattern**: Create test data easily (like FactoryBot)
- **Request Specs**: Integration test helpers
- **Parallel Tests**: Run tests concurrently
- **Coverage Reports**: Built-in test coverage
- **Mock Helpers**: Easy mocking for external services
- **Snapshot Testing**: Compare output against saved snapshots

### Deployment

- **Zero-Downtime Deploy**: Graceful restarts
- **Health Checks**: `/health` and `/ready` endpoints
- **Metrics Export**: Prometheus metrics endpoint
- **Docker**: `gort docker:init` - generate Dockerfile
- **Kubernetes**: Generate K8s manifests
- **Platform Configs**: Heroku, AWS, GCP deploy configs

### Monitoring & Observability

- **APM Integration**: New Relic, DataDog, etc.
- **Distributed Tracing**: OpenTelemetry support
- **Error Tracking**: Sentry, Rollbar integration
- **Query Analytics**: Slow query logging
- **Request Profiling**: Performance flame graphs

### Developer Tools

- **Schema Visualizer**: Generate ERD diagrams
- **Route Mapper**: Visual route explorer
- **Dependency Graph**: Show model relationships
- **Code Stats**: Lines of code, complexity metrics
- **Security Audit**: `gort audit` - check for vulnerabilities

---

_See ROADMAP.md for development phases and progress tracking._

---

## 📝 Naming Conventions

- **Files**: snake_case (`user_controller.go`, `create_users.go`)
- **Types/Structs**: PascalCase (`User`, `PostsController`)
- **Functions**: PascalCase for exported, camelCase for private
- **Database**: plural snake_case (`users`, `blog_posts`)
- **Routes**: RESTful conventions (`/users`, `/users/:id`)

---

## 🎓 Example Usage

```bash
# Create new app
gort new blog
cd blog

# Generate a Post model
gort generate model Post title:string body:text published:boolean

# Run migrations
gort db:migrate

# Generate controller
gort generate controller Posts index show create update destroy

# Start server
gort server
# Server running at http://localhost:3000
```

---

## 💡 Key Differentiators from Other Go Frameworks

1. **Full-stack opinionated** (not minimal like Echo/Gin)
2. **Generators for everything** (save typing, enforce conventions)
3. **Rails-like DX** (familiar for Ruby devs)
4. **Integrated ORM** (not just router + middleware)
5. **Convention over configuration** (sensible defaults)

---

Ready to build! 🚀
