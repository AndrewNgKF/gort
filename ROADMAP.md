# Go RT Framework - Development Roadmap

## 🚀 Phase 1: Foundation (MVP)

### CLI Infrastructure

- [x] Project setup with cobra CLI framework
- [x] `gort new <app_name>` - Create new application with full structure
- [x] `gort version` - Show framework version
- [ ] `gort server` - Start development server
- [ ] Hot reload integration with Air

### Routing & HTTP

- [x] Core router implementation with Rails-like DSL
- [x] RESTful routing (`r.Resources()`)
- [x] Named routes and URL helpers
- [x] Route parameter extraction (`:id`, `*path`)
- [x] Namespace and nested resources support
- [x] `gort routes` - Display all routes command

### Controllers

- [x] Base controller with common methods
- [x] `gort generate controller` - Controller generator
- [x] Render methods (JSON, HTML, text)
- [x] Parameter binding (JSON, form, query)
- [x] Response helpers (NotFound, BadRequest, etc.)
- [ ] Flash message system
- [ ] Session management

### Models & Database

- [x] Database connection manager (PostgreSQL, MySQL, SQLite)
- [x] `gort generate model` - Model generator with fields
- [x] Basic CRUD operations (Create, Read, Update, Delete)
- [x] Query builder (Where, Order, Limit, etc.)
- [x] Model timestamps (created_at, updated_at)
- [x] Lean model files (Rails-style - no CRUD methods in model)
- [x] Public database API (gort.Find, gort.Create, etc.)
- [ ] `gort db:create` and `gort db:drop`

### Migrations

- [x] Migration file structure and naming
- [x] `gort generate migration` - Migration generator
- [ ] Schema DSL (CreateTable, AddColumn, etc.)
- [x] `gort db:migrate` - Run migrations
- [x] `gort db:rollback` - Rollback migrations
- [x] Migration version tracking table (schema_migrations)
- [x] `gort db:reset` - Drop, create, migrate
- [x] `gort db:schema:dump` - Generate schema.go file
- [x] Auto-update schema.go after migrations

### Views & Templates

- [x] Template engine wrapper for html/template
- [x] Layout system (application.html default)
- [ ] Partial rendering (`{{ render "partial" }}`)
- [ ] View helpers (formatDate, truncate, etc.)
- [x] Template auto-reload in development
- [x] `gort generate scaffold` - Full CRUD generator (needs method signature refinement)

### Testing

- [ ] Test helpers and assertions
- [ ] Controller testing utilities
- [ ] Model testing utilities
- [ ] `gort test` - Run test suite
- [ ] Test database management

---

## 🎨 Phase 2: Enhancement

### Advanced ORM Features

- [ ] Associations: `belongs_to`
- [ ] Associations: `has_many`
- [ ] Associations: `has_one`
- [ ] Associations: `many_to_many` with join tables
- [ ] Eager loading (prevent N+1 queries)
- [ ] Query scopes
- [ ] Model validations (presence, format, length, uniqueness)
- [ ] Validation error messages
- [ ] Callbacks (BeforeSave, AfterCreate, etc.)
- [ ] Soft deletes (deleted_at)
- [ ] Custom query methods

### Middleware System

- [ ] Middleware pipeline architecture
- [ ] Built-in Logger middleware
- [ ] Built-in Recovery middleware
- [ ] CORS middleware
- [ ] Request timeout middleware
- [ ] Rate limiting middleware
- [ ] `gort generate middleware` - Custom middleware generator
- [ ] Before/After filters in controllers

### Security

- [ ] CSRF protection
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS protection in templates
- [ ] Secure session cookies
- [ ] Password hashing utilities (bcrypt)
- [ ] Environment variable management (.env support)

### Form Helpers

- [ ] Form builder for templates
- [ ] Input field helpers (text, email, password, etc.)
- [ ] Select/dropdown helpers
- [ ] Checkbox and radio helpers
- [ ] Form error display
- [ ] CSRF token injection

### Asset Management

- [ ] Static file serving from `public/`
- [ ] Asset fingerprinting for caching
- [ ] CSS/JS concatenation
- [ ] Development vs production asset modes

### Database Seeds

- [ ] Seed file structure
- [ ] `gort db:seed` - Run seeds
- [ ] Faker library integration for dummy data

---

## 🌟 Phase 3: Polish & Advanced Features

