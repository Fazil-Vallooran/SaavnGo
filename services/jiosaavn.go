package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"jioSaavnAPI/config"
	"jioSaavnAPI/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var cfg = config.LoadConfig()

// GetSongHandler retrieves detailed information about a song
// @Summary      Get song details
// @Description  Returns detailed information about a song including artists, album, download URLs, and images
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Song ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /song/{id} [get]
func GetSongHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing song ID",
		})
		return
	}
	url := fmt.Sprintf("%s?__call=song.getDetails&cc=in&_format=json&_marker=0&pids=%s", cfg.JioSaavnBaseURL, id)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch song",
		})
		return
	}
	fmt.Println(resp.Body)
	defer resp.Body.Close()
	var raw map[string]map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		fmt.Print(raw)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to parse response",
		})
		return
	}

	songData, ok := raw[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Song not found",
		})
		return
	}

	formatted := utils.FormatSongDetailed(songData)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []any{formatted}})
}
// GetSongFromTokenHandler retrieves song information using a token
// @Summary      Get song details from token
// @Description  Returns detailed information about a song using a token
// @Tags         Songs
// @Accept       json
// @Produce      json
func GetSongFromTokenHandler(c *gin.Context) {
    token := c.Param("token")
    if token == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   "Missing token",
        })
        return
    }
    
    url := fmt.Sprintf("%s?__call=webapi.get&token=%s&type=song&includeMetaTags=0&ctx=web6dot0&api_version=4&_format=json&_marker=0", cfg.JioSaavnBaseURL, token)
    
    resp, err := http.Get(url)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to fetch song",
        })
        return
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   fmt.Sprintf("API returned status: %d", resp.StatusCode),
        })
        return
    }
    
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to read response body",
        })
        return
    }
    
    var response struct {
        Songs []map[string]interface{} `json:"songs"`
    }
    
    if err := json.Unmarshal(bodyBytes, &response); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "error":   "Failed to parse response: " + err.Error(),
        })
        return
    }
    
    if len(response.Songs) == 0 {
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "error":   "Song not found",
        })
        return
    }
    
    songData := response.Songs[0]
    
    // Use the new formatting function
    formatted := utils.FormatSongFromToken(songData)
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    []interface{}{formatted},
    })
}
// GetAlbumHandler retrieves detailed information about an album
// @Summary      Get album details
// @Description  Returns detailed information about an album including songs and artists
// @Tags         Albums
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Album ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /album/{id} [get]
func GetAlbumHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing album ID",
		})
		return
	}

	url := fmt.Sprintf("%s?__call=content.getAlbumDetails&_format=json&cc=in&_marker=0&albumid=%s", cfg.JioSaavnBaseURL, id)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch album",
		})
		return
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to parse response",
		})
		return
	}

	// Extract album data from the "data" key
	albumData, ok := raw["data"].(map[string]interface{})
	if !ok {
		if _, hasTitle := raw["title"]; hasTitle {
			albumData = raw
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Album not found",
			})
			return
		}
	}

	if len(albumData) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Album data is empty",
		})
		return
	}

	formatted := utils.FormatAlbum(albumData)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": formatted})
}

// GetArtistHandler retrieves detailed information about an artist
// @Summary      Get artist details
// @Description  Returns detailed information about an artist including bio, top songs, and albums
// @Tags         Artists
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Artist ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /artist/{id} [get]
func GetArtistHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing artist ID",
		})
		return
	}

	url := fmt.Sprintf("%s?__call=artist.getArtistPageDetails&_format=json&cc=in&_marker=0&artistId=%s", cfg.JioSaavnBaseURL, id)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch artist",
		})
		return
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to parse response",
		})
		return
	}

	// Check if artist data exists
	if len(raw) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Artist not found",
		})
		return
	}

	// Format the artist details
	formatted := utils.FormatArtistDetails(raw)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": formatted})
}

// GetLyricsHandler retrieves lyrics for a song
// @Summary      Get song lyrics
// @Description  Returns lyrics for a specific song with proper line breaks
// @Tags         Lyrics
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Lyrics ID (usually same as song ID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /lyrics/{id} [get]
func GetLyricsHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing song ID for Lyrics",
		})
		return
	}

	url := fmt.Sprintf("%s?__call=lyrics.getLyrics&lyrics_id=%s&ctx=web6dot0&api_version=4&_format=json&_marker=0",
		cfg.JioSaavnBaseURL, id)

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch lyrics",
		})
		return
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to parse response",
		})
		return
	}

	// Clean up the lyrics text
	if lyrics, ok := raw["lyrics"].(string); ok {
		// Replace HTML entities with actual line breaks
		lyrics = strings.ReplaceAll(lyrics, "\u003cbr\u003e", "\n")
		raw["lyrics"] = lyrics
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": raw})
}

