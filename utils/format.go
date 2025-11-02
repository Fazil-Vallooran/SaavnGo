package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// FormatSong formats the raw song data from JioSaavn API
func FormatSong(data map[string]interface{}) map[string]interface{} {
	// Safe check for encrypted_media_url
	encryptedURL := ""
	if val, ok := data["encrypted_media_url"]; ok && val != nil {
		encryptedURL = fmt.Sprintf("%v", val)
	}
	mediaURL := DecryptURL(encryptedURL)
	// Check if 320kbps is available
	is320kbps := false
	if val, ok := data["320kbps"]; ok && val != nil {
		is320kbps = fmt.Sprintf("%v", val) == "true"
	}

	if !is320kbps {
		mediaURL = strings.Replace(mediaURL, "_320.mp4", "_160.mp4", 1)
	}
	previewURL := strings.Replace(mediaURL, "_320.mp4", "_96_p.mp4", 1)
	previewURL = strings.Replace(previewURL, "_160.mp4", "_96_p.mp4", 1)
	previewURL = strings.Replace(previewURL, "//aac.", "//preview.", 1)
	data["media_url"] = mediaURL
	data["media_preview_url"] = previewURL
	for _, key := range []string{"song", "music", "singers", "starring", "album", "primary_artists"} {
		if val, ok := data[key]; ok {
			data[key] = strings.TrimSpace(fmt.Sprintf("%v", val))
		}
	}
	if img, ok := data["image"].(string); ok {
		data["image"] = strings.Replace(img, "150x150", "500x500", 1)
	}
	return data
}

// FormatSongFromToken formats song data from the webapi.get endpoint
// This endpoint has a different structure than other JioSaavn endpoints
func FormatSongFromToken(data map[string]interface{}) map[string]interface{} {
	// Extract more_info object
	moreInfo, _ := data["more_info"].(map[string]interface{})

	// Get encrypted media URL and decrypt it
	encryptedURL := GetString(moreInfo, "encrypted_media_url")
	mediaURL := DecryptURL(encryptedURL)

	// Determine max available quality
	has320 := GetString(moreInfo, "320kbps") == "true"
	maxQuality := "_160.mp4"
	if has320 {
		maxQuality = "_320.mp4"
	}

	// Build download URLs
	downloadURLs := []map[string]string{
		{"quality": "96kbps", "url": strings.Replace(mediaURL, maxQuality, "_96.mp4", 1)},
		{"quality": "160kbps", "url": strings.Replace(mediaURL, maxQuality, "_160.mp4", 1)},
	}

	if has320 {
		downloadURLs = append(downloadURLs, map[string]string{
			"quality": "320kbps",
			"url":     mediaURL,
		})
	}

	// Build image array with multiple sizes
	imageURL := GetString(data, "image")
	// Replace 150x150 with different sizes (API returns 150x150 by default)
	images := []map[string]string{
		{"quality": "50x50", "url": strings.Replace(imageURL, "150x150", "50x50", 1)},
		{"quality": "150x150", "url": imageURL},
		{"quality": "500x500", "url": strings.Replace(imageURL, "150x150", "500x500", 1)},
	}

	// Parse duration
	durationStr := GetString(moreInfo, "duration")
	duration := 0
	if durationStr != "" {
		duration, _ = strconv.Atoi(durationStr)
	}

	// Parse explicit content
	explicitContent := GetString(data, "explicit_content") == "1"

	// Parse play count
	playCountStr := GetString(data, "play_count")
	playCount := 0
	if playCountStr != "" {
		playCount, _ = strconv.Atoi(playCountStr)
	}

	// Parse has lyrics
	hasLyrics := GetString(moreInfo, "has_lyrics") == "true"

	// Build artists from artistMap
	artistMap, _ := moreInfo["artistMap"].(map[string]interface{})

	primaryArtists := buildArtistsFromMap(artistMap, "primary_artists")
	featuredArtists := buildArtistsFromMap(artistMap, "featured_artists")
	allArtists := buildArtistsFromMap(artistMap, "artists")

	// Ensure arrays are not nil
	if primaryArtists == nil {
		primaryArtists = []map[string]interface{}{}
	}
	if featuredArtists == nil {
		featuredArtists = []map[string]interface{}{}
	}
	if allArtists == nil {
		allArtists = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"id":              GetString(data, "id"),
		"name":            GetString(data, "title"),
		"type":            "song",
		"year":            GetString(data, "year"),
		"releaseDate":     GetString(moreInfo, "release_date"),
		"duration":        duration,
		"label":           GetString(moreInfo, "label"),
		"explicitContent": explicitContent,
		"playCount":       playCount,
		"language":        GetString(data, "language"),
		"hasLyrics":       hasLyrics,
		"lyricsId":        nil,
		"url":             GetString(data, "perma_url"),
		"copyright":       GetString(moreInfo, "copyright_text"),
		"album": map[string]interface{}{
			"id":   GetString(moreInfo, "album_id"),
			"name": GetString(moreInfo, "album"),
			"url":  GetString(moreInfo, "album_url"),
		},
		"artists": map[string]interface{}{
			"primary":  primaryArtists,
			"featured": featuredArtists,
			"all":      allArtists,
		},
		"image":       images,
		"downloadUrl": downloadURLs,
	}
}

