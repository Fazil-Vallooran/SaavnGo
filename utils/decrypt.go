package utils

import (
	"crypto/des"
	"encoding/base64"
	"jioSaavnAPI/config"
	"strings"
)

var cfg = config.LoadConfig()

// DecryptURL decrypts the encrypted media URL from JioSaavn
func DecryptURL(encrypted string) string {

	if encrypted == "" {
		return ""
	}

	key := []byte(cfg.DecryptionKey)

	encData, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encrypted))
	if err != nil {
		return ""
	}

	// Validate data length is valid for DES
	if len(encData) == 0 || len(encData)%8 != 0 {
		return ""
	}

	block, err := des.NewCipher(key)
	if err != nil {
		return ""
	}

	decrypted := make([]byte, len(encData))
	for bs, be := 0, block.BlockSize(); bs < len(encData); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		if be > len(encData) {
			return "" // Invalid block size
		}
		block.Decrypt(decrypted[bs:be], encData[bs:be])
	}

	// Remove PKCS5 padding
	if len(decrypted) > 0 {
		padLen := int(decrypted[len(decrypted)-1])
		if padLen > 0 && padLen <= len(decrypted) {
			decrypted = decrypted[:len(decrypted)-padLen]
		}
	}

	url := string(decrypted)
	return strings.Replace(url, "_96.mp4", "_320.mp4", 1)
}
