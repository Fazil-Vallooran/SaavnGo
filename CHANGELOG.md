# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0] - 2024

### Added
- ✨ Complete project restructuring with proper separation of concerns
- 🔧 Configuration management system with environment variable support
- 🧪 Unit tests for utils and config packages
- 📦 Models package with proper data structures
- 🔀 Middleware package with CORS and Logger middleware
- 🐳 Docker and Docker Compose support
- 📝 Comprehensive documentation (README, API_EXAMPLES, etc.)
- 🛠️ Makefile for common development tasks
- 💚 Health check endpoint
- 📥 Download song handler with proper URL decryption
- 🎵 Lyrics endpoint
- 👤 Artist details endpoint
- 💿 Album details endpoint
- 🔍 All search endpoints (songs, albums, artists, playlists)
- 🎯 Better error handling throughout the application

### Changed
- ♻️ Refactored all inline route handlers to proper service functions
- 🔐 Moved hardcoded values to configuration
- 📊 Improved logging with custom middleware
- 🎨 Enhanced code organization and structure
- ✅ Better error handling in decrypt function

### Removed
- 🗑️ Duplicate route handlers
- 🗑️ Hardcoded configuration values

## [1.0.0] - Initial Release

### Added
- Basic API structure
- Song details endpoint
- Search functionality (basic)
- URL decryption utility
- Song formatting utility
