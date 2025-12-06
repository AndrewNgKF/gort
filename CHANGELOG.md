# Changelog

All notable changes to this project will be documented in this file.

## [0.1.0] - 2025-12-06

### Added

- Initial release of Gort framework
- CLI with `gort new` command to create new applications
- Scaffold generator (`gort generate scaffold`) for complete CRUD resources
- Model generator with database migrations
- Controller generator
- RESTful routing with `Resources()` method
- Database ORM with support for PostgreSQL, MySQL, and SQLite
- Automatic timestamps (created_at, updated_at)
- Migration system (migrate, rollback, reset, schema:dump)
- Template engine with layouts and partials
- Method override middleware for PUT/DELETE in HTML forms
- Base controller with Render, JSON, Redirect helpers
- Route printing with `gort routes`

### Features

- 7 RESTful actions per resource (Index, New, Show, Edit, Create, Update, Destroy)
- Rails-like project structure (app/, config/, db/, public/)
- Database configuration via YAML
- Form helpers with automatic CSRF protection via method override
- Development-friendly error messages

[0.1.0]: https://github.com/AndrewNgKF/gort/releases/tag/v0.1.0
