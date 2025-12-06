# Publishing Guide

## How to Publish Gort to the Public

Unlike npm or pip, Go doesn't have a central package registry. Instead, Go modules are distributed directly via version control (typically GitHub).

## Prerequisites

1. **GitHub Account** - Create one at https://github.com
2. **Git Installed** - `brew install git` (macOS) or download from git-scm.com
3. **Clean Repository** - All changes committed, no sensitive data

## Step-by-Step Publishing

### 1. Create GitHub Repository

```bash
# On GitHub.com:
# 1. Click "New Repository"
# 2. Name it "gort"
# 3. Make it public
# 4. Don't initialize with README (we have one)
```

### 2. Update Module Path

Already done! The `go.mod` should have:

```
module github.com/yourusername/gort
```

Replace `yourusername` with your actual GitHub username.

### 3. Initialize Git and Push

```bash
cd go_rt_framework

# Initialize git
git init

# Add all files
git add .

# Create first commit
git commit -m "Initial commit - Gort v0.1.0

- Rails-inspired Go web framework
- Complete CRUD scaffolding
- RESTful routing
- Database ORM (PostgreSQL, MySQL, SQLite)
- Migration system
- Template engine with layouts"

# Add GitHub remote (replace with your URL)
git remote add origin https://github.com/yourusername/gort.git

# Push to GitHub
git branch -M main
git push -u origin main
```

### 4. Create a Release Tag

```bash
# Tag the release
git tag -a v0.1.0 -m "Release v0.1.0 - Initial public release"

# Push the tag
git push origin v0.1.0
```

### 5. Create GitHub Release (Optional but Recommended)

1. Go to your repository on GitHub
2. Click "Releases" → "Create a new release"
3. Select tag `v0.1.0`
4. Title: "Gort v0.1.0 - Initial Release"
5. Description: Copy from CHANGELOG.md
6. Click "Publish release"

## How Users Install It

After publishing, anyone can install Gort:

```bash
# Install latest version
go install github.com/yourusername/gort/cmd/gort@latest

# Or specific version
go install github.com/yourusername/gort/cmd/gort@v0.1.0

# Verify
gort version
```

## Automatic Features

Once on GitHub, you get:

### 1. pkg.go.dev Documentation

- Automatically indexed at `https://pkg.go.dev/github.com/yourusername/gort`
- API documentation generated from your code
- Updates when you push new versions
- Free forever!

### 2. Go Proxy Caching

- Your module is cached at `proxy.golang.org`
- Fast downloads worldwide
- Immutable versions (v0.1.0 never changes)
- No action needed - automatic!

### 3. Module Checksums

- Security via `sum.golang.org`
- Prevents tampering
- Automatic integrity checks

## Publishing Future Updates

```bash
# Make your changes
git add .
git commit -m "Add new feature"
git push

# Create new version tag
git tag -a v0.2.0 -m "Release v0.2.0"
git push origin v0.2.0

# Users update with:
# go install github.com/yourusername/gort/cmd/gort@latest
```

## Versioning Rules

Go uses Semantic Versioning (semver):

- `v0.1.0` → `v0.2.0` - Minor version (new features, backward compatible)
- `v0.2.0` → `v0.2.1` - Patch version (bug fixes)
- `v0.9.0` → `v1.0.0` - Major version (breaking changes)

## Important Notes

1. **Never delete tags** - They're cached forever in the Go proxy
2. **Never force push** - Breaks user's `go.sum` checksums
3. **Use semantic versioning** - Required for Go modules
4. **Test before tagging** - Tags are immutable!

## Comparison with Other Ecosystems

| Feature        | npm           | pip            | Go Modules             |
| -------------- | ------------- | -------------- | ---------------------- |
| Registry       | npmjs.com     | pypi.org       | No central registry    |
| Hosting        | Centralized   | Centralized    | GitHub/GitLab/etc      |
| Publishing     | `npm publish` | `twine upload` | `git push` + `git tag` |
| Installation   | `npm install` | `pip install`  | `go install`           |
| Cost           | Free          | Free           | Free                   |
| Account Needed | npm account   | PyPI account   | GitHub/Git only        |

## Benefits of Go's Approach

- ✅ No separate account/registry needed
- ✅ Version control = package distribution
- ✅ Free forever (just needs GitHub)
- ✅ Can use private repos (with auth)
- ✅ Automatic caching worldwide
- ✅ Built-in security checksums
- ✅ No "npm left-pad" incidents

## Need Help?

Check the README for installation instructions or open an issue on GitHub!