// buildArtistsFromMap extracts artist data from the artistMap structure
func buildArtistsFromMap(artistMap map[string]interface{}, key string) []map[string]interface{} {
	if artistMap == nil {
		return []map[string]interface{}{}
	}

	artistsRaw, ok := artistMap[key]
	if !ok {
		return []map[string]interface{}{}
	}

	artistsArray, ok := artistsRaw.([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}

	result := make([]map[string]interface{}, 0, len(artistsArray))
	for _, artistRaw := range artistsArray {
		artist, ok := artistRaw.(map[string]interface{})
		if !ok {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":   GetString(artist, "id"),
			"name": GetString(artist, "name"),
			"role": GetString(artist, "role"),
			"image": []map[string]string{
				{
					"quality": "50x50",
					"url":     strings.Replace(GetString(artist, "image"), "150x150", "50x50", 1),
				},
				{
					"quality": "150x150",
					"url":     GetString(artist, "image"),
				},
				{
					"quality": "500x500",
					"url":     strings.Replace(GetString(artist, "image"), "150x150", "500x500", 1),
				},
			},
			"type": GetString(artist, "type"),
			"url":  GetString(artist, "perma_url"),
		})
	}

	return result
}

// FormatSongDetailed transforms raw JioSaavn API data into a clean, structured format
func FormatSongDetailed(data map[string]interface{}) map[string]interface{} {
	// First apply basic formatting
	data = FormatSong(data)

	// Now build the detailed structure
	baseURL := data["media_url"].(string)

	// Determine max available quality
	has320 := GetString(data, "320kbps") == "true"
	maxQuality := "_160.mp4"
	if has320 {
		maxQuality = "_320.mp4"
	}

	// Build download URLs - only include good quality options (96kbps and above)
	downloadURLs := []map[string]string{
		{"quality": "96kbps", "url": strings.Replace(baseURL, maxQuality, "_96.mp4", 1)},
		{"quality": "160kbps", "url": strings.Replace(baseURL, maxQuality, "_160.mp4", 1)},
	}

	// Only add 320kbps if actually available
	if has320 {
		downloadURLs = append(downloadURLs, map[string]string{
			"quality": "320kbps",
			"url":     baseURL,
		})
	}

	// Build image array with multiple sizes
	imageURL := GetString(data, "image")
	images := []map[string]string{
		{"quality": "50x50", "url": strings.Replace(imageURL, "500x500", "50x50", 1)},
		{"quality": "150x150", "url": strings.Replace(imageURL, "500x500", "150x150", 1)},
		{"quality": "500x500", "url": imageURL},
	}

	// Parse duration
	duration := GetInt(data, "duration")

	// Parse explicit content
	explicitContent := GetInt(data, "explicit_content") == 1

	// Parse play count
	playCount := GetInt(data, "play_count")

	// Parse has lyrics
	hasLyrics := GetString(data, "has_lyrics") == "true"

	// Build artists - avoid duplicates by using primary artists for "all" if no separate singers data
	primaryArtists := buildArtists(data, "primary_artists", "primary_artists_id", "primary_artists")
	featuredArtists := buildArtists(data, "featured_artists", "featured_artists_id", "featured_artists")

	// Use primary artists as "all" artists to avoid confusion
	allArtists := primaryArtists

	// Ensure featured artists is an array, not nil
	if len(featuredArtists) == 0 {
		featuredArtists = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"id":              GetString(data, "id"),
		"name":            GetString(data, "song"),
		"type":            "song",
		"year":            GetString(data, "year"),
		"releaseDate":     GetString(data, "release_date"),
		"duration":        duration,
		"label":           GetString(data, "label"),
		"explicitContent": explicitContent,
		"playCount":       playCount,
		"language":        GetString(data, "language"),
		"hasLyrics":       hasLyrics,
		"lyricsId":        nil,
		"url":             GetString(data, "perma_url"),
		"album": map[string]interface{}{
			"id":   GetString(data, "albumid"),
			"name": GetString(data, "album"),
			"url":  GetString(data, "album_url"),
		},
		"artists": map[string]interface{}{
			"primary":  primaryArtists,
			"featured": featuredArtists,
			"all":      allArtists,
		},
		"image":       images,
		"downloadUrl": downloadURLs,
	}
}

// BuildImageArray creates an array of image qualities from a base image URL
func BuildImageArray(imageURL string) []map[string]string {
	if imageURL == "" {
		return []map[string]string{}
	}
	imageURL = strings.Replace(imageURL, "150x150", "500x500", 1)
	return []map[string]string{
		{"quality": "50x50", "url": strings.Replace(imageURL, "500x500", "50x50", 1)},
		{"quality": "150x150", "url": strings.Replace(imageURL, "500x500", "150x150", 1)},
		{"quality": "500x500", "url": imageURL},
	}
}

type Song struct {
	ID string `json:"id"`
}

func FormatAlbumFromToken(listInterface interface{}) map[string]interface{} {
	// Initialize response structure with safe defaults
	result := map[string]interface{}{
		"id":       "",
		"name":     "",
		"type":     "album",
		"image":    "[]map[string]string{}",
		"songs":    []map[string]interface{}{},
	}

	// Type assert to slice
	listSlice, ok := listInterface.([]interface{})
	if !ok || len(listSlice) == 0 {
		return result
	}

	// Extract album metadata from first song
	// Token API embeds album info in each song, so we extract it once
	if firstSong, ok := listSlice[0].(map[string]interface{}); ok {
		// Extract album ID (may be in different locations)
		if albumID := GetString(firstSong, "id"); albumID != "" {
			result["id"] = albumID
		} 
		// Extract album name
		if albumName := GetString(firstSong, "title"); albumName != "" {
			albumName = strings.TrimSpace(albumName)
			result["title"] = strings.TrimSpace(albumName)
		}

		// Extract and format image URL
		if imageURL := GetString(firstSong, "image"); imageURL != "" {
			formattedImage := formatImageURL(imageURL)
			result["image"] = BuildImageArray(formattedImage)
		}
	}

	// Format all songs in the album
	songs := make([]map[string]interface{}, 0, len(listSlice))
	for _, item := range listSlice {
		songMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Use the existing formatAlbumSong function for consistency
		formattedSong := formatAlbumSongToken(songMap)
		songs = append(songs, formattedSong)
	}

	result["songs"] = songs
	result["songCount"] = len(songs)

	return result
}

