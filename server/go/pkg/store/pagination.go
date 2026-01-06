package store

import (
	"encoding/base64"
	"strconv"
)

func EncodePageToken(offset int) *string {
	if offset < 0 {
		return nil
	}
	s := strconv.Itoa(offset)
	enc := base64.RawURLEncoding.EncodeToString([]byte(s))
	return &enc
}

func DecodePageToken(token string) (int, error) {
	if token == "" {
		return 0, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(raw))
}
