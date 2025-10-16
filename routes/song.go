package routes

import (
	"github.com/gin-gonic/gin"
	"jioSaavnAPI/services"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(r *gin.Engine) {
	// Song routes - with trailing slash support
	r.GET("/song/:id", services.GetSongHandler)
	r.GET("/song/:id/", services.GetSongHandler)
	
	// Lyrics routes
	r.GET("/lyrics/:id", services.GetLyricsHandler)
	r.GET("/lyrics/:id/", services.GetLyricsHandler)
	
	// Album routes
	r.GET("/album/:id", services.GetAlbumHandler)
	r.GET("/album/:id/", services.GetAlbumHandler)
	
	// Artist routes
	r.GET("/artist/:id", services.GetArtistHandler)
	r.GET("/artist/:id/", services.GetArtistHandler)
	
	// Search routes
	r.GET("/search", services.FullSearchHandler)
	r.GET("/search/autocomplete", services.AutocompleteHandler)
}
