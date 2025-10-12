package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// GetString safely extracts a string from a map
func GetString(data map[string]interface{}, key string) string {
	if data == nil {
		return ""
	}
	if val, ok := data[key]; ok && val != nil {
		return strings.TrimSpace(fmt.Sprintf("%v", val))
	}
	return ""
}

// GetInt safely extracts an integer from a map
func GetInt(data map[string]interface{}, key string) int {
	if data == nil {
		return 0
	}
	
	// Try to get as float64 first (JSON numbers are parsed as float64)
	if val, ok := data[key].(float64); ok {
		return int(val)
	}

	// Fallback to string parsing
	val := GetString(data, key)
	if val == "" {
		return 0
	}
	num, _ := strconv.Atoi(val)
	return num
}

// GetBool safely extracts a boolean from a map
func GetBool(data map[string]interface{}, key string) bool {
	if data == nil {
		return false
	}
	
	if val, ok := data[key].(bool); ok {
		return val
	}
	
	// Handle string "true"/"false"
	valStr := GetString(data, key)
	return valStr == "true" || valStr == "1"
}

// SanitizeImageURL converts image URLs to desired size
func SanitizeImageURL(url string, toSize string) string {
	if url == "" {
		return ""
	}
	
	sizes := []string{"50x50", "150x150", "500x500"}
	for _, size := range sizes {
		url = strings.Replace(url, size, toSize, -1)
	}
	return url
}