// FormatPlaylistFromToken extracts playlist metadata and song list from a token response.
// Similar to albums, but playlists may have different metadata fields.
func FormatPlaylistFromToken(playlistInterface interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"id":         "",
		"title":      "",
		"subtitle":   "",
		"type":       "playlist",
		"image":      []map[string]string{},
		"language":   "",
		"url":        "",
		"songCount":  0,
		"description": "",
		"songs":      []map[string]interface{}{},
	}

	playlistMap, ok := playlistInterface.(map[string]interface{})
	if !ok {
		return result
	}

	// --- Extract top-level playlist metadata ---
	result["id"] = GetString(playlistMap, "id")
	result["title"] = GetString(playlistMap, "title")
	result["subtitle"] = GetString(playlistMap, "subtitle")
	result["language"] = GetString(playlistMap, "language")
	result["url"] = GetString(playlistMap, "perma_url")
	result["description"] = GetString(playlistMap, "header_desc")

	if imageURL := GetString(playlistMap, "image"); imageURL != "" {
		result["image"] = BuildImageArray(formatImageURL(imageURL))
	}

	// --- Extract songs from "list" field ---
	listRaw, ok := playlistMap["list"].([]interface{})
	if ok && len(listRaw) > 0 {
		songs := make([]map[string]interface{}, 0, len(listRaw))
		for _, songItem := range listRaw {
			if songMap, ok := songItem.(map[string]interface{}); ok {
				formatted := formatAlbumSong(songMap)
				songs = append(songs, formatted)
			}
		}
		result["songs"] = songs
		result["songCount"] = len(songs)
	}

	return result
}


// formatImageURL ensures the image URL uses the 500x500 resolution.
// Handles three common patterns in the API:
// 1. URLs with "150x150" that need replacement
// 2. URLs with "50x50" that need replacement  
// 3. URLs without resolution suffix that need "-500x500" appended
func formatImageURL(imageURL string) string {
	if imageURL == "" {
		return ""
	}

	// Already has 500x500, return as-is
	if strings.Contains(imageURL, "500x500") {
		return imageURL
	}

	// Replace common smaller sizes with 500x500
	if strings.Contains(imageURL, "150x150") {
		return strings.Replace(imageURL, "150x150", "500x500", 1)
	}

	if strings.Contains(imageURL, "50x50") {
		return strings.Replace(imageURL, "50x50", "500x500", 1)
	}

	// Append 500x500 before .jpg extension if no size specified
	if strings.HasSuffix(imageURL, ".jpg") {
		return strings.Replace(imageURL, ".jpg", "-500x500.jpg", 1)
	}

	return imageURL
}

func FormatPlaylistFromContents(contents interface{}) []Song {
	songs := []Song{}

	contentStr, ok := contents.(string)
	if !ok || contentStr == "" {
		return songs
	}

	ids := strings.Split(contentStr, ",")
	for _, id := range ids {
		songs = append(songs, Song{ID: id})
	}
	return songs
}

// FormatArtistDetails formats artist details response
func FormatArtistDetails(data map[string]interface{}) map[string]interface{} {
	imageURL := GetString(data, "image")
	images := BuildImageArray(imageURL)

	followerCount := GetInt(data, "follower_count")

	// Extract top songs if available
	topSongs := []map[string]interface{}{}
	if songs, ok := data["topSongs"].([]interface{}); ok {
		for _, song := range songs {
			if songMap, ok := song.(map[string]interface{}); ok {
				topSongs = append(topSongs, FormatSongDetailed(songMap))
			}
		}
	}

	// Extract top albums if available
	topAlbums := []map[string]interface{}{}
	if albums, ok := data["topAlbums"].([]interface{}); ok {
		for _, album := range albums {
			if albumMap, ok := album.(map[string]interface{}); ok {
				topAlbums = append(topAlbums, map[string]interface{}{
					"id":    GetString(albumMap, "albumid"),
					"name":  GetString(albumMap, "title"),
					"image": BuildImageArray(GetString(albumMap, "image")),
					"url":   GetString(albumMap, "perma_url"),
					"year":  GetString(albumMap, "year"),
				})
			}
		}
	}

	return map[string]interface{}{
		"id":               GetString(data, "artistId"),
		"name":             GetString(data, "name"),
		"url":              GetString(data, "perma_url"),
		"type":             "artist",
		"followerCount":    followerCount,
		"fanCount":         GetString(data, "fan_count"),
		"isVerified":       GetString(data, "isVerified") == "true",
		"dominantLanguage": GetString(data, "dominantLanguage"),
		"dominantType":     GetString(data, "dominantType"),
		"bio":              GetString(data, "bio"),
		"dob":              GetString(data, "dob"),
		"fb":               GetString(data, "fb"),
		"twitter":          GetString(data, "twitter"),
		"wiki":             GetString(data, "wiki"),
		"image":            images,
		"topSongs":         topSongs,
		"topAlbums":        topAlbums,
	}
}