// AutocompleteSongsHandler provides fast, lightweight song search results
// @Summary      Fast song autocomplete
// @Description  Lightweight song search optimized for quick results (returns only essential fields)
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        q      query     string  true   "Search query"
// @Param        limit  query     int     false  "Number of results (max 50)" default(10)
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /search/songs/autocomplete [get]
func AutocompleteHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing query parameter",
		})
		return
	}

	limit := 3
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 10 {
			limit = parsed
		}
	}

	// Use autocomplete endpoint for speed
	url := fmt.Sprintf("%s?__call=autocomplete.get&_format=json&_marker=0&query=%s&type=song",
		cfg.JioSaavnBaseURL, utils.EscapeString(query))

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch results",
		})
		return
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to parse response",
		})
		return
	}

	// Extract songs from the nested structure
	songs := []map[string]interface{}{}

	// First check topquery for best match
	if topQuery, ok := raw["topquery"].(map[string]interface{}); ok {
		if topData, ok := topQuery["data"].([]interface{}); ok && len(topData) > 0 {
			if topSong, ok := topData[0].(map[string]interface{}); ok {
				if topSong["type"] == "song" {
					songs = append(songs, formatLightweightSong(topSong))
				}
			}
		}
	}

	// Then get songs from songs section
	if songsSection, ok := raw["songs"].(map[string]interface{}); ok {
		if songsData, ok := songsSection["data"].([]interface{}); ok {
			for _, song := range songsData {
				if songMap, ok := song.(map[string]interface{}); ok {
					formatted := formatLightweightSong(songMap)
					// Avoid duplicates (check if same ID already in results)
					isDuplicate := false
					for _, existing := range songs {
						if existing["id"] == formatted["id"] {
							isDuplicate = true
							break
						}
					}
					if !isDuplicate {
						songs = append(songs, formatted)
					}
				}
			}
		}
	}

	// Limit results
	if len(songs) > limit {
		songs = songs[:limit]
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"total":   len(songs),
			"results": songs,
		},
	})
}

// formatLightweightSong formats song with only essential fields for quick search
func formatLightweightSong(data map[string]interface{}) map[string]interface{} {
	// Extract more_info if available
	moreInfo, _ := data["more_info"].(map[string]interface{})

	// Get basic info
	songID := utils.GetString(data, "id")
	title := utils.GetString(data, "title")
	album := utils.GetString(data, "album")

	// Get image - prefer larger size
	imageURL := utils.GetString(data, "image")
	imageURL = strings.Replace(imageURL, "50x50", "150x150", 1)

	// Get artists
	primaryArtists := utils.GetString(moreInfo, "primary_artists")
	if primaryArtists == "" {
		primaryArtists = utils.GetString(data, "primary_artists")
	}

	singers := utils.GetString(moreInfo, "singers")
	if singers == "" {
		singers = primaryArtists
	}

	// Get language
	language := utils.GetString(moreInfo, "language")
	if language == "" {
		language = utils.GetString(data, "language")
	}

	// Get song URL
	songURL := utils.GetString(data, "url")

	// Return lightweight response
	return map[string]interface{}{
		"id":          songID,
		"title":       strings.TrimSpace(title),
		"album":       strings.TrimSpace(album),
		"artists":     strings.TrimSpace(singers),
		"image":       imageURL,
		"url":         songURL,
		"language":    language,
		"description": fmt.Sprintf("%s · %s", strings.TrimSpace(singers), strings.TrimSpace(album)),
	}
}

// GetFullSearchResults uses search.getResults for paginated, comprehensive search
func GetFullSearchResults(query string, searchType string) (map[string]interface{}, error) {
	if query == "" {
		return nil, errors.New("missing query parameter")
	}

	// Determine which search endpoint to use based on type
	var callType string
	switch searchType {
	case "song":
		callType = "search.getResults"
	case "album":
		callType = "search.getAlbumResults"
	case "artist":
		callType = "search.getArtistResults"
	case "playlist":
		callType = "search.getPlaylistResults"
	default:
		callType = "search.getResults" // Default to songs
	}

	apiURL := fmt.Sprintf("%s?p=1&q=%s&_format=json&_marker=0&api_version=4&ctx=web6dot0&n=20&__call=%s",
		cfg.JioSaavnBaseURL, utils.EscapeString(query), callType)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch search results: %w", err)
	}
	defer resp.Body.Close()

	var raw map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract the data object which contains results
	if dataObj, ok := raw["data"].(map[string]interface{}); ok {
		return dataObj, nil
	}

	// If no data object, return raw results
	return raw, nil
}

// FullSearchHandler provides comprehensive paginated search results
// @Summary      Full search with pagination
// @Description  Comprehensive search results with pagination support for songs, albums, artists, and playlists
// @Tags         Search
// @Accept       json
// @Produce      json
// @Param        q      query     string  true   "Search query"
// @Param        type   query     string  false  "Search type: song, album, artist, playlist" default(song)
// @Param        page   query     int     false  "Page number" default(1)
// @Param        limit  query     int     false  "Results per page (max 50)" default(20)
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /search [get]
func FullSearchHandler(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Missing query parameter",
		})
		return
	}

	searchType := c.DefaultQuery("type", "song")

	// Get raw results from API
	results, err := GetFullSearchResults(query, searchType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Format results based on search type
	var formatted map[string]interface{}

	switch searchType {
	case "song":
		formatted = utils.FormatSongSearch(results)
	case "album":
		formatted = utils.FormatAlbumSearch(results)
	case "artist":
		formatted = utils.FormatArtistSearch(results)
	case "playlist":
		formatted = utils.FormatPlaylistSearch(results)
	default:
		formatted = utils.FormatSongSearch(results)
	}

	c.JSON(http.StatusOK, formatted)
}

// Helper to extract play_count safely
func getPlayCountFromSong(song map[string]interface{}) int {
	// Try direct field
	if playCount, ok := song["play_count"].(string); ok {
		if count, err := strconv.Atoi(playCount); err == nil {
			return count
		}
	}

	// Try nested in more_info
	if moreInfo, ok := song["more_info"].(map[string]interface{}); ok {
		if playCount, ok := moreInfo["play_count"].(string); ok {
			if count, err := strconv.Atoi(playCount); err == nil {
				return count
			}
		}
	}

	return 0
}