### Console/REPL

- [ ] Interactive console with app context
- [ ] `gort console` - Start REPL
- [ ] Access to models and database
- [ ] Helper commands in console
- [ ] Load development/production environment

### Background Jobs

- [ ] Job queue system
- [ ] `gort generate job` - Job generator
- [ ] Worker process management
- [ ] Retry logic and error handling
- [ ] Scheduled jobs (cron-like)
- [ ] Job monitoring/dashboard

### WebSockets & Real-time

- [ ] WebSocket connection handling
- [ ] Channel/room system (like ActionCable)
- [ ] Broadcast helpers
- [ ] `gort generate channel` - Channel generator
- [ ] Authentication for WebSocket connections

### Mailers

- [ ] Email sending abstraction
- [ ] `gort generate mailer` - Mailer generator
- [ ] HTML and text email templates
- [ ] Email preview in development
- [ ] SMTP/SendGrid/AWS SES support
- [ ] Background email sending

### Authentication & Authorization

- [ ] `gort generate auth` - Full auth scaffold
- [ ] User registration/login/logout
- [ ] Password reset flow
- [ ] Session-based authentication
- [ ] JWT token support
- [ ] Authorization helpers (current_user, etc.)
- [ ] Role-based access control

### API Features

- [ ] API versioning support
- [ ] JSON API serializers
- [ ] API authentication (token-based)
- [ ] Rate limiting for API endpoints
- [ ] API documentation generator
- [ ] Pagination helpers

### Internationalization (i18n)

- [ ] Multi-language support
- [ ] Translation file structure (YAML)
- [ ] Template translation helpers
- [ ] Locale detection and switching
- [ ] Pluralization rules

### Caching

- [ ] Cache store abstraction (memory, Redis)
- [ ] Fragment caching in views
- [ ] Query result caching
- [ ] Cache invalidation helpers
- [ ] HTTP caching headers

### Logging & Monitoring

- [ ] Structured logging
- [ ] Log levels (debug, info, warn, error)
- [ ] Request logging with timing
- [ ] Query logging
- [ ] Log rotation
- [ ] Integration with monitoring tools

---

## 📦 Phase 4: Production & Deployment

### Build & Deploy

- [ ] `gort build` - Production binary builder
- [ ] Environment configuration management
- [ ] Database connection pooling
- [ ] Graceful shutdown handling
- [ ] Health check endpoints
- [ ] Docker support and Dockerfile generator

### Performance

- [ ] Database connection pooling optimization
- [ ] Response compression (gzip)
- [ ] HTTP/2 support
- [ ] Static asset CDN integration
- [ ] Query performance monitoring
- [ ] Memory profiling tools

### Documentation

- [ ] Comprehensive getting started guide
- [ ] API reference documentation
- [ ] Generator documentation
- [ ] Migration guide from Rails
- [ ] Video tutorials
- [ ] Example applications (blog, e-commerce, etc.)

### Community & Ecosystem

- [ ] Plugin system architecture
- [ ] Community plugin registry
- [ ] Contribution guidelines
- [ ] Code of conduct
- [ ] Issue templates
- [ ] CI/CD pipeline setup

---

## 🎯 Version Milestones

### v0.1.0 - Alpha (MVP)

Complete Phase 1 - Basic CRUD applications possible

### v0.5.0 - Beta

Complete Phase 2 - Production-ready for simple apps

### v1.0.0 - Stable Release

Complete Phase 3 - Full-featured framework

### v2.0.0 - Enterprise

Complete Phase 4 - Enterprise-grade features

---

## 📊 Progress Tracking

**Overall Completion: 20%**

- Phase 1: 29/48 tasks complete (60%)
- Phase 2: 0/38 tasks complete (0%)
- Phase 3: 0/42 tasks complete (0%)
- Phase 4: 0/18 tasks complete (0%)

**Total: 29/146 tasks complete**

### Recent Achievements ✨

- ✅ Built complete CLI with Cobra
- ✅ Implemented Rails-like router with RESTful routing
- ✅ Created controller and model generators
- ✅ Set up template rendering with layouts
- ✅ Database connection manager for PostgreSQL/MySQL/SQLite
- ✅ Full migration system (migrate, rollback, reset)
- ✅ Rails-style lean models with public API
- ✅ schema.go generation and auto-updates

---

Last Updated: December 5, 2025