// Helper function to build artist arrays
func buildArtists(data map[string]interface{}, nameKey, idKey, role string) []map[string]interface{} {
	names := GetString(data, nameKey)
	ids := GetString(data, idKey)

	if names == "" || ids == "" {
		return []map[string]interface{}{}
	}

	nameList := strings.Split(names, ", ")
	idList := strings.Split(ids, ", ")

	artists := []map[string]interface{}{}
	for i, name := range nameList {
		name = strings.TrimSpace(name)
		id := ""
		if i < len(idList) {
			id = strings.TrimSpace(idList[i])
		}

		artists = append(artists, map[string]interface{}{
			"id":    id,
			"name":  name,
			"role":  role,
			"image": []string{},
			"type":  "artist",
		})
	}

	return artists
}

// FormatSearchSong formats search result songs to match the detailed song format
func FormatSearchSong(data map[string]interface{}) map[string]interface{} {
	// Get more_info nested object
	moreInfo, _ := data["more_info"].(map[string]interface{})

	// Decrypt and build download URLs
	encryptedURL := GetString(moreInfo, "encrypted_media_url")
	if encryptedURL == "" {
		encryptedURL = GetString(data, "encrypted_media_url")
	}

	baseURL := DecryptURL(encryptedURL)
	has320 := GetString(moreInfo, "320kbps") == "true"

	downloadURLs := []map[string]string{
		{"quality": "96kbps", "url": strings.Replace(baseURL, "_320.mp4", "_96.mp4", 1)},
		{"quality": "160kbps", "url": strings.Replace(baseURL, "_320.mp4", "_160.mp4", 1)},
	}
	if has320 {
		downloadURLs = append(downloadURLs, map[string]string{
			"quality": "320kbps",
			"url":     baseURL,
		})
	}

	// Build image array
	imageURL := GetString(data, "image")
	imageURL = strings.Replace(imageURL, "150x150", "500x500", 1)
	images := BuildImageArray(imageURL)

	// Parse duration
	duration := GetInt(moreInfo, "duration")
	if duration == 0 {
		duration = GetInt(data, "duration")
	}

	// Parse explicit content
	explicitContent := GetString(data, "explicit_content") == "1"
	if !explicitContent {
		explicitContent = GetString(moreInfo, "explicit_content") == "1"
	}

	// Parse play count
	playCount := GetInt(data, "play_count")

	// Parse has lyrics
	hasLyrics := GetString(moreInfo, "has_lyrics") == "true"

	// Build artists from artistMap in more_info
	primaryArtists := []map[string]interface{}{}
	featuredArtists := []map[string]interface{}{}
	allArtists := []map[string]interface{}{}

	if artistMap, ok := moreInfo["artistMap"].(map[string]interface{}); ok {
		// Primary artists
		if primary, ok := artistMap["primary_artists"].([]interface{}); ok {
			for _, artist := range primary {
				if a, ok := artist.(map[string]interface{}); ok {
					artistImage := GetString(a, "image")
					primaryArtists = append(primaryArtists, map[string]interface{}{
						"id":    GetString(a, "id"),
						"name":  GetString(a, "name"),
						"role":  GetString(a, "role"),
						"type":  "artist",
						"image": BuildImageArray(artistImage),
						"url":   GetString(a, "perma_url"),
					})
				}
			}
		}

		// Featured artists
		if featured, ok := artistMap["featured_artists"].([]interface{}); ok {
			for _, artist := range featured {
				if a, ok := artist.(map[string]interface{}); ok {
					artistImage := GetString(a, "image")
					featuredArtists = append(featuredArtists, map[string]interface{}{
						"id":    GetString(a, "id"),
						"name":  GetString(a, "name"),
						"role":  GetString(a, "role"),
						"type":  "artist",
						"image": BuildImageArray(artistImage),
						"url":   GetString(a, "perma_url"),
					})
				}
			}
		}

		// All artists
		if artists, ok := artistMap["artists"].([]interface{}); ok {
			for _, artist := range artists {
				if a, ok := artist.(map[string]interface{}); ok {
					artistImage := GetString(a, "image")
					allArtists = append(allArtists, map[string]interface{}{
						"id":    GetString(a, "id"),
						"name":  GetString(a, "name"),
						"role":  GetString(a, "role"),
						"type":  "artist",
						"image": BuildImageArray(artistImage),
						"url":   GetString(a, "perma_url"),
					})
				}
			}
		}
	}

	// Ensure arrays are not nil
	if len(featuredArtists) == 0 {
		featuredArtists = []map[string]interface{}{}
	}
	if len(primaryArtists) == 0 {
		primaryArtists = []map[string]interface{}{}
	}
	if len(allArtists) == 0 {
		allArtists = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"id":              GetString(data, "id"),
		"name":            strings.TrimSpace(GetString(data, "title")),
		"type":            "song",
		"year":            GetString(data, "year"),
		"releaseDate":     GetString(moreInfo, "release_date"),
		"duration":        duration,
		"label":           GetString(moreInfo, "label"),
		"explicitContent": explicitContent,
		"playCount":       playCount,
		"language":        GetString(data, "language"),
		"hasLyrics":       hasLyrics,
		"lyricsId":        nil,
		"url":             GetString(data, "perma_url"),
		"copyright":       GetString(moreInfo, "copyright_text"),
		"album": map[string]interface{}{
			"id":   GetString(moreInfo, "album_id"),
			"name": GetString(moreInfo, "album"),
			"url":  GetString(moreInfo, "album_url"),
		},
		"artists": map[string]interface{}{
			"primary":  primaryArtists,
			"featured": featuredArtists,
			"all":      allArtists,
		},
		"image":       images,
		"downloadUrl": downloadURLs,
	}
}

