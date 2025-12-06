# Gort Framework

A Rails-inspired web framework for Go, bringing developer happiness and convention over configuration to the Go ecosystem.

## ⚡ Try It Now (5 Minutes)

```bash
# Install Gort
go install github.com/AndrewNgKF/gort/cmd/gort@latest

# Create your first app
gort new myblog && cd myblog

# Setup everything
go mod tidy && gort db:setup

# Generate a blog scaffold
gort generate scaffold Post title:string body:text published:boolean

# Run migrations and start server
gort db:migrate && gort server
```

🎉 Open http://localhost:3000 to see your app!

## Installation

### Option 1: Install Globally (Recommended)

Install Gort so you can use it from anywhere, just like Rails:

```bash
# Install from GitHub
go install github.com/AndrewNgKF/gort/cmd/gort@latest

# Or install from local source
cd go_rt_framework
go install ./cmd/gort

# Verify installation
gort version

# Now you can use 'gort' from any directory!
```

The `gort` command will be installed to `$GOPATH/bin` (usually `~/go/bin`). Make sure this is in your PATH:

```bash
# Add to ~/.zshrc or ~/.bashrc if needed
export PATH=$PATH:$HOME/go/bin
```

### Option 2: Development Mode

If you're working on the Gort framework itself:

```bash
cd go_rt_framework
go build -o bin/gort ./cmd/gort
# Use ./bin/gort for commands
```

## Quick Start

```bash
# Create a new application
gort new myblog
cd myblog

# Install dependencies
go mod tidy

# Setup database (creates database and runs migrations)
gort db:setup

# Generate a complete CRUD scaffold
gort generate scaffold Post title:string body:text published:boolean

# Run new migrations
gort db:migrate

# Start the server
gort server
```

Visit http://localhost:3000 to see the welcome page, or http://localhost:3000/posts to see your scaffold!

## Available Commands

### Application Generation

```bash
gort new <app_name>              # Create new application
```

### Generators

```bash
gort generate model <name> [fields...]           # Generate model + migration
gort generate controller <name> [actions...]     # Generate controller
gort generate scaffold <name> [fields...]        # Generate full CRUD resource
```

**Supported Field Types:**

| Type | SQL Type | Go Type | Use Case |
|------|----------|---------|----------|
| `string` | `VARCHAR(255)` | `string` | Short text (names, emails, titles) |
| `text` | `TEXT` | `string` | Long text (descriptions, content) |
| `integer` / `int` | `INTEGER` | `int` | Numbers |
| `bigint` / `int64` | `BIGINT` | `int64` | Large numbers |
| `boolean` / `bool` | `BOOLEAN` | `bool` | True/false flags |
| `float` / `decimal` | `DECIMAL(10,2)` | `float64` | Money, ratings, percentages |
| `datetime` / `timestamp` | `TIMESTAMP` | `time.Time` | Date and time |
| `date` | `DATE` | `time.Time` | Date only |
| `uuid` | `UUID` | `string` | Unique identifiers |
| `json` | `JSON` | `interface{}` | JSON data |
| `jsonb` | `JSONB` | `interface{}` | Binary JSON (faster, PostgreSQL) |
| `array` | `TEXT[]` | `[]string` | List of values |

**Field Modifiers:**
- `:unique` - Add unique constraint
- `:index` - Add database index
- `:nullable` / `:null` - Allow NULL values
- `:ref:ModelName` - Foreign key reference

**Examples:**
```bash
# Basic types
gort generate model User name:string bio:text age:integer active:boolean

# Advanced types with modifiers
gort generate model Article title:string:unique slug:string:index body:text views:bigint rating:float published_at:datetime tags:array metadata:jsonb

# With relationships
gort generate model Comment body:text user_id:integer:ref:User post_id:integer:ref:Post
```

### Database Commands

Both `db:command` (Rails style) and `db command` (Go style) syntaxes are supported!

```bash
# Database setup
gort db:create                   # Create the database
gort db:drop                     # Drop the database
gort db:setup                    # Create database and run migrations (recommended!)

# Migrations
gort db:migrate                  # Run pending migrations
gort db:rollback                 # Rollback last migration
gort db:reset                    # Reset database (rollback all + migrate)

# Schema
gort db:schema:dump              # Update schema.go with current database structure
```

### Server Command

```bash
gort server                      # Start the development server (default: port 3000)
gort server -p 8080              # Start server on custom port
gort s                           # Short alias
```

### Other Commands

```bash
gort routes                      # Display all registered routes
gort version                     # Show version information
```

## Project Structure

