# API Usage Examples

## Table of Contents
- [Health Check](#health-check)
- [Song Operations](#song-operations)
- [Search Operations](#search-operations)
- [Artist Operations](#artist-operations)
- [Album Operations](#album-operations)

## Health Check

Check if the API is running:

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "message": "JioSaavn API is running",
  "status": "ok"
}
```

## Song Operations

### Get Song Details

Get detailed information about a specific song:

```bash
curl http://localhost:8080/song/IkuS3Pj6
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "IkuS3Pj6",
      "song": "Tum Hi Ho",
      "album": "Aashiqui 2",
      "year": "2013",
      "duration": "262",
      "language": "hindi",
      "has_lyrics": "true",
      "primary_artists": "Arijit Singh",
      "singers": "Arijit Singh",
      "image": "https://c.saavncdn.com/191/Aashiqui-2-Hindi-2013-500x500.jpg",
      "media_url": "https://aac.saavncdn.com/191/xyz_320.mp4",
      "media_preview_url": "https://preview.saavncdn.com/191/xyz_96_p.mp4"
    }
  ]
}
```

### Download Song

Get the download URL for a song:

```bash
curl http://localhost:8080/download/IkuS3Pj6
```

**Response:**
```json
{
  "success": true,
  "download_url": "https://aac.saavncdn.com/191/xyz_320.mp4",
  "song_name": "Tum Hi Ho",
  "quality": "320kbps"
}
```

### Get Song Lyrics

Fetch lyrics for a song:

```bash
curl http://localhost:8080/lyrics/IkuS3Pj6
```

**Response:**
```json
{
  "success": true,
  "data": {
    "lyrics": "Baatein teri\nYaadein teri...",
    "snippet": "Baatein teri...",
    "copyright": "© 2013 T-Series"
  }
}
```

## Search Operations

### Search All

Search across all categories:

```bash
curl "http://localhost:8080/search?q=tum%20hi%20ho"
# Or with spaces (will be auto-encoded):
curl "http://localhost:8080/search?q=tum hi ho"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "songs": [...],
    "albums": [...],
    "artists": [...],
    "playlists": [...]
  }
}
```

### Search Songs Only

```bash
curl "http://localhost:8080/search/songs?q=tum%20hi%20ho"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "results": [
      {
        "id": "IkuS3Pj6",
        "title": "Tum Hi Ho",
        "album": "Aashiqui 2",
        "image": "https://c.saavncdn.com/191/Aashiqui-2-Hindi-2013-150x150.jpg",
        "url": "https://www.jiosaavn.com/song/..."
      }
    ]
  }
}
```

### Search Albums

```bash
curl "http://localhost:8080/search/albums?q=aashiqui"
```

### Search Artists

```bash
curl "http://localhost:8080/search/artist?q=arijit%20singh"
```

### Search Playlists

```bash
curl "http://localhost:8080/search/playlist?q=romantic"
```

## Artist Operations

### Get Artist Details

Get detailed information about an artist:

```bash
curl http://localhost:8080/artist/459320
```

**Response:**
```json
{
  "success": true,
  "data": {
    "artistId": "459320",
    "name": "Arijit Singh",
    "image": "https://c.saavncdn.com/artists/Arijit_Singh_500x500.jpg",
    "follower_count": "12345678",
    "type": "artist",
    "topSongs": [...],
    "topAlbums": [...],
    "bio": "..."
  }
}
```

## Album Operations

### Get Album Details

Get detailed information about an album:

```bash
curl http://localhost:8080/album/1134888
```

**Response:**
```json
{
  "success": true,
  "data": {
    "albumid": "1134888",
    "title": "Aashiqui 2",
    "image": "https://c.saavncdn.com/191/Aashiqui-2-Hindi-2013-500x500.jpg",
    "primary_artists": "Arijit Singh, Ankit Tiwari",
    "year": "2013",
    "language": "hindi",
    "songs": [
      {
        "id": "IkuS3Pj6",
        "song": "Tum Hi Ho",
        ...
      },
      ...
    ]
  }
}
```

## Using with Different Tools

### cURL with Pretty Print

```bash
curl -s http://localhost:8080/song/IkuS3Pj6 | jq
```

### HTTPie

```bash
http GET http://localhost:8080/song/IkuS3Pj6
```

### JavaScript (Fetch API)

```javascript
fetch('http://localhost:8080/search/songs?q=tum%20hi%20ho')
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error('Error:', error));
```

### Python (requests)

```python
import requests

response = requests.get('http://localhost:8080/song/IkuS3Pj6')
data = response.json()
print(data)
```

### Node.js (axios)

```javascript
const axios = require('axios');

axios.get('http://localhost:8080/song/IkuS3Pj6')
  .then(response => console.log(response.data))
  .catch(error => console.error('Error:', error));
```

## Error Responses

### Missing Parameter

```bash
curl http://localhost:8080/search/songs
```

**Response:**
```json
{
  "error": "Missing query parameter"
}
```

### Not Found

```bash
curl http://localhost:8080/song/invalid_id
```

**Response:**
```json
{
  "error": "Song not found"
}
```

### Server Error

```json
{
  "error": "Failed to fetch song"
}
```

## Rate Limiting Considerations

While this API doesn't implement rate limiting, be considerate when making requests to avoid overwhelming the upstream JioSaavn API:

- Implement client-side caching
- Use reasonable request intervals
- Handle errors gracefully with exponential backoff

## Tips

1. **URL Encoding**: Always URL-encode search queries with special characters
2. **ID Format**: Song/Album/Artist IDs are alphanumeric strings
3. **Image Quality**: Images can be resized by changing the URL (150x150, 500x500, etc.)
4. **Download Quality**: Songs are available in 320kbps (if available) or 160kbps