// FormatSearchArtist formats search result artists
func FormatSearchArtist(data map[string]interface{}) map[string]interface{} {
	// Build image array
	imageURL := GetString(data, "image")
	images := BuildImageArray(imageURL)

	return map[string]interface{}{
		"id":    GetString(data, "id"),
		"name":  strings.TrimSpace(GetString(data, "name")),
		"role":  GetString(data, "role"),
		"type":  "artist",
		"image": images,
		"url":   GetString(data, "perma_url"),
	}
}

// FormatSearchPlaylist formats search result playlists
func FormatSearchPlaylist(data map[string]interface{}) map[string]interface{} {
	// Get more_info nested object
	moreInfo, _ := data["more_info"].(map[string]interface{})

	// Build image array
	imageURL := GetString(data, "image")
	// Handle error cases where image might be HTML error message
	if strings.HasPrefix(imageURL, "<!doctype") || strings.HasPrefix(imageURL, "<html") {
		imageURL = "" // Use empty string for error cases
	}
	images := BuildImageArray(imageURL)

	// Parse song count
	songCount := 0
	if songCountStr := GetString(moreInfo, "song_count"); songCountStr != "" {
		songCount = GetInt(moreInfo, "song_count")
		// If GetInt failed, try parsing the string
		if songCount == 0 {
			parsed, _ := strconv.Atoi(songCountStr)
			songCount = parsed
		}
	}

	// Parse explicit content
	explicitContent := GetString(data, "explicit_content") == "1"

	// Get language
	language := GetString(moreInfo, "language")
	if language == "" {
		language = GetString(data, "language")
	}

	return map[string]interface{}{
		"id":              GetString(data, "id"),
		"name":            strings.TrimSpace(GetString(data, "title")),
		"type":            "playlist",
		"image":           images,
		"url":             GetString(data, "perma_url"),
		"songCount":       songCount,
		"language":        language,
		"explicitContent": explicitContent,
	}
}

// FormatAlbumDetailed transforms raw JioSaavn album data into a clean, structured format
func FormatAlbumDetailed(data map[string]interface{}) map[string]interface{} {
	// Build image array
	imageURL := GetString(data, "image")
	// Check if image is already in the correct format or needs conversion
	if !strings.Contains(imageURL, "500x500") && !strings.Contains(imageURL, "150x150") {
		imageURL = strings.Replace(imageURL, ".jpg", "-500x500.jpg", 1)
	} else {
		imageURL = strings.Replace(imageURL, "150x150", "500x500", 1)
	}
	images := BuildImageArray(imageURL)

	// Parse play count
	playCount := 0
	if pc := GetString(data, "play_count"); pc != "" {
		playCount = GetInt(data, "play_count")
	}

	// Parse explicit content
	explicitContent := GetString(data, "explicit_content") == "1"

	// Extract artists from more_info if available
	primaryArtists := []map[string]interface{}{}
	featuredArtists := []map[string]interface{}{}
	allArtists := []map[string]interface{}{}

	if moreInfo, ok := data["more_info"].(map[string]interface{}); ok {
		if artistMap, ok := moreInfo["artistMap"].(map[string]interface{}); ok {
			// Primary artists
			if primary, ok := artistMap["primary_artists"].([]interface{}); ok {
				for _, artist := range primary {
					if a, ok := artist.(map[string]interface{}); ok {
						primaryArtists = append(primaryArtists, map[string]interface{}{
							"id":    GetString(a, "id"),
							"name":  GetString(a, "name"),
							"role":  GetString(a, "role"),
							"type":  "artist",
							"image": []string{},
						})
					}
				}
			}

			// Featured artists
			if featured, ok := artistMap["featured_artists"].([]interface{}); ok {
				for _, artist := range featured {
					if a, ok := artist.(map[string]interface{}); ok {
						featuredArtists = append(featuredArtists, map[string]interface{}{
							"id":    GetString(a, "id"),
							"name":  GetString(a, "name"),
							"role":  GetString(a, "role"),
							"type":  "artist",
							"image": []string{},
						})
					}
				}
			}

			// All artists
			if artists, ok := artistMap["artists"].([]interface{}); ok {
				for _, artist := range artists {
					if a, ok := artist.(map[string]interface{}); ok {
						allArtists = append(allArtists, map[string]interface{}{
							"id":    GetString(a, "id"),
							"name":  GetString(a, "name"),
							"role":  GetString(a, "role"),
							"type":  "artist",
							"image": []string{},
						})
					}
				}
			}
		}
	}

	// Ensure arrays are not nil
	if len(featuredArtists) == 0 {
		featuredArtists = []map[string]interface{}{}
	}
	if len(primaryArtists) == 0 {
		primaryArtists = []map[string]interface{}{}
	}
	if len(allArtists) == 0 {
		allArtists = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"id":              GetString(data, "id"),
		"name":            strings.TrimSpace(GetString(data, "title")),
		"description":     GetString(data, "description"),
		"type":            "album",
		"year":            GetString(data, "year"),
		"playCount":       playCount,
		"language":        GetString(data, "language"),
		"explicitContent": explicitContent,
		"url":             GetString(data, "perma_url"),
		"image":           images,
		"artists": map[string]interface{}{
			"primary":  primaryArtists,
			"featured": featuredArtists,
			"all":      allArtists,
		},
	}
}

