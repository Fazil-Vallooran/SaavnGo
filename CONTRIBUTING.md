# Contributing to JioSaavn API

Thank you for considering contributing to this project! Here are some guidelines to help you get started.

## Code of Conduct

- Be respectful and considerate
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Keep discussions on topic

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Git
- Basic understanding of REST APIs
- Familiarity with the Gin framework (helpful but not required)

### Setting Up Development Environment

1. Fork the repository

2. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/jioSaavnAPI.git
cd jioSaavnAPI
```

3. Install dependencies:
```bash
make deps
```

4. Run the application:
```bash
make run
```

5. Run tests:
```bash
make test
```

## Development Workflow

### Branching Strategy

- `main` - Stable production branch
- `develop` - Development branch
- `feature/*` - Feature branches
- `bugfix/*` - Bug fix branches
- `hotfix/*` - Urgent fixes

### Making Changes

1. Create a new branch:
```bash
git checkout -b feature/your-feature-name
```

2. Make your changes following the coding standards

3. Add tests for your changes

4. Run tests and ensure they pass:

5. Format your code:

6. Commit your changes:
```bash
git add .
git commit -m "feat: add new feature"
```

7. Push to your fork:
```bash
git push origin feature/your-feature-name
```

8. Open a Pull Request

## Coding Standards

### Go Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting (run `make fmt`)
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

### Example:

```go
// GetSongByID fetches song details from JioSaavn API
func GetSongByID(id string) (*Song, error) {
    if id == "" {
        return nil, errors.New("song ID cannot be empty")
    }
    
    // Implementation
    return song, nil
}
```

### Project Structure

```
jioSaavnAPI/
├── config/          # Configuration management
├── middleware/      # HTTP middleware
├── models/          # Data models
├── routes/          # Route definitions
├── services/        # Business logic
├── utils/           # Helper functions
└── main.go          # Entry point
```

### Adding New Endpoints

1. Define the model in `models/` (if needed)
2. Create the service handler in `services/`
3. Register the route in `routes/`
4. Add tests
5. Update documentation

Example:

```go
// In services/jiosaavn.go
func GetPlaylistHandler(c *gin.Context) {
    id := c.Param("id")
    // Implementation
}

// In routes/song.go
r.GET("/playlist/:id", services.GetPlaylistHandler)
```

## Testing Guidelines

### Writing Tests

- Write tests for all new functions
- Aim for good test coverage
- Use table-driven tests where appropriate
- Test both success and failure cases

### Example Test:

```go
func TestGetSongByID(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        wantErr bool
    }{
        {
            name:    "Valid ID",
            id:      "abc123",
            wantErr: false,
        },
        {
            name:    "Empty ID",
            id:      "",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := GetSongByID(tt.id)
            if (err != nil) != tt.wantErr {
                t.Errorf("GetSongByID() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./services -v
```

## Documentation

- Update README.md for user-facing changes
- Update API_EXAMPLES.md for new endpoints
- Add inline comments for complex logic
- Update CHANGELOG.md

## Commit Message Guidelines

Follow the Conventional Commits specification:

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

Examples:
```
feat: add playlist endpoint
fix: handle nil pointer in decrypt function
docs: update API examples
test: add tests for format utils
```

## Pull Request Process

1. Ensure all tests pass
2. Update documentation
3. Add a clear description of changes
4. Reference any related issues
5. Wait for review and address feedback

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe how you tested your changes

## Checklist
- [ ] Code follows project style
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] All tests pass
```

## Reporting Bugs

### Before Reporting

- Check if the bug has already been reported
- Try to reproduce the bug
- Gather relevant information (logs, environment, etc.)

### Bug Report Template

```markdown
**Describe the bug**
A clear description of the bug

**To Reproduce**
Steps to reproduce:
1. Call endpoint '...'
2. With parameters '...'
3. See error

**Expected behavior**
What you expected to happen

**Environment**
- OS: [e.g., Ubuntu 20.04]
- Go version: [e.g., 1.21]
- API version: [e.g., 2.0.0]

**Additional context**
Any other relevant information
```

## Feature Requests

We welcome feature requests! Please:

1. Check if the feature has already been requested
2. Provide a clear use case
3. Explain the expected behavior
4. Consider submitting a PR if possible

## Questions?

- Open an issue for questions
- Tag with `question` label
- Be specific and provide context

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

## Recognition

Contributors will be recognized in the project README. Thank you for your contributions! 🎉
