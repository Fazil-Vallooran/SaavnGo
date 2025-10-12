# Quick Start Guide

Get up and running with JioSaavn API in less than 5 minutes!

## Installation

### Option 1: Run with Go (Recommended for Development)

```bash
# Clone the repository
git clone <repository-url>
cd jioSaavnAPI

# Install dependencies
go mod download

# Run the application
go run main.go
```

### Option 2: Using the Helper Script

```bash
# Make the script executable (first time only)
chmod +x run.sh

# Run the application
./run.sh
```

### Option 3: Using Make

```bash
# Run on default port (8080)
make run

# Run on alternative port (3000)
make run-alt
```

### Option 4: Using Docker

```bash
# Build and run with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the container
docker-compose down
```

### Option 5: Build Binary

```bash
# Build the binary
make build

# Run the binary
./bin/jiosaavn-api
```

## First Request

Once the server is running, test it with a health check:

```bash
curl http://localhost:8080/health
```

You should see:
```json
{
  "status": "ok",
  "message": "JioSaavn API is running"
}
```

## Try These Examples

### 1. Get Song Details

```bash
curl http://localhost:8080/song/IkuS3Pj6
```

### 2. Search for Songs

```bash
curl "http://localhost:8080/search/songs?q=tum%20hi%20ho"
```

### 3. Download a Song

```bash
curl http://localhost:8080/download/IkuS3Pj6
```

### 4. Get Lyrics

```bash
curl http://localhost:8080/lyrics/IkuS3Pj6
```

## Common Issues

### Port Already in Use

If port 8080 is already in use, run on a different port:

```bash
SERVER_PORT=3000 go run main.go
```

Or use the Makefile:

```bash
make run-alt
```

### Dependencies Not Found

Install Go dependencies:

```bash
go mod download
go mod tidy
```

### Permission Denied (run.sh)

Make the script executable:

```bash
chmod +x run.sh
```

## Configuration

Create a `.env` file (copy from `.env.example`):

```bash
cp .env.example .env
```

Edit the values as needed:

```env
SERVER_PORT=8080
JIOSAAVN_BASE_URL=https://www.jiosaavn.com/api.php
DECRYPTION_KEY=38346591
```

Then run:

```bash
# The app will automatically load .env (if using a package like godotenv)
# Or export them manually:
export $(cat .env | xargs)
go run main.go
```

## Testing

Run all tests:

```bash
make test
```

Run tests with coverage:

```bash
make test-coverage
```

## Next Steps

1. 📖 Read the full [README.md](README.md)
2. 🔍 Explore [API Examples](API_EXAMPLES.md)
3. 🤝 Check [Contributing Guidelines](CONTRIBUTING.md)
4. 🐛 Report issues on GitHub

## API Endpoints Overview

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/health` | GET | Health check |
| `/song/:id` | GET | Get song details |
| `/download/:id` | GET | Get download URL |
| `/lyrics/:id` | GET | Get song lyrics |
| `/search` | GET | Search all |
| `/search/songs` | GET | Search songs |
| `/search/albums` | GET | Search albums |
| `/search/artist` | GET | Search artists |
| `/search/playlist` | GET | Search playlists |
| `/artist/:id` | GET | Get artist details |
| `/album/:id` | GET | Get album details |

## Tips

1. **Use jq for pretty JSON output:**
   ```bash
   curl http://localhost:8080/song/IkuS3Pj6 | jq
   ```

2. **Test with HTTPie (simpler than curl):**
   ```bash
   http GET http://localhost:8080/song/IkuS3Pj6
   ```

3. **Watch logs in real-time:**
   ```bash
   go run main.go | tee app.log
   ```

4. **Enable production mode:**
   ```bash
   GIN_MODE=release go run main.go
   ```

## Support

- 📖 Documentation: Check all `.md` files
- 🐛 Issues: Report on GitHub
- 💬 Questions: Open a GitHub issue with `question` tag

Happy coding! 🚀