func FormatAlbum(data map[string]interface{}) map[string]interface{} {
	// Build image array
	imageURL := GetString(data, "image")
	if !strings.Contains(imageURL, "500x500") {
		imageURL = strings.Replace(imageURL, ".jpg", "-500x500.jpg", 1)
	} else {
		imageURL = strings.Replace(imageURL, "150x150", "500x500", 1)
	}
	images := BuildImageArray(imageURL)

	// Parse songs - Songs from album endpoint have flat structure
	songs := []map[string]interface{}{}
	if rawSongs, ok := data["songs"].([]interface{}); ok {
		for _, s := range rawSongs {
			if songMap, ok := s.(map[string]interface{}); ok {
				// Songs from album endpoint need special handling
				// They have artistMap instead of separate artist fields
				formattedSong := formatAlbumSong(songMap)
				songs = append(songs, formattedSong)
			}
		}
	}

	// Parse explicit content
	explicitContent := GetString(data, "explicit_content") == "1"

	// Extract artists from primary_artists fields
	primaryArtists := []map[string]interface{}{}

	primaryArtistNames := GetString(data, "primary_artists")
	primaryArtistIds := GetString(data, "primary_artists_id")

	if primaryArtistNames != "" && primaryArtistIds != "" {
		nameList := strings.Split(primaryArtistNames, ", ")
		idList := strings.Split(primaryArtistIds, ", ")
		for i, name := range nameList {
			id := ""
			if i < len(idList) {
				id = strings.TrimSpace(idList[i])
			}
			artist := map[string]interface{}{
				"id":    id,
				"name":  strings.TrimSpace(name),
				"role":  "primary_artists",
				"type":  "artist",
				"image": []map[string]string{},
			}
			primaryArtists = append(primaryArtists, artist)
		}
	}

	return map[string]interface{}{
		"id":              GetString(data, "albumid"),
		"name":            strings.TrimSpace(GetString(data, "name")),
		"type":            "album",
		"year":            GetString(data, "year"),
		"language":        GetString(data, "language"),
		"explicitContent": explicitContent,
		"songCount":       len(songs),
		"url":             GetString(data, "perma_url"),
		"image":           images,
		"artists": map[string]interface{}{
			"primary": primaryArtists,
		},
		"songs": songs,
	}
}

// formatAlbumSong formats songs from album endpoint which have a different structure
func formatAlbumSong(data map[string]interface{}) map[string]interface{} {
	// Decrypt and build download URLs
	encryptedURL := GetString(data, "encrypted_media_url")
	baseURL := DecryptURL(encryptedURL)
	has320 := GetString(data, "320kbps") == "true"

	maxQuality := "_160.mp4"
	if has320 {
		maxQuality = "_320.mp4"
		baseURL = strings.Replace(baseURL, "_160.mp4", "_320.mp4", 1)
	}

	downloadURLs := []map[string]string{
		{"quality": "96kbps", "url": strings.Replace(baseURL, maxQuality, "_96.mp4", 1)},
		{"quality": "160kbps", "url": strings.Replace(baseURL, maxQuality, "_160.mp4", 1)},
	}
	if has320 {
		downloadURLs = append(downloadURLs, map[string]string{
			"quality": "320kbps",
			"url":     baseURL,
		})
	}

	// Build image array
	imageURL := GetString(data, "image")
	imageURL = strings.Replace(imageURL, "150x150", "500x500", 1)
	images := BuildImageArray(imageURL)

	// Parse duration
	duration := GetInt(data, "duration")

	// Parse explicit content
	explicitContent := GetInt(data, "explicit_content") == 1

	// Parse play count
	playCount := GetInt(data, "play_count")

	// Parse has lyrics
	hasLyrics := GetString(data, "has_lyrics") == "true"

	return map[string]interface{}{
		"id":              GetString(data, "id"),
		"name":            strings.TrimSpace(GetString(data, "song")),
		"type":            "song",
		"year":            GetString(data, "year"),
		"releaseDate":     GetString(data, "release_date"),
		"duration":        duration,
		"label":           GetString(data, "label"),
		"explicitContent": explicitContent,
		"playCount":       playCount,
		"language":        GetString(data, "language"),
		"hasLyrics":       hasLyrics,
		"lyricsId":        nil,
		"url":             GetString(data, "perma_url"),
		"album": map[string]interface{}{
			"id":   GetString(data, "albumid"),
			"name": GetString(data, "album"),
			"url":  GetString(data, "album_url"),
		},
		"image":       images,
		"downloadUrl": downloadURLs,
	}
}

// formatAlbumSong formats songs from album endpoint which have a different structure
func formatAlbumSongToken(data map[string]interface{}) map[string]interface{} {

	return map[string]interface{}{
		"id":              GetString(data, "id"),
		"name":            strings.TrimSpace(GetString(data, "title")),
		"type":            "song",
		
	}
}
// EscapeString URL-encodes special characters in a string
func EscapeString(str string) string {
	replacer := strings.NewReplacer(
		" ", "%20", "!", "%21", "\"", "%22", "#", "%23", "$", "%24",
		"%", "%25", "&", "%26", "'", "%27", "(", "%28", ")", "%29",
		"*", "%2A", "+", "%2B", ",", "%2C", "-", "%2D", ".", "%2E", "/", "%2F",
		":", "%3A", ";", "%3B", "<", "%3C", "=", "%3D", ">", "%3E",
		"?", "%3F", "@", "%40", "[", "%5B", "\\", "%5C", "]", "%5D", "^", "%5E",
		"_", "%5F", "`", "%60", "{", "%7B", "|", "%7C", "}", "%7D", "~", "%7E",
	)
	return replacer.Replace(str)
}

