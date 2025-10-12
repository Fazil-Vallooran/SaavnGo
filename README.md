# JioSaavn API

A REST API wrapper for JioSaavn built with Go and Gin framework.

## Features

- 🎵 Get song details by ID
- 🔍 Search for songs, albums, artists, and playlists
- 📥 Get download links for songs
- 📝 Fetch song lyrics
- 👤 Get artist information
- 💿 Get album details
- 🔐 Automatic URL decryption for media files

## Installation

### Prerequisites

- Go 1.16 or higher
- Git

### Setup

1. Clone the repository:
```bash
git clone <repository-url>
cd jioSaavnAPI
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run main.go
```

The server will start on `http://localhost:8080` by default.

## Configuration

You can configure the application using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | Port to run the server on | `8080` |
| `JIOSAAVN_BASE_URL` | JioSaavn API base URL | `https://www.jiosaavn.com/api.php` |
| `DECRYPTION_KEY` | Key for decrypting media URLs | `38346591` |

Example:
```bash
export SERVER_PORT=3000
go run main.go
```

## API Endpoints

### Health Check

```
GET /health
```

Returns the API health status.

### Song Details

```
GET /song/:id
```

Get detailed information about a song.

**Parameters:**
- `id` - Song ID

**Example:**
```bash
curl http://localhost:8080/song/abc123
```

### Search

#### Search All

```
GET /search?q=query
```

Search across all categories (songs, albums, artists, playlists).

#### Search Songs

```
GET /search/songs?q=query
```

Search for songs only.

#### Search Albums

```
GET /search/albums?q=query
```

Search for albums only.

#### Search Artists

```
GET /search/artist?q=query
```

Search for artists only.

#### Search Playlists

```
GET /search/playlist?q=query
```

Search for playlists only.

**Example:**
```bash
curl "http://localhost:8080/search/songs?q=tum%20hi%20ho"
```

### Download Song

```
GET /download/:id
```

Get the download URL for a song.

**Parameters:**
- `id` - Song ID

**Example:**
```bash
curl http://localhost:8080/download/abc123
```

### Lyrics

```
GET /lyrics/:id
```

Get lyrics for a song.

**Parameters:**
- `id` - Song ID

**Example:**
```bash
curl http://localhost:8080/lyrics/abc123
```

### Artist Details

```
GET /artist/:id
```

Get detailed information about an artist.

**Parameters:**
- `id` - Artist ID

**Example:**
```bash
curl http://localhost:8080/artist/abc123
```

### Album Details

```
GET /album/:id
```

Get detailed information about an album.

**Parameters:**
- `id` - Album ID

**Example:**
```bash
curl http://localhost:8080/album/abc123
```

## Project Structure

```
jioSaavnAPI/
├── config/          # Configuration management
├── middleware/      # Custom middleware (CORS, Logger)
├── models/          # Data models
├── routes/          # Route definitions
├── services/        # Business logic and API handlers
├── utils/           # Utility functions (encryption, formatting)
├── main.go          # Application entry point
└── README.md        # Documentation
```

## Response Format

### Success Response

```json
{
  "success": true,
  "data": {
    // Response data
  }
}
```

### Error Response

```json
{
  "error": "Error message"
}
```

## Building for Production

To build a production binary:

```bash
go build -o jiosaavn-api
./jiosaavn-api
```

For cross-platform builds:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o jiosaavn-api-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o jiosaavn-api.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o jiosaavn-api-mac
```

## Development

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
golangci-lint run
```

## License

This project is for educational purposes only. Please respect JioSaavn's terms of service.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Disclaimer

This API wrapper is not officially affiliated with JioSaavn. Use at your own risk and respect the original service's terms of use.