```
myapp/
├── app/
│   ├── controllers/            # Request handlers
│   │   ├── application_controller.go
│   │   └── posts_controller.go
│   ├── models/                 # Data models
│   │   └── post.go
│   ├── views/                  # HTML templates
│   │   ├── layouts/
│   │   │   └── application.html
│   │   └── posts/
│   │       ├── index.html
│   │       ├── show.html
│   │       ├── new.html
│   │       ├── edit.html
│   │       └── _form.html      # Partial template
│   └── middleware/             # Custom middleware
├── config/
│   ├── application.go          # App configuration
│   ├── database.yml            # Database settings
│   ├── routes.go               # Route definitions
│   └── environments/           # Environment-specific config
├── db/
│   ├── migrations/             # Database migrations
│   ├── schema.go               # Current database schema
│   └── seeds.go                # Seed data
├── public/
│   └── assets/                 # Static files (CSS, JS, images)
├── test/                       # Test files
└── main.go                     # Application entry point
```

## Documentation

See [SPEC.md](SPEC.md) for the complete specification and [ROADMAP.md](ROADMAP.md) for development progress.

## Features

✅ **CLI Generators** - Scaffold models, controllers, and full CRUD resources  
✅ **RESTful Routing** - Rails-like routing DSL with Resources()  
✅ **Database ORM** - CRUD operations with automatic timestamps  
✅ **Migration System** - Version-controlled database schema changes  
✅ **Template Engine** - Layouts and partials with html/template  
✅ **Method Override** - PUT/DELETE support via `_method` parameter  
✅ **Multi-Database** - PostgreSQL, MySQL, and SQLite support  
✅ **Rich Data Types** - Support for string, text, integer, bigint, boolean, float, datetime, date, uuid, json, jsonb, array  
✅ **Dev Server** - Built-in `gort server` command with hot reload support  
✅ **Rails-style Commands** - Both `db:migrate` and `db migrate` syntaxes supported  
✅ **Beautiful Welcome Page** - Minimal, Rails-inspired landing page with community globe illustration

## Example: Creating a Blog in 5 Commands

```bash
# 1. Create application
gort new blog
cd blog

# 2. Install dependencies
go mod tidy

# 3. Setup database (creates DB and runs migrations)
gort db:setup

# 4. Generate Post scaffold with rich field types
gort generate scaffold Post title:string body:text author:string:unique published:boolean views:integer rating:float published_at:datetime tags:array

# 5. Run new migrations and start server
gort db:migrate
gort server
```

Now visit:
- **http://localhost:3000** - Beautiful welcome page with Rails-style globe
- **http://localhost:3000/posts** - Full CRUD application with:
  - List all posts (Index)
  - Create new posts (New/Create)
  - View individual posts (Show)
  - Edit posts (Edit/Update)
  - Delete posts (Destroy)

## Development Workflow

### Working on the Gort Framework Itself

If you're developing the Gort framework and want to test changes:

1. **Make changes to Gort framework code**
2. **Rebuild Gort:** `go build -o bin/gort ./cmd/gort`
3. **In your test app, add a replace directive:**
   ```bash
   # In your app directory (e.g., blog/)
   echo "" >> go.mod
   echo "replace gort => ../path/to/go_rt_framework" >> go.mod
   ```
4. **Use the local version:** `../bin/gort generate scaffold ...`

This is similar to using `gem 'rails', path: '../rails'` in a Gemfile for Rails development.

### Using Published/Installed Gort

Once installed globally with `go install`, just use:

```bash
gort new myapp        # No ./ or ../bin/ prefix needed!
cd myapp
gort generate scaffold Post title:string body:text
gort db migrate
```

## Database Configuration

Edit `config/database.yml`:

```yaml
development:
  adapter: postgresql
  database: myapp_development
  host: localhost
  port: 5432
  username: your_username
  password: your_password
  sslmode: disable

production:
  adapter: postgresql
  database: myapp_production
  # ... production settings
```

## Routing Example

In `config/routes.go`:

```go
func Routes() *gort.Router {
    r := gort.NewRouter()

    // Middleware
    r.Use(gort.Logger())
    r.Use(gort.Recovery())

    // Root route
    r.Root(controllers.HomeIndex)

    // RESTful resources (generates 7 routes)
    r.Resources("posts", &controllers.PostController{})

    // Custom routes
    r.Get("/about", controllers.AboutPage)

    return r
}
```

## Model Example

```go
type Post struct {
    ID        int64     `db:"id"`
    Title     string    `db:"title"`
    Body      string    `db:"body"`
    Published bool      `db:"published"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

func (Post) TableName() string {
    return "posts"
}
```

## Controller Example

```go
func (c *PostController) Index(w http.ResponseWriter, r *http.Request) {
    var posts []models.Post

    if err := gort.FindAll("posts", &posts); err != nil {
        c.InternalServerError(w, err.Error())
        return
    }

    c.Render(w, r, "posts/index", map[string]interface{}{
        "posts": posts,
    })
}
```

## License

MIT