func formatSongSearchItem(data map[string]interface{}) map[string]interface{} {
	// Extract nested fields safely
	album := map[string]interface{}{}
	if rawAlbum, ok := data["album"].(map[string]interface{}); ok {
		album = map[string]interface{}{
			"id":   GetString(rawAlbum, "id"),
			"name": GetString(rawAlbum, "name"),
			"url":  GetString(rawAlbum, "url"),
		}
	}

	// Build primary artists
	primaryArtists := []map[string]interface{}{}
	if rawArtists, ok := data["artists"].(map[string]interface{}); ok {
		if primary, ok := rawArtists["primary"].([]interface{}); ok {
			for _, a := range primary {
				if artist, ok := a.(map[string]interface{}); ok {
					primaryArtists = append(primaryArtists, map[string]interface{}{
						"id":    GetString(artist, "id"),
						"name":  GetString(artist, "name"),
						"role":  GetString(artist, "role"),
						"type":  "artist",
						"image": artist["image"],
						"url":   GetString(artist, "url"),
					})
				}
			}
		}
	}

	// Build image array
	images := []map[string]string{}
	if imageList, ok := data["image"].([]interface{}); ok {
		for _, img := range imageList {
			if imgMap, ok := img.(map[string]interface{}); ok {
				images = append(images, map[string]string{
					"quality": GetString(imgMap, "quality"),
					"url":     GetString(imgMap, "url"),
				})
			}
		}
	}

	// Build download URLs
	downloadUrls := []map[string]string{}
	if dls, ok := data["downloadUrl"].([]interface{}); ok {
		for _, d := range dls {
			if dlMap, ok := d.(map[string]interface{}); ok {
				downloadUrls = append(downloadUrls, map[string]string{
					"quality": GetString(dlMap, "quality"),
					"url":     GetString(dlMap, "url"),
				})
			}
		}
	}

	return map[string]interface{}{
		"id":              GetString(data, "id"),
		"name":            GetString(data, "name"),
		"type":            GetString(data, "type"),
		"year":            GetString(data, "year"),
		"releaseDate":     GetString(data, "releaseDate"),
		"duration":        GetInt(data, "duration"),
		"label":           GetString(data, "label"),
		"explicitContent": data["explicitContent"] == true,
		"playCount":       GetInt(data, "playCount"),
		"language":        GetString(data, "language"),
		"hasLyrics":       data["hasLyrics"] == true,
		"lyricsId":        nil,
		"url":             GetString(data, "url"),
		"album":           album,
		"artists": map[string]interface{}{
			"primary": primaryArtists,
		},
		"image":       images,
		"downloadUrl": downloadUrls,
	}
}

// FormatSongSearch formats search response containing multiple songs
func FormatSongSearch(data map[string]interface{}) map[string]interface{} {
	resultsData, _ := data["results"].([]interface{})
	start := GetInt(data, "start")
	total := GetInt(data, "total")

	var formattedResults []map[string]interface{}

	for _, item := range resultsData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		// Get more_info nested object
		moreInfo, _ := itemMap["more_info"].(map[string]interface{})

		// Build image array
		imageURL := GetString(itemMap, "image")
		if imageURL == "" {
			imageURL = GetString(moreInfo, "image")
		}
		imageURL = strings.Replace(imageURL, "150x150", "500x500", 1)
		images := BuildImageArray(imageURL)

		// Parse duration
		duration := GetInt(moreInfo, "duration")
		if duration == 0 {
			duration = GetInt(itemMap, "duration")
		}

		// Parse explicit content
		explicitContent := GetString(itemMap, "explicit_content") == "1"

		// Parse play count
		playCount := GetInt(itemMap, "play_count")

		// Build primary artists
		primaryArtists := []map[string]interface{}{}
		if artistMap, ok := moreInfo["artistMap"].(map[string]interface{}); ok {
			if primary, ok := artistMap["primary_artists"].([]interface{}); ok {
				for _, artist := range primary {
					if a, ok := artist.(map[string]interface{}); ok {
						primaryArtists = append(primaryArtists, map[string]interface{}{
							"id":   GetString(a, "id"),
							"name": GetString(a, "name"),
							"role": GetString(a, "role"),
							"type": "artist",
							"url":  GetString(a, "perma_url"),
						})
					}
				}
			}
		}

		// Build album info
		album := map[string]interface{}{
			"id":   GetString(moreInfo, "album_id"),
			"name": GetString(moreInfo, "album"),
			"url":  GetString(moreInfo, "album_url"),
		}

		formatted := map[string]interface{}{
			"id":              GetString(itemMap, "id"),
			"name":            strings.TrimSpace(GetString(itemMap, "title")),
			"type":            "song",
			"year":            GetString(itemMap, "year"),
			"language":        GetString(itemMap, "language"),
			"explicitContent": explicitContent,
			"duration":        duration,
			"playCount":       playCount,
			"url":             GetString(itemMap, "perma_url"),
			"image":           images,
			"album":           album,
			"artists": map[string]interface{}{
				"primary": primaryArtists,
			},
		}

		formattedResults = append(formattedResults, formatted)
	}
	// Sort by playCount (descending - most played first)
	sort.Slice(formattedResults, func(i, j int) bool {
		playCountI, _ := formattedResults[i]["playCount"].(int)
		playCountJ, _ := formattedResults[j]["playCount"].(int)
		return playCountI > playCountJ
	})

	return map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"total":   total,
			"start":   start,
			"results": formattedResults,
		},
	}
}

func FormatArtistSearch(data interface{}) map[string]interface{} {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"total":   0,
				"start":   0,
				"results": []map[string]interface{}{},
			},
		}
	}

	resultsData, _ := dataMap["results"].([]interface{})
	start := GetInt(dataMap, "start")
	total := GetInt(dataMap, "total")

	var formattedResults []map[string]interface{}

	for _, item := range resultsData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		formatted := map[string]interface{}{
			"id":              GetString(itemMap, "id"),
			"name":            GetString(itemMap, "name"),
			"description":     GetString(itemMap, "description"),
			"type":            GetString(itemMap, "type"),
			"url":             GetString(itemMap, "perma_url"),
			"image":           getImageArray(itemMap["image"]),
			"explicitContent": itemMap["explicitContent"],
		}

		// Handle followers/follower_count
		if followers, ok := itemMap["followers"]; ok {
			formatted["followers"] = followers
		} else if followerCount, ok := itemMap["follower_count"]; ok {
			formatted["followers"] = followerCount
		}

		formattedResults = append(formattedResults, formatted)
	}

	return map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"total":   total,
			"start":   start,
			"results": formattedResults,
		},
	}
}

func FormatPlaylistSearch(data interface{}) map[string]interface{} {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"total":   0,
				"start":   0,
				"results": []map[string]interface{}{},
			},
		}
	}

	resultsData, _ := dataMap["results"].([]interface{})
	start := GetInt(dataMap, "start")
	total := GetInt(dataMap, "total")

	var formattedResults []map[string]interface{}

	for _, item := range resultsData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		formatted := map[string]interface{}{
			"id":              GetString(itemMap, "id"),
			"name":            GetString(itemMap, "title"),
			"description":     GetString(itemMap, "description"),
			"type":            GetString(itemMap, "type"),
			"url":             GetString(itemMap, "perma_url"),
			"image":           getImageArray(itemMap["image"]),
			"explicitContent": itemMap["explicitContent"],
			"language":        GetString(itemMap, "language"),
		}

		// Handle song count
		if songCount, ok := itemMap["songCount"]; ok {
			formatted["songCount"] = songCount
		} else if songCountStr := GetString(itemMap, "song_count"); songCountStr != "" {
			if count, err := strconv.Atoi(songCountStr); err == nil {
				formatted["songCount"] = count
			}
		}

		// Extract artists if available
		if artistsMap, ok := itemMap["artists"].(map[string]interface{}); ok {
			formatted["artists"] = map[string]interface{}{
				"primary":  extractArtists(artistsMap["primary"]),
				"featured": extractArtists(artistsMap["featured"]),
				"all":      extractArtists(artistsMap["all"]),
			}
		}

		formattedResults = append(formattedResults, formatted)
	}

	return map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"total":   total,
			"start":   start,
			"results": formattedResults,
		},
	}
}

func FormatAlbumSearch(data interface{}) map[string]interface{} {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"total":   0,
				"start":   0,
				"results": []map[string]interface{}{},
			},
		}
	}

	resultsData, _ := dataMap["results"].([]interface{})
	start := GetInt(dataMap, "start")
	total := GetInt(dataMap, "total")

	var formattedResults []map[string]interface{}

	for _, item := range resultsData {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		formatted := map[string]interface{}{
			"id":              GetString(itemMap, "id"),
			"name":            GetString(itemMap, "title"),
			"description":     GetString(itemMap, "description"),
			"type":            GetString(itemMap, "type"),
			"year":            GetString(itemMap, "year"),
			"language":        GetString(itemMap, "language"),
			"url":             GetString(itemMap, "perma_url"),
			"image":           getImageArray(itemMap["image"]),
			"explicitContent": itemMap["explicitContent"],
		}

		// Handle play count
		if playCount, ok := itemMap["playCount"]; ok {
			formatted["playCount"] = playCount
		} else if pc := GetInt(itemMap, "play_count"); pc > 0 {
			formatted["playCount"] = pc
		}

		// Handle song count
		if songCount, ok := itemMap["songCount"]; ok {
			formatted["songCount"] = songCount
		} else if sc := GetInt(itemMap, "song_count"); sc > 0 {
			formatted["songCount"] = sc
		}

		// Extract artists if available
		if artistsMap, ok := itemMap["artists"].(map[string]interface{}); ok {
			formatted["artists"] = map[string]interface{}{
				"primary":  extractArtists(artistsMap["primary"]),
				"featured": extractArtists(artistsMap["featured"]),
				"all":      extractArtists(artistsMap["all"]),
			}
		}

		formattedResults = append(formattedResults, formatted)
	}

	return map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"total":   total,
			"start":   start,
			"results": formattedResults,
		},
	}
}

// Extract and format artist array
func extractArtists(data interface{}) []map[string]interface{} {
	artistList, ok := data.([]interface{})
	if !ok {
		return []map[string]interface{}{}
	}

	var artists []map[string]interface{}
	for _, a := range artistList {
		artist, ok := a.(map[string]interface{})
		if !ok {
			continue
		}

		artistObj := map[string]interface{}{
			"id":   GetString(artist, "id"),
			"name": GetString(artist, "name"),
			"role": GetString(artist, "role"),
			"type": "artist",
			"url":  GetString(artist, "perma_url"),
		}

		// Handle image properly
		if img := artist["image"]; img != nil {
			artistObj["image"] = getImageArray(img)
		} else {
			artistObj["image"] = []map[string]string{}
		}

		artists = append(artists, artistObj)
	}
	return artists
}

// Handle image arrays safely
func getImageArray(data interface{}) []map[string]string {
	imageList, ok := data.([]interface{})
	if !ok {
		return []map[string]string{}
	}

	var images []map[string]string
	for _, img := range imageList {
		if imgMap, ok := img.(map[string]interface{}); ok {
			images = append(images, map[string]string{
				"quality": GetString(imgMap, "quality"),
				"url":     GetString(imgMap, "url"),
			})
		}
	}
	return images
}
